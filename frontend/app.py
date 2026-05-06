import streamlit as st
import httpx
import os
import random
from streamlit_carousel_uui import uui_carousel
import uuid
from dotenv import load_dotenv
import logging

load_dotenv()

# Configure basic logging to stdout
logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - frontend - %(levelname)s - %(message)s"
)

# --- Configuration ---
BACKEND_URL = os.getenv("BACKEND_URL", "http://auth-server:8000/api/v1/chat")
AUTH_URL = os.getenv("AUTH_URL", "http://auth-server:8081/api/v1/auth")
DEFAULT_AUTH_TOKEN = ""

# --- Session State ---
if "messages" not in st.session_state:
    st.session_state.messages = []
if "session_id" not in st.session_state:
    st.session_state.session_id = str(uuid.uuid4())
if "auth_token" not in st.session_state:
    st.session_state.auth_token = DEFAULT_AUTH_TOKEN


def on_session_change():
    selected = st.session_state.session_selector["id"]
    if selected == "New Session":
        st.session_state.session_id = str(uuid.uuid4())
        st.session_state.messages = []
    else:
        st.session_state.session_id = selected
        url = BACKEND_URL
        token = st.session_state.auth_token
        base = url.replace("/chat", "") if url.endswith("/chat") else url
        try:
            res = httpx.get(
                f"{base}/chat/sessions/{selected}", headers={"Authorization": f"Bearer {token}"}
            )
            if res.status_code == 200:
                st.session_state.messages = res.json().get("messages", [])
        except Exception:
            st.session_state.messages = []


# --- UI Components ---
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
    name = product.get("name", "Unknown Product")

    images = product.get("images") or []
    if not images:
        for v in variants:
            images.extend(v.get("images", []))

    with st.container(border=True):
        if images:
            slides = [{"title": "", "text": "", "image": url, "interval": 3000} for url in images]
            uui_carousel(
                items=slides,
                variant="md",
                key=f"carousel_{code}_{random.randint(0, 999999)}",
            )

        st.markdown(f"**{name}**")
        st.caption(f"{category}")
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


# --- Page Layout ---
st.set_page_config(page_title="Acomm MVP UI", page_icon="🛍️", layout="wide")
st.title("🛍️ Acomm Purchasing Agent")
st.markdown("Communicate with the commercial backend via REST API.")

# --- Sidebar ---
with st.sidebar:
    st.header("Authentication")
    if not st.session_state.auth_token:
        st.info("Please log in to use the agent.")
        auth_mode = st.radio("Mode", ["Login", "Register"], horizontal=True)
        email = st.text_input("Email")
        password = st.text_input("Password", type="password")

        if st.button(auth_mode):
            if not email or not password:
                st.error("Email and password required.")
            else:
                endpoint = f"{AUTH_URL}/login" if auth_mode == "Login" else f"{AUTH_URL}/register"
                try:
                    payload = {"email": email, "password": password}
                    logging.info(
                        f"Sending {auth_mode} request to: {endpoint} with payload: {payload}"
                    )
                    res = httpx.post(endpoint, json=payload)
                    logging.info(f"Received response from auth-server: Status {res.status_code}")
                    if res.status_code in (200, 201):
                        if auth_mode == "Login":
                            st.session_state.auth_token = res.json().get("token")
                            st.success("Logged in successfully!")
                            st.rerun()
                        else:
                            st.success("Registered successfully! Please log in.")
                    else:
                        logging.warning(f"Auth failed. Response: {res.text}")
                        st.error(f"Auth failed: {res.text}")
                except Exception as e:
                    logging.error(f"Failed to connect to auth-server: {e}")
                    st.error(f"Connection failed: {e}")
    else:
        st.success("Logged in")
        if st.button("Logout"):
            st.session_state.auth_token = ""
            st.rerun()

        st.divider()
        st.header("Session Settings")

        # Fetch available sessions
        sessions = []
        base_url = (
            BACKEND_URL.replace("/chat", "")
            if BACKEND_URL.endswith("/chat")
            else BACKEND_URL
        )
        try:
            res = httpx.get(
                f"{base_url}/chat/sessions",
                headers={"Authorization": f"Bearer {st.session_state.auth_token}"},
            )
            if res.status_code == 200:
                sessions = res.json().get("sessions", [])
        except Exception as e:
            st.warning("Could not load sessions.")

        options = [{"id": "New Session", "title": "New Session"}] + sessions
        try:
            idx = next(i for i, s in enumerate(options) if s["id"] == st.session_state.session_id)
        except StopIteration:
            idx = 0

        st.selectbox(
            "Select Session",
            options,
            index=idx,
            format_func=lambda x: f"{x['title']} ({x['id'][:8]})" if x["id"] != "New Session" else x["title"],
            key="session_selector",
            on_change=on_session_change,
        )
        st.caption(f"Current ID: {st.session_state.session_id}")

        if st.button("Clear Current Chat"):
            st.session_state.messages = []
            st.rerun()

# --- Chat History ---
for message in st.session_state.messages:
    with st.chat_message(message["role"]):
        _render_message(message)

# --- Chat Input & Backend Communication ---
if prompt := st.chat_input("Ask for a product...", disabled=not st.session_state.auth_token):
    # Render user message
    st.session_state.messages.append({"role": "user", "content": prompt})
    with st.chat_message("user"):
        st.markdown(prompt)

    # Call Backend
    with st.chat_message("assistant"):
        with st.spinner("Agent is thinking..."):
            try:
                headers = {"Authorization": f"Bearer {st.session_state.auth_token}"}
                payload = {"message": prompt, "session_id": st.session_state.session_id}

                response = httpx.post(
                    BACKEND_URL, json=payload, headers=headers, timeout=60.0
                )

                if response.status_code == 200:
                    data = response.json()
                    reply_text = data.get("response", "")
                    products = data.get("products", [])

                    if reply_text:
                        st.markdown(reply_text)
                    if products:
                        render_product_cards(products)

                    # Save to state
                    msg = {"role": "assistant", "content": reply_text}
                    if products:
                        msg["products"] = products
                    st.session_state.messages.append(msg)
                else:
                    st.error(f"Backend Error [{response.status_code}]: {response.text}")

            except Exception as e:
                st.error(f"Connection failed: {e}")
