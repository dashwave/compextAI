import openai
from openai._types import NOT_GIVEN

def get_client(api_key):
    return openai.OpenAI(
        api_key=api_key
    )

def chat_completion(api_key:str, model:str, messages:list, temperature:float, timeout:int, max_completion_tokens:int, response_format:dict):
    client = get_client(api_key)
    if response_format is None or response_format == {}:
        response_format = NOT_GIVEN

    response = client.chat.completions.create(
        model=model,
        messages=messages,
        temperature=temperature,
        timeout=timeout,
        max_completion_tokens=max_completion_tokens,
        response_format=response_format,
    )
    return response.model_dump_json()
