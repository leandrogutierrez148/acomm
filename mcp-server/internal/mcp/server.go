package mcp

import (
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/server"
)

// MCPHandler is an interface that any domain-specific MCP handler must implement
// to register its own tools onto the global MCPServer instance.
type MCPHandler interface {
	RegisterTools(srv *server.MCPServer)
}

// MCPServer represents the MCP server wrapper.
type MCPServer struct {
	srv      *server.MCPServer
	handlers []MCPHandler
}

// NewMCPServer initializes a new MCP server.
func NewMCPServer(handlers ...MCPHandler) *MCPServer {
	s := server.NewMCPServer(
		"catalog-server",
		"1.0.0",
	)

	mcpSrv := &MCPServer{
		srv:      s,
		handlers: handlers,
	}

	mcpSrv.registerAllTools()

	return mcpSrv
}

func (m *MCPServer) registerAllTools() {
	for _, h := range m.handlers {
		h.RegisterTools(m.srv)
	}
}

// ServeStdio runs the MCP server on standard input/output.
func (m *MCPServer) ServeStdio() error {
	log.Println("Starting MCP server on stdio...!!!!")
	return server.ServeStdio(m.srv)
}

// ServeSSE runs the MCP server over HTTP using Server-Sent Events.
func (m *MCPServer) ServeSSE(port string) error {
	sse := server.NewSSEServer(m.srv)
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Starting MCP server on SSE at %s...", addr)
	return sse.Start(addr)
}
