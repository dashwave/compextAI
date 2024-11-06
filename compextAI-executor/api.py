import fastapi
import uvicorn
import os
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from openai_models import chat_completion

app = fastapi.FastAPI()

@app.get("/")
def read_root():
    return {"pong"}

class ChatCompletionRequest(BaseModel):
    """
    Request body for the chat completion endpoint.
    """
    api_key: str
    model: str
    messages: list[dict]
    temperature: float = 0.5
    timeout: int = 600
    max_completion_tokens: int = 10000

@app.post("/chatcompletion/openai")
def chat_completion_openai(request: ChatCompletionRequest):
    try:
        response = chat_completion(request.api_key, request.model, request.messages, request.temperature, request.timeout, request.max_completion_tokens)
        return JSONResponse(status_code=200, content={"role": "assistant", "content": response})
    except Exception as e:
        print(e)
        return JSONResponse(status_code=500, content={"error": str(e)})

if __name__ == "__main__":
    port = 8889
    if os.getenv("SERVER_PORT"):
        port = int(os.getenv("SERVER_PORT"))
    uvicorn.run(app, host="0.0.0.0", port=port)
