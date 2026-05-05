from typing import List, Optional

from pydantic import BaseModel


class ChatRequest(BaseModel):
    message: str
    session_id: str


class ProductVariant(BaseModel):
    sku: Optional[str] = None
    name: Optional[str] = None
    price: Optional[float] = None
    images: Optional[List[str]] = None


class Product(BaseModel):
    code: str
    name: str
    category: Optional[str] = None
    price: float
    images: Optional[List[str]] = None
    variants: Optional[List[ProductVariant]] = None


class ChatResponse(BaseModel):
    response: str
    products: Optional[List[Product]] = None
