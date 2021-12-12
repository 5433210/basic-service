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

	"H4sIAAAAAAAC/+xXT2/bPgz9KgZ/P2AXr/bW7eJbuxVDukOGATsVQSHbTKJEkTxK7poF/u6DpDpxYzlN",
	"sXbogN1UkeGf9/RIdwOFWlVKojQasg0Q6kpJje6PCyJF9lAoaVAae2RVJXjBDFcyWWgl7Z0u5rhi9vQ/",
	"4RQy+C/ZRU28VScu2te7+NA0TQwl6oJ4ZYNB5tNFtPWIYbx8suzj5aHU42UnbxPfxdyBsP1ttoGKVIVk",
	"uMeoUKW7NesKIQNtiMuZrb1khnUMKl9gYaxhhVqzWehHTQyE32tOWEJ25UPv/CcxGG4Etkhta4r7OS5u",
	"sah9b/sFT7nkeo7lNTPBunkZvF6o/HrApGpT1eFghKYmeT0IkjaMzFAtTafhbT+DzfqXuk+OnPJZkAXJ",
	"VkdQ4LziNs5kvyBFoXouVd4vZfA5lFyzXGAX2lwpgUxaK3Z6Oyiu1m+YwYGOY5B46+CfKlpZImyt+Npw",
	"1zohK8dSrCEzVGMcYF+icw2Wr35IpDDzxRzLWoQrsgF/KnkEQbyENksnZidCB+BdrR1cO5xa3gJ0dibH",
	"8eo/KPI2YSdyL6/143Kq2gHICstRb3AtVB61bVN09mUEu/DuHmK4QdLe+81J6lipULKKQwanJ+lJCjFU",
	"zMxdQ8lC5e4wQ/cm7mcjNMTxBiMmROQ8XTByA3lUQgaf0JwJcelN95bJ2zR91DB/5JjlBlf6IZVYhpst",
	"1IyIrR8mq8/L3u74bGO88+2Fsm9h8CvQer9/hLfNWCkdoKMgZAYjZqnoMfHBGf2btopBbc5VuX6yheqw",
	"bLwe+zQf7my8fG7Imtg/5WSzUPmobDx4Ag32YfT3AzB+dMYWxpfXaBxW6gyNbyjK15Efkvs6vVT5+XpU",
	"/mGdHiHP35ejkwwjtkKDpCG72gC3NjvmoN2E4F4GdPeJX3G73vbTT9ykLOZ9vOuqHJbiN2f8J0UvxQTb",
	"j7nhRbN9vq901HEPv+IzIS66Pi9u7ey+Xv/O5fMcSnL/ZSHdtFFrEpDB3JgqSxKhCibmSpvsNE1TaCbN",
	"rwAAAP//24yhR6gOAAA=",
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
