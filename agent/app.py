import asyncio
import os
import json
import random
import streamlit as st
import traceback
from streamlit_carousel_uui import uui_carousel

from anthropic import Anthropic
from mcp import ClientSession
from mcp.client.sse import sse_client
from dotenv import load_dotenv

load_dotenv()

# Session state
if "messages" not in st.session_state:
    st.session_state.messages = []
if "api_key" not in st.session_state:
    st.session_state.api_key = os.getenv("ANTHROPIC_API_KEY", "")


def get_anthropic_client():
    """Return a configured Anthropic client, or None if no API key is set."""
    if not st.session_state.api_key:
        return None
    return Anthropic(api_key=st.session_state.api_key)


def get_system_prompt():
    """Load the agent system prompt from disk, falling back to a default."""
    prompt_path = os.path.join(os.path.dirname(__file__), "system.prompt")
    try:
        with open(prompt_path, "r", encoding="utf-8") as f:
            return f.read()
    except Exception:
        return """You are Uncle Donald, you are old and angry, you dont like to work"""


# Product rendering constants
PRODUCT_TOOLS = {"search_products", "get_product_by_code", "create_product"}

PLACEHOLDER_IMAGES = [
    "https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=600&q=80",
    "https://images.unsplash.com/photo-1585386959984-a4155224a1ad?w=600&q=80",
    "https://images.unsplash.com/photo-1491553895911-0055eca6402d?w=600&q=80",
    "https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f?w=600&q=80",
]


def extract_products_from_result(tool_name: str, content_str: str) -> list[dict]:
    """Deserialise a tool result payload into a list of product dicts."""
    try:
        data = json.loads(content_str)
    except json.JSONDecodeError:
        return []

    if tool_name == "search_products":
        products = data.get("products", [])
        if isinstance(products, list):
            return products

    elif tool_name in ("get_product_by_code", "create_product"):
        if "code" in data:
            return [data]

    return []


def render_product_cards(products: list[dict]) -> None:
    """Render a grid of product cards from a list of product dicts."""
    if not products:
        return

    cols_per_row = 3
    for row_start in range(0, len(products), cols_per_row):
        row_products = products[row_start : row_start + cols_per_row]
        cols = st.columns(cols_per_row, gap="medium")
        for col, product in zip(cols, row_products):
            with col:
                _render_single_card(product)


def _render_single_card(product: dict) -> None:
    """Render a single product card including image carousel, price, and variants."""
    code = product.get("code", "—")
    category = product.get("category", "Uncategorized")
    price = product.get("price", 0.0)
    variants = product.get("variants") or []

    with st.container(border=True):
        slides = [
            dict(title="", text="", image=url, interval=3000)
            for url in PLACEHOLDER_IMAGES
        ]
        uui_carousel(
            items=slides,
            variant="md",
            key=f"carousel_{code}_{random.randint(0, 999999)}",
        )

        st.markdown(f"**{category}**")
        st.caption(f"Code: `{code}`")
        st.markdown(
            f"<span style='font-size:1.4rem;font-weight:700'>${price:,.2f}</span>",
            unsafe_allow_html=True,
        )

        if variants:
            with st.expander(f"{len(variants)} variant(s)"):
                for v in variants:
                    sku = v.get("sku") or v.get("name", "—")
                    v_price = v.get("price", price)
                    st.write(f"• **{sku}** — ${v_price:,.2f}")


def _render_message(message: dict) -> None:
    """Render a persisted chat message, including any embedded product cards."""
    text = message.get("content", "")
    if text:
        st.markdown(text)
    products = message.get("products")
    if products:
        render_product_cards(products)


# Page layout
st.set_page_config(page_title="Acomm", page_icon="🛍️", layout="wide")
st.title("🛍️ E-Commerce for Agents")
st.markdown(
    "Chat with the AI agent to explore your product catalog and categories using the MCP protocol"
)

