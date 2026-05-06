import json
import os

from langchain_anthropic import ChatAnthropic
from langchain_core.messages import SystemMessage, ToolMessage

from app.agent.state import AgentState
from app.core.config import settings
from app.services.mcp_client import mcp_manager

# Initialize the LLM
llm = ChatAnthropic(
    model="claude-haiku-4-5-20251001",
    api_key=settings.ANTHROPIC_API_KEY,
    temperature=0,
)


def get_system_prompt() -> str:
    prompt_path = os.path.join(os.path.dirname(__file__), "system.prompt")
    try:
        with open(prompt_path, "r", encoding="utf-8") as f:
            return f.read()
    except Exception:
        return "You are a helpful purchasing assistant."


def get_llm_with_tools():
    """
    Bind the MCP tools to the LLM so Claude knows which tools are available.
    Returns the bound LLM instance.
    """
    mcp_tools = mcp_manager.get_langchain_tools()
    if mcp_tools:
        return llm.bind_tools(mcp_tools)
    return llm


async def call_model(state: AgentState) -> dict:
    """
    Invoke the LLM with the current conversation history.
    The LLM is bound with MCP tools so it can request tool calls.
    """
    messages = state["messages"]
    sys_msg = SystemMessage(content=get_system_prompt())

    # Bind tools dynamically (tools may change if new MCP servers connect)
    bound_llm = get_llm_with_tools()

    response = await bound_llm.ainvoke([sys_msg] + list(messages))

    # We clear collected_products ONLY at the start of a new turn (when last message was from human)
    # This prevents wiping out products collected by the tool_node in the middle of a loop.
    is_new_turn = len(messages) > 0 and messages[-1].type == "human"

    update = {"messages": [response]}
    if is_new_turn:
        update["collected_products"] = []

    return update


async def tool_node(state: AgentState) -> dict:
    """
    Execute the tool calls requested by the LLM via the MCP connection manager.
    Each tool call is routed to the correct MCP server and the result is returned
    as a ToolMessage so LangGraph can feed it back to the agent.
    """
    messages = state["messages"]
    last_message = messages[-1]

    tool_messages = []
    collected_products = []

    for tool_call in last_message.tool_calls:
        tool_name = tool_call["name"]
        tool_args = tool_call["args"]
        tool_call_id = tool_call["id"]

        try:
            result = await mcp_manager.call_tool(tool_name, tool_args)

            # MCP returns a result object; extract text content
            if hasattr(result, "content"):
                # result.content is a list of content blocks
                content_parts = []
                for block in result.content:
                    if hasattr(block, "text"):
                        content_parts.append(block.text)
                    else:
                        content_parts.append(str(block))
                result_text = "\n".join(content_parts)
            else:
                result_text = json.dumps(result) if not isinstance(result, str) else result

            # Try to parse the result to extract products for the UI
            try:
                parsed_result = json.loads(result_text)
                print(f"DEBUG: Parsed result type: {type(parsed_result)}")
                if isinstance(parsed_result, dict):
                    print(f"DEBUG: Parsed result keys: {parsed_result.keys()}")
                    if "products" in parsed_result and isinstance(parsed_result["products"], list):
                        collected_products.extend(parsed_result["products"])
                        print(
                            f"DEBUG: Extracted {len(parsed_result['products'])} products via 'products' key"
                        )
                    elif "code" in parsed_result and "name" in parsed_result:
                        collected_products.append(parsed_result)
                        print("DEBUG: Extracted 1 product via code/name")
                elif isinstance(parsed_result, list):
                    # For endpoints returning a raw list, check if it looks like a product/item
                    if (
                        len(parsed_result) > 0
                        and isinstance(parsed_result[0], dict)
                        and ("sku" in parsed_result[0] or "code" in parsed_result[0])
                    ):
                        collected_products.extend(parsed_result)
                        print(f"DEBUG: Extracted {len(parsed_result)} products via list")
            except (json.JSONDecodeError, TypeError) as e:
                print(f"DEBUG: Error parsing result: {e}")

            tool_messages.append(
                ToolMessage(
                    content=result_text,
                    tool_call_id=tool_call_id,
                )
            )
        except Exception as e:
            tool_messages.append(
                ToolMessage(
                    content=f"Error executing tool '{tool_name}': {str(e)}",
                    tool_call_id=tool_call_id,
                )
            )

    return {"messages": tool_messages, "collected_products": collected_products}
