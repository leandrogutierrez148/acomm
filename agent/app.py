import asyncio
import os
import streamlit as st
import traceback

from contextlib import AsyncExitStack
from anthropic import Anthropic
from mcp import ClientSession, StdioServerParameters
from mcp.client.sse import sse_client
from dotenv import load_dotenv

load_dotenv()

# Initialize session state for chat history and connection
if "messages" not in st.session_state:
    st.session_state.messages = []
if "api_key" not in st.session_state:
    st.session_state.api_key = os.getenv("ANTHROPIC_API_KEY", "")

# We keep the Anthropic client lazily loaded based on the API key
def get_anthropic_client():
    if not st.session_state.api_key:
        return None
    return Anthropic(api_key=st.session_state.api_key)

# Basic UI Setup
st.set_page_config(page_title="A-Commerce", page_icon="🛍️")
st.title("🛍️ A-Commerce")
st.markdown("Chat with the AI agent to explore your product catalog and categories using the MCP protocol")

# Sidebar for configuration
with st.sidebar:
    st.header("Configuration")
    st.text_input("Anthropic API Key", type="password", key="api_key_input", value=st.session_state.api_key)
    
    if st.session_state.api_key_input:
        key_val = st.session_state.api_key_input.strip()
        if "\n" in key_val or len(key_val) > 250:
            st.error("Invalid API Key format detected. Please do not paste multiline text or tracebacks.")
            st.session_state.api_key = ""
        else:
            st.session_state.api_key = key_val
        
    if not st.session_state.api_key:
        st.warning("Please provide a valid Anthropic API Key to use the agent.")
    
    mcp_url = st.text_input("MCP Server SSE URL", value=os.getenv("MCP_SERVER_URL", "http://mcp-server:8080/sse"))
    st.session_state.mcp_url = mcp_url

    if st.button("Clear Chat History"):
        st.session_state.messages = []
        st.session_state.anthropic_messages = []
        st.rerun()

# Display chat messages from history
for message in st.session_state.messages:
    if message["role"] != "system":
        with st.chat_message(message["role"]):
            st.markdown(message["content"])

async def process_tool_calls_with_claude(anthropic, anthropic_tools, session, prompt):
    if "anthropic_messages" not in st.session_state:
        st.session_state.anthropic_messages = []
    
    # Append the user prompt
    st.session_state.anthropic_messages.append({"role": "user", "content": prompt})

    with st.chat_message("assistant"):
        with st.spinner("Connecting to MCP and Thinking..."):
            iteration = 0
            final_text = ""

            while iteration < 5:
                iteration += 1
                
                response = anthropic.messages.create(
                    model="claude-haiku-4-5-20251001",
                    max_tokens=2048,
                    messages=st.session_state.anthropic_messages,
                    tools=anthropic_tools,
                    system="You are an AI assistant managing a product catalog. Use the available tools to answer queries."
                )

                st.session_state.anthropic_messages.append({
                    "role": "assistant",
                    "content": response.content
                })

                if response.stop_reason == "tool_use":
                    tool_results = []
                    for block in response.content:
                        if block.type == "tool_use":
                            tool_name = block.name
                            tool_args = block.input
                            
                            st.sidebar.info(f"Calling Tool: `{tool_name}` with args: {tool_args}")
                            
                            # Execute via MCP
                            result = await session.call_tool(tool_name, arguments=tool_args)
                            
                            # Format result back to Claude
                            content_str = ""
                            if result.content:
                                content_str = getattr(result.content[0], "text", str(result.content))
                                
                            tool_results.append({
                                "type": "tool_result",
                                "tool_use_id": block.id,
                                "content": content_str
                            })
                    
                    st.session_state.anthropic_messages.append({
                        "role": "user",
                        "content": tool_results
                    })
                    # Loop continues, pushing tool results back to Claude
                else:
                    # Final response
                    for block in response.content:
                        if block.type == "text":
                            final_text += block.text
                    break

        if final_text:
            st.markdown(final_text)
            st.session_state.messages.append({"role": "assistant", "content": final_text})

async def handle_agent_response(prompt: str):
    anthropic = get_anthropic_client()
    st.session_state.messages.append({"role": "user", "content": prompt})
    with st.chat_message("user"):
        st.markdown(prompt)

    try:
        async with sse_client(url=st.session_state.mcp_url) as streams:
            async with ClientSession(streams[0], streams[1]) as session:
                await session.initialize()

                tools_response = await session.list_tools()
                
                # Convert tools to Anthropic format
                anthropic_tools = []
                for tool in tools_response.tools:
                    anthropic_tools.append({
                        "name": tool.name,
                        "description": tool.description,
                        "input_schema": tool.inputSchema
                    })

                await process_tool_calls_with_claude(anthropic, anthropic_tools, session, prompt)
    except Exception as e:
        st.error(f"Failed to connect or talk to Claude.\nError: {e}")
        st.code(traceback.format_exc())

# Process User Input
if prompt := st.chat_input("Ask a question about the catalog..."):
    if not st.session_state.api_key:
        st.info("Please set an Anthropic API key first.")
    else:
        asyncio.run(handle_agent_response(prompt))
