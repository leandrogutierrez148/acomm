from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from langgraph.checkpoint.postgres.aio import AsyncPostgresSaver
from psycopg_pool import AsyncConnectionPool

from app.api.routes import api_router
from app.core.config import settings


@asynccontextmanager
async def lifespan(app: FastAPI):
    dsn = settings.DATABASE_URL

    from app.services.mcp_client import mcp_manager

    # Initialize MCP Connections
    try:
        await mcp_manager.connect_to_server("store", "http://mcp-server:8080/sse", "internal-token")
    except Exception as e:
        print(f"MCP Connection failed: {e}")

    # Initialize connection pool for LangGraph memory
    async with AsyncConnectionPool(dsn, max_size=20, kwargs={"autocommit": True}) as pool:
        app.state.db_pool = pool

        # Initialize checkpointer and create tables if they don't exist
        checkpointer = AsyncPostgresSaver(pool)
        await checkpointer.setup()
        app.state.checkpointer = checkpointer

        # Initialize chat_sessions table
        async with pool.connection() as conn:
            async with conn.cursor() as cur:
                await cur.execute("""
                    CREATE TABLE IF NOT EXISTS chat_sessions (
                        session_id TEXT PRIMARY KEY,
                        user_id TEXT,
                        title TEXT,
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                    )
                """)

        yield

    # Shutdown: Clean up resources
    print("Shutting down Acomm Backend API...")
    try:
        await mcp_manager.close_all()
    except Exception:
        pass


app = FastAPI(
    title="Acomm Purchasing Agent API",
    description="Backend for the Acomm commercial purchasing agent",
    version="1.0.0",
    lifespan=lifespan,
)

# CORS configuration
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.BACKEND_CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(api_router, prefix=settings.API_V1_STR)


@app.get("/health")
async def health_check():
    return {"status": "ok"}
