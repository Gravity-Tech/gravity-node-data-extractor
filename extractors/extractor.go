package extractors

import (
	"context"
	"errors"
)

const (
	String DataType = "string"
	Int64  DataType = "int64"
	Base64 DataType = "base64"
)

var (
	NotFoundErr = errors.New("data not found")
)

type DataType string

type Data struct {
	Type  DataType
	Value string
}

type ExtractorInfo struct {
	Description string `json:"description"`
	Tag         string `json:"tag"`
}

type IExtractor interface {
	Info() *ExtractorInfo
	Extract(ctx context.Context) (*Data, error)
	Aggregate(values []Data) (*Data, error)
}
