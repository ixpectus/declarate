# Declarate: library for declarative testing of APIs, CLIs and anything else

Declarate is a library for declarative testing of APIs, CLIs, and any other services using YAML files. Inspired by [gonkey](https://github.com/lamoda/gonkey).

- Simple, readable declarative syntax
- Easy to write tests: send a request, compare the response
- Extensible and flexible — add your own action types
- Early stage of development

## License

[MIT](LICENSE)

## Quick Start

Describe your test in YAML, specify commands, and Declarate will execute them sequentially.

```yaml
- name: set variables
  variables:
    path: users/42

- name: get user
  method: GET
  path: /api/{{$path}}
  response:
    name: Alice
    age: 30
  responseStatus: 200
```

## Commands

### Request — HTTP Requests

#### Light Mode
- `path` — request path
- `method` — HTTP method (GET, POST, PUT, etc.)
- `response` — expected response body
- `responseStatus` — expected status code

#### Extended Mode
- `fullResponse` — full response (body + status + all headers), useful for storing status in a variable

### Database
- `db_conn` — database connection string (optional; uses the default if omitted)
- `db_query` — SQL query
- `db_response` — expected result in JSON format

### Script
- `script.path` — path to an executable script

### Shell
- `shell_cmd` — shell command; its output is saved to a variable

## Variables

### Setting variables
```yaml
- name: set vars
  variables:
    name: Tom
    count: 5
```

### Extracting from responses
Uses [gjson](https://github.com/tidwall/gjson) syntax:
```yaml
- name: get user
  path: /api/users
  method: GET
  response:
    name: "{{$name}}"
  variables:
    user_id: id        # extracts the JSON field "id"
```

### Variables from request body
Variables from previous steps can be referenced via `request`:
```yaml
- name: create user
  method: POST
  path: /api/users
  request: |         # {{$db2}} — variable from the previous step
    {"name": "{{$db2}}"}
```

## Response Comparison

### Comparison parameters
- `allowArrayExtraItems` — allows extra items in array comparison
- `ignoreArraysOrdering` — ignores element order in arrays

### Modifiers
Use modifiers in expected responses for flexible comparison:

| Modifier | Description | Example |
|---|---|---|
| `$matchRegexp(...)` | Regular expression match | `{"version": "$matchRegexp(^15.+$)"}` |
| `$any` | Value must be present, content can be anything | `{"name": "$any"}` |
| `$notEmpty` | Value must not be empty | `{"name": "$notEmpty"}` |
| `$num` | Must be a numeric value | `{"age": "$num"}` |
| `$oneOf(v1, v2, ...)` | Must be one of the listed values | `{"status": "$oneOf(200, 300)"}` |

## Test Flow Control

### Steps
Group commands within a single test:
```yaml
- name: verify database
  steps:
    - name: check API
      path: /api/data
      response: []
    - name: check DB
      db_query: select count(*) from items
      db_response: '[{"count":0}]'
```

### Tags
Filter tests by tags:
```yaml
- definition:
    tags: [base, smoke]
```

### Conditions
Skip a step or an entire test when a condition is not met:
```yaml
- name: optional step
  condition: 'notEmpty("{{$HOME}}")'
```

### Polling
Send periodic requests until the expected state is reached:
```yaml
- name: wait for resource ready
  method: GET
  path: /api/resource/{{$id}}
  response: {"status": "ready"}
  poll:
    response_regexp: ".+pending.+"   # continue while this matches
    duration: 100s                   # timeout
    interval: 1s                     # polling interval
```

### Poll (Object Comparison)
Instead of a regexp, you can specify an expected response object:
```yaml
  poll:
    response: {"status": "ready"}
    duration: 100s
    interval: 1s
```