# Sidebar
with st.sidebar:
    st.header("Configuration")
    st.text_input(
        "Anthropic API Key",
        type="password",
        key="api_key_input",
        value=st.session_state.api_key,
    )

    if st.session_state.api_key_input:
        key_val = st.session_state.api_key_input.strip()
        if "\n" in key_val or len(key_val) > 250:
            st.error(
                "Invalid API Key format detected. Please do not paste multiline text or tracebacks."
            )
            st.session_state.api_key = ""
        else:
            st.session_state.api_key = key_val

    if not st.session_state.api_key:
        st.warning("Please provide a valid Anthropic API Key to use the agent.")

    mcp_url = st.text_input(
        "MCP Server SSE URL",
        value=os.getenv("MCP_SERVER_URL", "http://mcp-server:8080/sse"),
    )
    st.session_state.mcp_url = mcp_url

    if st.button("Clear Chat History"):
        st.session_state.messages = []
        st.session_state.anthropic_messages = []
        st.rerun()

# Chat history
for message in st.session_state.messages:
    if message["role"] != "system":
        with st.chat_message(message["role"]):
            _render_message(message)


# Agent logic
async def process_tool_calls_with_claude(anthropic, anthropic_tools, session, prompt):
    """Drive a multi-turn Claude conversation, dispatching MCP tool calls until a final response is produced."""
    if "anthropic_messages" not in st.session_state:
        st.session_state.anthropic_messages = []

    st.session_state.anthropic_messages.append({"role": "user", "content": prompt})

    with st.chat_message("assistant"):
        with st.spinner("Connecting to MCP and Thinking..."):
            iteration = 0
            final_text = ""
            collected_products: list[dict] = []

            while iteration < 5:
                iteration += 1

                response = anthropic.messages.create(
                    model="claude-haiku-4-5-20251001",
                    max_tokens=2048,
                    messages=st.session_state.anthropic_messages,
                    tools=anthropic_tools,
                    system=get_system_prompt(),
                )

                st.session_state.anthropic_messages.append(
                    {"role": "assistant", "content": response.content}
                )

                if response.stop_reason == "tool_use":
                    tool_results = []
                    for block in response.content:
                        if block.type != "tool_use":
                            continue

                        tool_name = block.name
                        tool_args = block.input

                        st.sidebar.info(
                            f"Calling Tool: `{tool_name}` with args: {tool_args}"
                        )

                        result = await session.call_tool(tool_name, arguments=tool_args)

                        content_str = ""
                        if result.content:
                            content_str = getattr(
                                result.content[0], "text", str(result.content)
                            )

                        tool_results.append(
                            {
                                "type": "tool_result",
                                "tool_use_id": block.id,
                                "content": content_str,
                            }
                        )

                        if tool_name in PRODUCT_TOOLS:
                            products = extract_products_from_result(
                                tool_name, content_str
                            )
                            collected_products.extend(products)

                    st.session_state.anthropic_messages.append(
                        {"role": "user", "content": tool_results}
                    )

                else:
                    for block in response.content:
                        if block.type == "text":
                            final_text += block.text
                    break

        if final_text:
            st.markdown(final_text)

        if collected_products:
            render_product_cards(collected_products)

        msg: dict = {"role": "assistant", "content": final_text}
        if collected_products:
            msg["products"] = collected_products
        st.session_state.messages.append(msg)


async def handle_agent_response(prompt: str):
    """Establish an MCP session and delegate the user prompt to the Claude agent."""
    anthropic = get_anthropic_client()
    st.session_state.messages.append({"role": "user", "content": prompt})
    with st.chat_message("user"):
        st.markdown(prompt)

    try:
        async with sse_client(url=st.session_state.mcp_url) as streams:
            async with ClientSession(streams[0], streams[1]) as session:
                await session.initialize()

                tools_response = await session.list_tools()

                anthropic_tools = [
                    {
                        "name": tool.name,
                        "description": tool.description,
                        "input_schema": tool.inputSchema,
                    }
                    for tool in tools_response.tools
                ]

                await process_tool_calls_with_claude(
                    anthropic, anthropic_tools, session, prompt
                )
    except Exception as e:
        st.error(f"Failed to connect or talk to Claude.\nError: {e}")
        st.code(traceback.format_exc())


# Entry point
if prompt := st.chat_input("Ask a question about the catalog..."):
    if not st.session_state.api_key:
        st.info("Please set an Anthropic API key first.")
    else:
        asyncio.run(handle_agent_response(prompt))
