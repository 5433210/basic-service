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

	"H4sIAAAAAAAC/+xXTW/jNhD9K8K0QC9qpDTtRbekDQqnBxcL7CkIAkoa27RpUjukvPEa+u8LkpGsWJTj",
	"YJNFFtibohm/mXlvPpQdFGpdKYnSaMh2QKgrJTW6P66JFNmHQkmD0thHVlWCF8xwJZOlVtK+08UC18w+",
	"/Uo4gwx+SfaoibfqxKF9eMSHpmliKFEXxCsLBpkPF1HnEcN09WrRp6tjoaerXtwmfsTck9D9NttBRapC",
	"MtxzVKjSvTXbCiEDbYjLuc29ZIb1DCpfYmGsYY1as3noR00MhJ9qTlhCduuh9/53MRhuBLZMdTnFwxjX",
	"D1jUvrbDhGdccr3A8p6ZYN68DL5eqvx+xKRqU9VhMEJTk7wfJUkbRmYslwM6eAldGn0uulJHefBNfKib",
	"nPF5UCDJ1ieo47ziFmeQkKJQPjcqH6Yy2ikl1ywX2Gc9V0ogk9aKvdqOzl3rNy7uSMUxSHxwyswUra1G",
	"Nlf83XBXOiErp1JsITNUYxxoDInONZi++iyRwk1RLLCsRTgjC/hFSTyxX3yUHmYPoUfwPtcerz1NrW4B",
	"OXtL5fTFcHT+24A95EFc68flTLW7kRVWo8FOW6o8asum6PL/Cezh3XuIYYOkvff5WepUqVCyikMGF2fp",
	"WQoxVMwsXEHJUuXuYY6uJ55GIzTEcYMREyJyng6M3K6elJDBv2guhbjxpid35o80fdGef+EG5gbX+rkp",
	"sQo3HdWMiG2fF2uoy8FZ+c9i/OnLC0XvaPDX0Xr/9QJvG7FSOiBHQcgMRsxKMVDib2f0PW0nBrW5UuX2",
	"1W6t47Lx8ziU+Xhl09VbU9bEvpWT3VLlk7Lx5Ak0OKTRvx+h8R9nbGl8f4XG4Umdo/EFRfk28kvycE5v",
	"VH61nZTfeU5PGM9vH0c3MozYGg2Shux2B9za7JqD9hKC6wzo3xN/4va1HYa/c5uyWAz5rqtyfBQ/OuPP",
	"UfSjmGD7MTd+aLr2/U1HPfdwF18Kcd33eXdnZ//1+mMen7eYJPcPGNKmRa1JQAYLY6osSYQqmFgobbKL",
	"NE2TzTk0d83XAAAA//+sHF63xg4AAA==",
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
