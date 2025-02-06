from litellm import Router
from litellm.utils import token_counter, get_model_info
import json
import litellm
import requests

litellm.vertex_location = "us-east5"
litellm.vertex_project = "dashwave"

# litellm.set_verbose = True

AZURE_LOCATION = "eastus"
AZURE_VERSION = "2024-08-01-preview"

def get_model_list(api_keys:dict):
    return [
    {
        "model_name": "gpt4",
        "litellm_params": {
            "model": "gpt-4",
            "api_key": api_keys.get("openai", "")
        }
    },
    {
        "model_name": "o1",
        "litellm_params": {
            "model": "o1",
            "api_key": api_keys.get("openai", "")
        }
    },
    {
        "model_name": "o1-preview",
        "litellm_params": {
            "model": "o1-preview",
            "api_key": api_keys.get("openai", "")
        }
    },
    {
        "model_name": "o1-mini",
        "litellm_params": {
            "model": "o1-mini",
            "api_key": api_keys.get("openai", "")
        }
    },
    {
        "model_name": "gpt-4o",
        "litellm_params": {
            "model": "azure/gpt-4o",
            "api_key": api_keys.get("azure", ""),
            "api_base": api_keys.get("azure_endpoint", ""),
            "api_version": AZURE_VERSION
        }
    },
    {
        "model_name": "gpt-4o",
        "litellm_params": {
            "model": "gpt-4o",
            "api_key": api_keys.get("openai", "")
        }
    },
    {
        "model_name": "claude-3-5-sonnet",
        "litellm_params": {
            "model": "vertex_ai/claude-3-5-sonnet-v2@20241022",
            "vertex_credentials": json.dumps(api_keys.get("google_service_account_creds", {})),
        }
    },
    {
        "model_name": "claude-3-5-sonnet",
        "litellm_params": {
            "model": "claude-3-5-sonnet-20240620",
            "api_key": api_keys.get("anthropic", "")
        }
    },
    ]

router = Router(
    routing_strategy="latency-based-routing",
    routing_strategy_args={
        "ttl": 10,
        "lowest_latency_buffer": 0.5
    },
    enable_pre_call_checks=True,
    redis_host="redis",
    redis_port=6379,
    redis_password="mysecretpassword",
    cache_responses=True,
    cooldown_time=3600
)

def chat_completion(api_keys:dict, model_name:str, messages:list, temperature:float, timeout:int, max_completion_tokens:int, response_format:dict, tools:list[dict]):
    router.set_model_list(get_model_list(api_keys))

    available_models = router.get_model_list()

    selected_model_name = None
    for model in available_models:
        if model["model_name"] == model_name:
            selected_model_name = model["litellm_params"]["model"]
            break

    if selected_model_name is None:
        raise Exception(f"Model {model_name} not found")
    
    # remove the provider name from the model name if it exists
    selected_model_name = selected_model_name.split("/")[-1]
    max_allowed_input_tokens = get_model_info(selected_model_name)["max_input_tokens"]

    while True:
        messages_tokens = token_counter(
            model=selected_model_name,
            messages=messages,
        )
        if messages_tokens > max_allowed_input_tokens:
            user_msg_indices = [i for i, message in enumerate(messages) if message["role"] == "user"]
            
            # remove all the messages from top until the second user message
            if len(user_msg_indices) > 1:
                messages = messages[user_msg_indices[1]:]
        else:
            break

    response = router.completion(
        model=model_name,
        messages=messages,
        temperature=temperature,
        timeout=timeout,
        max_completion_tokens=max_completion_tokens if max_completion_tokens else None,
        response_format=response_format if response_format else None,
        tools=tools if tools else None
    )

    return response.model_dump_json()
