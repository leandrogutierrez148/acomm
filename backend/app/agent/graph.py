from langgraph.graph import END, StateGraph

from app.agent.nodes import call_model, tool_node
from app.agent.state import AgentState


def should_continue(state: AgentState):
    """
    Determine if the agent should continue to execute tools or end the turn.
    """
    messages = state["messages"]
    last_message = messages[-1]

    # If the LLM makes a tool call, we transition to the 'tools' node.
    # Otherwise, we end the graph.
    if last_message.tool_calls:
        return "tools"
    return END


# Define the graph
workflow = StateGraph(AgentState)

# Define the nodes
workflow.add_node("agent", call_model)
workflow.add_node("tools", tool_node)

# Define the edges
workflow.set_entry_point("agent")

# We use a conditional edge from the agent node to either the tools node or END
workflow.add_conditional_edges("agent", should_continue, {"tools": "tools", END: END})

# After tools run, we always return to the agent to process the results
workflow.add_edge("tools", "agent")


# To compile the graph with PostgreSQL memory, we need to pass a checkpointer.
# This checkpointer will be injected from the FastAPI route or dependency.
def create_agent_app(checkpointer=None):
    if checkpointer:
        return workflow.compile(checkpointer=checkpointer)
    return workflow.compile()
