#!/bin/bash
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}' | ./dist/mcp-server

echo '{"jsonrpc": "2.0", "id": 12, "method": "tools/list", "params": {}}' | ./dist/mcp-server

echo '{"jsonrpc": "2.0", "id": 14, "method": "tools/call", "params": {"name": "create_order", "arguments": {"order_json": "{\"customer_email\": \"john@ejemplo.com\", \"customer_name\": \"John Doe\", \"customer_address\": \"123 Main St\", \"customer_phone\": \"555-1234\", \"items\": [{\"item_id\": 1, \"quantity\": 2}]}"}}}' | ./dist/mcp-server
