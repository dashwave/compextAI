from anthropic import Anthropic
from anthropic._types import NOT_GIVEN
def get_client(api_key):
    return Anthropic(api_key=api_key)

def chat_completion(api_key, system_prompt, model, messages, temperature, timeout, max_tokens):
    if not system_prompt:
        system_prompt = NOT_GIVEN
    client = get_client(api_key)
    response = client.messages.create(
        model=model,
        system=system_prompt,
        messages=messages,
        temperature=temperature,
        timeout=timeout,
        max_tokens=max_tokens,
    )
    
    return response.model_dump_json()
