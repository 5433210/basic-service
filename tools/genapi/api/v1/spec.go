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

	"H4sIAAAAAAAC/+xXT2/bPgz9KgZ/P2AXr/bW7eJbuxVDukOGATsVQSHbTKJEkTRK7poF/u6DpDpxYztN",
	"sXbogN1UkeGf9/RIdwOFWmklUVoD2QYIjVbSoP/jgkiROxRKWpTWHZnWghfMciWThVHS3ZlijivmTv8T",
	"TiGD/5Jd1CRYTeKjfb2LD3Vdx1CiKYhrFwyykC6irUcM4+WTZR8vD6UeL1t56/gu5g6E7W+zDWhSGsny",
	"gFGhSn9r1xohA2OJy5mrvWSWtQwqX2BhnWGFxrBZ34/qGAi/V5ywhOwqhN75T2Kw3ApskNrWFHdzXNxi",
	"UYXe9guecsnNHMtrZnvr5mXv9ULl1wMmVVld9QcjtBXJ60GQjGVkh2qpWw1v+xlsNrzUfXLklM96WZBs",
	"dQQF3itu4kz2C1LUV8+lyrulDD6HkhuWC2xDmyslkElnxVZvB8XV+A0zONBxDBJvPfxTRStHhKsVX1vu",
	"Wydk5ViKNWSWKox72P8hkfrJLeZYVqI/qQv/U8kjOOAltNxbgN1V3gKpRZAjoYeb1hg4XsoHFdskbEXu",
	"5HV+XE5VM81Y4QDvTKGFyqMGNIrOvoxgF97fQww3SCZ4vzlJPf4aJdMcMjg9SU9SiEEzO/cNJQuV+8MM",
	"PcH3sxFa4niDERMi8p4+GPnpOiohg09oz4S4DKZ7m+Ftmj5qMj9yZnKLK/PQk3cM11uoGRFbP0xWl5e9",
	"RfDZxXgX2uvLvoUh7DPn/f4R3i6jVqaHjoKQWYyYo6LDxAdvDG/aaQONPVfl+sm2o8eyDsrr0ny4s/Hy",
	"uSGr4/CUk81C5aOyDuAJtNiFMdwPwPjRGxsYX16jcb9SZ2hDQ1G+jvw47Oj0UuXn61H5h3V6hDx/X45e",
	"MozYCi2SgexqA9zZ3JiDZq2BfxnQ3hxhX+16208/8ZOymHfxrnQ5LMVv3vhPikGKCTZfZsOLZvt8X5mo",
	"5d7/is+EuGj7vLi1s/sU/TuXz3Moyf/LhHTTRK1IQAZza3WWJEIVTMyVsdlpmqZQT+pfAQAA//9GC8oe",
	"dQ4AAA==",
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