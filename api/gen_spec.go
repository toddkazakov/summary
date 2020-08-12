// Package Summary provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package Summary

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"strings"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
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

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}

