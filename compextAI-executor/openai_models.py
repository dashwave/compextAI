import openai

def get_client(api_key):
    return openai.OpenAI(
        api_key=api_key
    )

def chat_completion(api_key, model, messages, temperature, timeout, max_completion_tokens, response_format):
    client = get_client(api_key)
    response = client.chat.completions.create(
        model=model,
        messages=messages,
        temperature=temperature,
        timeout=timeout,
        max_completion_tokens=max_completion_tokens,
        response_format=response_format,
    )
    return response.model_dump_json()
