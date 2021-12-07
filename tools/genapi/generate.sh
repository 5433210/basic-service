#/bin/bash

base_dir=./
module=api
version=v1
package=$module$version
spec="/Users/zhangweili/Desktop/rbac/api/openapi/captcha.yaml"

# go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen
# oapi-codegen petstore-expanded.yaml  > petstore.gen.go
output_dir=$base_dir/$module/$version
mkdir -p $output_dir

target=$output_dir/types.go
oapi-codegen -generate types -o $target -package $package $spec
target=$output_dir/server.go
oapi-codegen -generate server -o $target -package $package $spec
target=$output_dir/spec.go
oapi-codegen -generate spec -o $target -package $package $spec
target=$output_dir/client.go
oapi-codegen -generate client -o $target -package $package $spec
