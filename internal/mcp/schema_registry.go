package mcp

import (
	"embed"
	"fmt"
	"strings"
)

const (
	CreateOrderSchemaName = "create_order"
)

//go:embed schemas/*.json
var toolSchemasFS embed.FS

func GetToolSchema(name string) (string, error) {
	normalizedName := strings.TrimSpace(name)
	if normalizedName == "" {
		return "", fmt.Errorf("schema name is required")
	}

	fileName := fmt.Sprintf("schemas/%s.schema.json", normalizedName)
	schemaBytes, err := toolSchemasFS.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to load schema %q: %w", normalizedName, err)
	}

	return string(schemaBytes), nil
}

func MustGetToolSchema(name string) string {
	schema, err := GetToolSchema(name)
	if err != nil {
		panic(err)
	}
	return schema
}
