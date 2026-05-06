from fastapi import APIRouter, Depends, HTTPException, Request
from langchain_core.messages import HumanMessage

from app.agent.graph import create_agent_app
from app.api.deps import User, get_current_user
from app.models.chat import ChatRequest, ChatResponse

router = APIRouter()


@router.post("", response_model=ChatResponse)
async def chat_endpoint(
    request: Request, payload: ChatRequest, current_user: User = Depends(get_current_user)
):
    """
    Process a chat message from the user and route it through the LangGraph agent.
    """

    # Retrieve checkpointer from application state
    checkpointer = request.app.state.checkpointer

    # Instantiate agent per request (compilation is fast, but cacheable if needed)
    agent_app = create_agent_app(checkpointer)

    # Initialize state
    initial_state = {
        "messages": [HumanMessage(content=payload.message)],
        "user_id": current_user.id,
        "session_id": payload.session_id,
        "context": {},
        "collected_products": [],
    }

    config = {"configurable": {"thread_id": payload.session_id}}

    # Check if this is the first message by checking if a checkpoint already exists
    existing_checkpoint = await checkpointer.aget(config)
    is_first_message = existing_checkpoint is None

    if is_first_message:
        # Store session title if it's the first message
        pool = request.app.state.db_pool
        try:
            async with pool.connection() as conn:
                async with conn.cursor() as cur:
                    title = payload.message[:50] + ("..." if len(payload.message) > 50 else "")
                    await cur.execute(
                        "INSERT INTO chat_sessions (session_id, user_id, title) VALUES (%s, %s, %s) ON CONFLICT (session_id) DO NOTHING",
                        (payload.session_id, current_user.id, title),
                    )
        except Exception as e:
            print(f"Failed to save session title: {e}")

    try:
        # Run the agent graph
        final_state = await agent_app.ainvoke(initial_state, config=config)

        # Extract the final AI message
        messages = final_state.get("messages", [])
        if not messages:
            raise HTTPException(status_code=500, detail="No response from agent")

        last_message = messages[-1]

        return ChatResponse(
            response=last_message.content
            if isinstance(last_message.content, str)
            else str(last_message.content),
            products=final_state.get("collected_products", []),
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/sessions")
async def get_sessions(request: Request, current_user: User = Depends(get_current_user)):
    """List all available session IDs."""
    pool = request.app.state.db_pool
    try:
        async with pool.connection() as conn:
            async with conn.cursor() as cur:
                await cur.execute(
                    "SELECT session_id, title FROM chat_sessions WHERE user_id = %s ORDER BY created_at DESC",
                    (current_user.id,),
                )
                rows = await cur.fetchall()
                return {"sessions": [{"id": row[0], "title": row[1]} for row in rows]}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/sessions/{session_id}")
async def get_session_history(
    session_id: str, request: Request, current_user: User = Depends(get_current_user)
):
    """Fetch the message history for a specific session."""
    checkpointer = request.app.state.checkpointer
    config = {"configurable": {"thread_id": session_id}}
    try:
        checkpoint = await checkpointer.aget(config)
        if not checkpoint:
            return {"messages": []}

        state = checkpoint.get("channel_values", {})
        messages = state.get("messages", [])

        formatted_messages = []
        for msg in messages:
            if msg.type == "system":
                continue
            role = "user" if msg.type == "human" else "assistant"
            content = msg.content if isinstance(msg.content, str) else str(msg.content)

            # Simple conversion for MVP history
            formatted_messages.append({"role": role, "content": content})

        return {"messages": formatted_messages}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
