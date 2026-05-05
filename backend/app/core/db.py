from typing import AsyncGenerator

from psycopg_pool import AsyncConnectionPool

from app.core.config import settings

# LangGraph's AsyncPostgresSaver expects a standard Postgres connection string
dsn = settings.DATABASE_URL


async def get_db_pool() -> AsyncConnectionPool:
    """
    Returns an AsyncConnectionPool to be used throughout the application.
    This pool should be created on app startup and closed on shutdown.
    """
    pool = AsyncConnectionPool(conninfo=dsn, max_size=20, kwargs={"autocommit": True})
    return pool


async def get_db_connection() -> AsyncGenerator:
    """
    A dependency that yields a database connection from the pool.
    """
    pool = await get_db_pool()
    async with pool.connection() as conn:
        yield conn
