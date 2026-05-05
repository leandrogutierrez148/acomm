# Acomm Frontend

This directory contains the user interface for the Acomm purchasing agent, built with Streamlit. It provides a chat-like experience where users can interact with the agent, browse dynamically rendered product carousels, and manage their sessions.

## Features

- **Chat Interface**: Communicate with the LangGraph agent in real-time.
- **Product Carousels**: Rich rendering of products and variants returned by the agent.
- **Session Management**: Automatically restores chat history using UUID-based sessions.

## Development

The frontend is built with Python and uses `pyproject.toml` for dependency management.

### Installation

```bash
# Create and activate a virtual environment
uv venv .venv
source .venv/bin/activate  # On Windows use: .venv\Scripts\activate

# Install dependencies
uv pip install -e ".[dev]"
```

### Environment Variables

The frontend relies on the following environment variables to locate the APIs:

- `BACKEND_URL`: Complete URL to the backend's chat endpoint (default: `http://localhost:8000/api/v1/chat`).
- `AUTH_URL`: Complete URL to the auth server's authentication endpoint (default: `http://localhost:8081/api/v1/auth`).

### Running Locally

Make sure the backend is running and accessible (by default it expects the backend at `http://localhost:8000/api/v1/chat`). You can adjust this using the `BACKEND_URL` environment variable.

```bash
streamlit run app.py
```

### Formatting & Linting

```bash
ruff check .
ruff format .
mypy app.py
```
