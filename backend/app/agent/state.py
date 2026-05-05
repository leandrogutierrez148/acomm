from typing import Annotated, Any, Dict, Sequence, TypedDict

from langchain_core.messages import BaseMessage
from langgraph.graph.message import add_messages


class AgentState(TypedDict):
    """
    The state of the purchasing agent.
    - messages: The conversation history with the user (managed by LangGraph's add_messages reducer).
    - user_id: The ID of the authenticated user.
    - session_id: The active session ID.
    - context: Any additional context (e.g., active platform, preferred language).
    - collected_products: Temporary storage for products found during the turn.
    """

    messages: Annotated[Sequence[BaseMessage], add_messages]
    user_id: str
    session_id: str
    context: Dict[str, Any]
    collected_products: list[dict]
