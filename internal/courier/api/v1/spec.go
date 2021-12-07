// Package apiv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.3 DO NOT EDIT.
package apiv1

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xXTW/bOBD9KwF3j9xI27QX3dKgQIMiTRDnFhgFI41jBuJHSbqIEei/F/ywREmMLBXp",
	"oUBvMjnz5s3jcIZ+QaVgUnDgRqPiBSnQUnAN7scnpYSyH6XgBrixn0TKmpbEUMGzJy24XdPlFhixX/8q",
	"2KAC/ZN1qJnf1ZlDuw34qGkajCrQpaLSgqHChztRnQUOyJ4MI7S2H1IJCcpQzzGiNhncel8E2wajjRJs",
	"ls95VSnQ2vookPX+mxHWjxpgeimA2UtABSJKkb39rXcPT1A67mFLG0X5o7N9qzCO+PcdVVCh4t4n7uC7",
	"+LhVcY2Roaa27l7vFk140wajXjhLPXY4rONxRr0jGJ0jExUsOcQra99gJMm+FqSyroLD9QYV9zNA7oDJ",
	"mhiLMG19B8/mqNFnw2rUrIdSu5Q6huuhlhg9/0e5AcVJjYoNqTUMZLoKovQkjvdek7nNb6RzRQxJFhyt",
	"EsuDjGiFsEcYFkobMFUwvYufuMEhySGjAdUOkIHW5BGO8y29Rgf7mHSPU4K0O9MR121YnQ7rrKJgDisR",
	"42YrOHzdsQdQ8SnHy4kDjrav2/7Rp8lbzKmyjeO4DB6p7+dTTrfO6sLdvkHaAQAfwkcKjDkn5IigIzWi",
	"1YQYK+DVyhCz0wsKi1bADd3QIPtwe3Z5RTh4otYijomsV0z/8lRbMb1wpo0PYtmsSftPDRw3aRLzxead",
	"luPVGRGpMq/Xr5he2OnX9vxnjKKO5lXqKoTWn077kN9k9sPWP9hJXYUo2ZF0ye6OEScsfUUkUcRXw6yy",
	"iGLfWM+jVeEmSQjSV2dqkozC9AXq7yUkcmc80saE1enL7qwiqg5rxNF6Ub4Rh3Ilrj2PnroXYqcoqJPz",
	"m0vUQZZ+FWH0A5T2lv+f5pa5kMCJpKhAZ6f5ae7EM1tHP4P2ZSy0S6QfSwOvTiC85mze7vF+WYXOdHjn",
	"2VxBm4+i2r/dk99hN43XMvpn8S7PFwVZ+GLom+t2Psyr5a5dH6vigLzGb/dGid5X48rqH+z1FxvovZcy",
	"lVEruf/vZa0/LLB2ITMd5tPrtaVdJx9Xlu/wv6OuLPLfqvpjq8r+9QVle5yb4DtVowJtjZFFltWiJPVW",
	"aFOc5XluB/LPAAAA//8fSPkYoxAAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
