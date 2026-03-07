#!/bin/bash
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}' | ./dist/mcp-server

echo '{"jsonrpc": "2.0", "id": 12, "method": "tools/list", "params": {}}' | ./dist/mcp-server

echo '{"jsonrpc": "2.0", "id": 14, "method": "tools/call", "params": {"name": "create_order", "arguments": {"order_json": "{\"customerEmail\": \"john@ejemplo.com\", \"customerName\": \"John Doe\", \"customerAddress\": \"123 Main St\", \"customerPhone\": \"555-1234\", \"items\": [{\"itemId\": 1, \"quantity\": 2}]}"}}}' | ./dist/mcp-server