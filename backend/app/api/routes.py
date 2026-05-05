from fastapi import APIRouter

from app.api.endpoints import chat

api_router = APIRouter()

# We will include routers here once they are implemented
# api_router.include_router(auth.router, prefix="/auth", tags=["auth"])
api_router.include_router(chat.router, prefix="/chat", tags=["chat"])
