## Why CompextAI?

Most AI companies rely on almost the same publicly available models. Yet, what makes each of them so different from other are mainly 2 things:

1. **Input prompts** - both system and user.
2. **Context** - Providing the model with the correct information.

While the iteration speeds for testing input prompts have become quicker by deploying the right tools in place, **managing context is still a problem**.

To better understand and improve our model’s performance, its critical to have an observability over the information being provided to the model. It is critical that the model has enough context to evaluate the inputs correctly, while ensuring that we don’t flood the model with too much information that it starts to hallucinate (and you end up paying too much money 😕)

## What is CompextAI?

CompextAI is an LLMOps tool designed for AI Pipelines where developers and product-managers can conveniently monitor their LLM conversations and executions. Not only can they monitor, but they can perform powerful actions on their existing conversations.

Some actions include:

1. Summarising the previous conversation messages while providing a new prompt (new message in the conversation) to reduce input prompt token length.
2. Execute a conversation (previously run on a different model in your AI pipeline) on a different model with different parameters to compare your outputs.
3. Editing messages in an existing conversation to evaluate for better feedback results.

You can maintain your AI conversations similar to how OpenAI Threads are maintained. Except with CompextAI, they are extremely powerful and compatible across all models. No more managing the conversations in a local/global array variable 😮‍💨

## What else can CompextAI do?

CompextAI has one goal - To let developers and product managers both be a part of the AI Pipeline development and review cycle. And to make this whole process a lot quicker.

**It can do a lot of developer friendly stuff, for example:**

```python
import compextAI.api.api as compextAPI
import compextAI.execution as compextExecution
import compextAI.threads as compextThreads
import compextAI.params as compextParams
import compextAI.messages as compextMessages
# import the libraries from compext

# initialise the client
client = compextAPI.APIClient(
    base_url="https://api.compextai.dev",
    api_key="xxxxxxxxxxxxxxx"
)

# create a thread - start a conversation
thread: compextThreads.Thread = compextThreads.create(
    client=client,
    project_name="demo-project",
    title="Say hello in docs!",
    metadata={
        "date": "2024-11-11"
    }
)
print("thread id: ", thread.thread_id)
```

Now, you can save this `thread_id` associated with this conversation. Later you can use this to add/edit messages and execute this thread. 

We can add messages to this empty thread:

```python
compextMessages.create(
    client=client,
    thread_id=thread.thread_id,
    messages=[
        compextMessages.Message(
            role="system",
            content="You are a helpful assistant who speaks in Hindi."
        ),
        compextMessages.Message(
            role="user",
            content="Say hello!"
        )
    ]
)
```

Now let’s fetch the model parameters to execute this thread from a defined `Prompt Config` on the compextAI dashboard:

```python
params = compextParams.retrieve(
    client=client,
    name="demo",
    environment="development",
    project_name="demo-project"
)
```

Now, let’s finally execute it: ⛳

```python
thread_execution = thread.execute(
    client=client,
    thread_exec_param_id=params.thread_execution_param_id,
    append_assistant_response=True
)

# wait for completion
wait_for_completion(thread_execution.thread_execution_id)

thread_execution_result: dict = compextExecution.get_thread_execution_response(
    client=client,
    thread_execution_id=thread_execution.thread_execution_id
)
# the content key returns the content of the response by the assistant
print(thread_execution_result['content'])
# नमस्ते!
```

**It can also do a lot of product-manager friendly stuff, for example:**

One can view this conversation on the compextAI dashboard:

// image 

// image

You can view the details of the execution as well:

// image

The product manager can view all details of this execution and even re-execute this thread with the same inputs with different configurations (different model, model params etc.)

The `Prompt Configs` can be configured from the dashboard:

// image

The developer can go on to add more messages to this conversation/edit conversation messages, and then execute the conversation again.

You can find the entire docs [here](https://compextai.notion.site/Docs-13b5ef52981080b4bdd9dcad34bbc394)
