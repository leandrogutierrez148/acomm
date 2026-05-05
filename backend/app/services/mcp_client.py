from mcp import ClientSession
from mcp.client.sse import sse_client


class MCPConnectionManager:
    """
    Manages connections to external MCP servers (e.g., VTEX MCP Bridge, Mercado Libre MCP Bridge).
    This allows the LangGraph agent to discover and use tools from multiple platforms.
    """

    def __init__(self):
        # Map of platform -> (session, sse_client)
        self.connections = {}
        self.tools = {}

    async def connect_to_server(self, platform_name: str, url: str, user_token: str):
        """
        Connects to a specific MCP server.
        Passes the delegated user_token (e.g. via headers) to authorize the session on the bridge.
        """
        # In a real implementation, you'd pass the auth token in the headers or connection params
        # headers = {"Authorization": f"Bearer {user_token}"}

        try:
            # We must hold references to the stream contexts to keep them open
            sse_ctx = sse_client(url=url)
            streams = await sse_ctx.__aenter__()

            session_ctx = ClientSession(streams[0], streams[1])
            session = await session_ctx.__aenter__()

            await session.initialize()
            tools_response = await session.list_tools()

            self.connections[platform_name] = {
                "sse_ctx": sse_ctx,
                "session_ctx": session_ctx,
                "session": session,
            }

            self.tools[platform_name] = tools_response.tools
            print(f"Successfully connected to MCP Server: {platform_name}")

        except Exception as e:
            print(f"Failed to connect to MCP Server {platform_name}: {e}")

    async def get_all_tools(self) -> list:
        """
        Aggregates tools from all connected MCP servers and formats them for Anthropic.
        """
        anthropic_tools = []
        for platform, tools in self.tools.items():
            for tool in tools:
                # We prefix the tool name with the platform to avoid collisions
                anthropic_tools.append(
                    {
                        "name": f"{platform}_{tool.name}",
                        "description": tool.description,
                        "input_schema": tool.inputSchema,
                    }
                )
        return anthropic_tools

    def get_langchain_tools(self) -> list:
        """
        Returns tools formatted for LangChain's bind_tools().
        Each tool is a dict with 'type': 'function' and nested 'function' details.
        """
        lc_tools = []
        for platform, tools in self.tools.items():
            for tool in tools:
                lc_tools.append(
                    {
                        "type": "function",
                        "function": {
                            "name": f"{platform}_{tool.name}",
                            "description": tool.description or "",
                            "parameters": tool.inputSchema,
                        },
                    }
                )
        return lc_tools

    async def call_tool(self, tool_call_name: str, arguments: dict):
        """
        Routes the tool call back to the correct MCP server.
        """
        # Extract platform prefix
        parts = tool_call_name.split("_", 1)
        if len(parts) != 2:
            raise ValueError(f"Invalid tool name format: {tool_call_name}")

        platform_name, original_tool_name = parts[0], parts[1]

        if platform_name not in self.connections:
            raise ValueError(f"No active connection to platform: {platform_name}")

        session = self.connections[platform_name]["session"]
        result = await session.call_tool(original_tool_name, arguments=arguments)
        return result

    async def close_all(self):
        """Clean up connections."""
        for platform, conns in self.connections.items():
            await conns["session_ctx"].__aexit__(None, None, None)
            await conns["sse_ctx"].__aexit__(None, None, None)


# Singleton instance for the app
mcp_manager = MCPConnectionManager()
