# Envoy sidecar

This sidecar will be used to provide `basic-auth`, `stats`, and additional `prometheus` stats for the ollama API endpoint.

## Current State

- Getting a `upstream timeout` when make calls to `/api/generate|chat`. Yet `/api/pull` works when using port-forward

Pull in a model
```
curl --user ollama:ollama http://127.0.0.1:8089/api/show -d '{
  "model": "gemma2:2b"
}'
```


Send chat - getting timeout within seconds. (investigating this issue) 
```
 curl --user ollama:ollama -X POST http://localhost:8089/api/chat -d '{
  "model": "gemma2:2b",
  "messages": [
    {
      "role": "user",
      "content": "why is the sky blue?"
    }
  ],
  "stream": true
}'
upstream request timeout
```

Send copy 
```
curl --user ollama:ollama http://127.0.0.1:8089/api/copy -d '{
  "source": "gemma2:2b",
  "destination": "gemma2:2b-backup"
}'

```

list models using openai compat.
```
 curl -s --user ollama:ollama http://127.0.0.1:8089/v1/models | jq .
{
  "object": "list",
  "data": [
    {
      "id": "gemma2:2b-backup",
      "object": "model",
      "created": 1733680957,
      "owned_by": "library"
    },
    {
      "id": "gemma2:2b",
      "object": "model",
      "created": 1733680343,
      "owned_by": "library"
    }
  ]
}

```

list models using ollama api/tags

```
curl -s --user ollama:ollama http://127.0.0.1:8089/api/tags | jq .
{
  "models": [
    {
      "name": "gemma2:2b-backup",
      "model": "gemma2:2b-backup",
      "modified_at": "2024-12-08T18:02:37.627859521Z",
      "size": 1629518495,
      "digest": "8ccf136fdd5298f3ffe2d69862750ea7fb56555fa4d5b18c04e3fa4d82ee09d7",
      "details": {
        "parent_model": "",
        "format": "gguf",
        "family": "gemma2",
        "families": [
          "gemma2"
        ],
        "parameter_size": "2.6B",
        "quantization_level": "Q4_0"
      }
    },
    {
      "name": "gemma2:2b",
      "model": "gemma2:2b",
      "modified_at": "2024-12-08T17:52:23.188259846Z",
      "size": 1629518495,
      "digest": "8ccf136fdd5298f3ffe2d69862750ea7fb56555fa4d5b18c04e3fa4d82ee09d7",
      "details": {
        "parent_model": "",
        "format": "gguf",
        "family": "gemma2",
        "families": [
          "gemma2"
        ],
        "parameter_size": "2.6B",
        "quantization_level": "Q4_0"
      }
    }
  ]
}
```

show info about model
```
curl --user ollama:ollama http://127.0.0.1:8089/api/show -d '{
  "model": "gemma2:2b"
}'
...
```