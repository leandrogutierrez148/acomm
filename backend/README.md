# Acomm Backend

This is the commercial backend for the Acomm purchasing agent, built with FastAPI and LangGraph. It handles the core AI workflows, maintains conversation state using PostgreSQL, and interacts with the MCP Server to retrieve e-commerce data.

## Development

This project uses modern Python packaging via `pyproject.toml`.

To install dependencies with pip:
```bash
# Create and activate a virtual environment
uv venv .venv
source .venv/bin/activate  # On Windows use: .venv\Scripts\activate

# Install dependencies in editable mode
uv pip install -e ".[dev]"
```

### Formatting & Linting
We use `ruff` for fast linting and formatting, and `mypy` for static type checking.

```bash
ruff check .
ruff format .
mypy app/
```

### Running Locally

To run the FastAPI server locally (outside of Docker):

```bash
uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload
```
