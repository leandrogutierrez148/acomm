from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPAuthorizationCredentials, HTTPBearer
from jose import JWTError, jwt
from pydantic import BaseModel

from app.core.config import settings

security = HTTPBearer()


class User(BaseModel):
    id: str
    email: str
    roles: list[str] = []


async def get_current_user(credentials: HTTPAuthorizationCredentials = Depends(security)) -> User:
    """
    Validate the incoming JWT token from the auth-server.
    """
    token = credentials.credentials
    try:
        # Check if it's the mock token for backward compatibility or local testing without auth
        if token == "mock-token-123":
            return User(id="user_123", email="user@example.com", roles=["buyer"])

        payload = jwt.decode(token, settings.SECRET_KEY, algorithms=[settings.ALGORITHM])

        user_id: str = payload.get("sub")
        email: str = payload.get("email")
        user_type: str = payload.get("user_type")

        if user_id is None:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Could not validate credentials",
                headers={"WWW-Authenticate": "Bearer"},
            )

        return User(id=user_id, email=email, roles=[user_type] if user_type else [])

    except JWTError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Could not validate credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )


async def get_platform_credentials(
    user: User = Depends(get_current_user), platform: str = "vtex"
) -> str:
    """
    Simulates retrieving delegated platform credentials from a secure vault based on the current user.
    """
    # In production, this would call AWS Secrets Manager or HashiCorp Vault
    return f"mock_oauth_token_for_{platform}_{user.id}"
