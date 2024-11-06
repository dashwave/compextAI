from anthropic import Anthropic

def get_client(api_key):
    return Anthropic(api_key=api_key)

def chat_completion(api_key, system_prompt, model, messages, temperature, timeout, max_tokens):
    print(f"system_prompt: {system_prompt}")
    client = get_client(api_key)
    response = client.messages.create(
        model=model,
        system=system_prompt,
        messages=messages,
        temperature=temperature,
        timeout=timeout,
        max_tokens=max_tokens,
    )
    if len(response.content) > 0:
        return response.content[0].text
    else:
        return ""
