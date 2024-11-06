import fastapi
import uvicorn
import os
from fastapi.responses import JSONResponse
from pydantic import BaseModel
import openai_models as openai
import anthropic_models as anthropic

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
    max_tokens: int = 10000
    max_completion_tokens: int = 10000
    response_format: dict = None
    system_prompt: str = None

@app.post("/chatcompletion/openai")
def chat_completion_openai(request: ChatCompletionRequest):
    try:
        response = openai.chat_completion(request.api_key, request.model, request.messages, request.temperature, request.timeout, request.max_completion_tokens, request.response_format)
        return JSONResponse(status_code=200, content={"role": "assistant", "content": response})
    except Exception as e:
        print(e)
        return JSONResponse(status_code=500, content={"error": str(e)})
    
@app.post("/chatcompletion/anthropic")
def chat_completion_anthropic(request: ChatCompletionRequest):
    try:
        response = anthropic.chat_completion(request.api_key, request.system_prompt, request.model, request.messages, request.temperature, request.timeout, request.max_tokens)
        return JSONResponse(status_code=200, content={"role": "assistant", "content": response})
    except Exception as e:
        print(e)
        return JSONResponse(status_code=500, content={"error": str(e)})

if __name__ == "__main__":
    port = 8889
    if os.getenv("SERVER_PORT"):
        port = int(os.getenv("SERVER_PORT"))
    uvicorn.run(app, host="0.0.0.0", port=port)

