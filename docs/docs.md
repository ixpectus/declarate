# Tests 

## Commands

### Request 

#### Light mode example
```yaml
- name: check request in light mode
  method: GET
  path: /loadName
  responseStatus: 200
  response: |
      {"age": 28,"name":"name1", "items":[1, 2, 3]}
```

#### Extended mode example
```yaml
- name: check request
  method: GET
  path: /loadName
  fullResponse: |
    {
      "body":{"age": 28,"name":"{{$name}}", "items":[1, 2, 3, 4]} 
      "status": 200
    }
```
Extended mode is useful when it is neccesary save status to variable

### Database 
#### Example
```yaml
- name: check values in table
  db_conn: '{{$full_connection}}'
  db_query: select id from t1
  db_response: >
    [{"id":"q"}]
```
##### Properties
- `db_conn` database connection string, if it not set default database connection string will be used.

### Script 
#### Example
```yaml
- name: script check
  script: 
    path: "./tests/scripts/echo.sh"
```

### Shell 
#### Example
```yaml
- name: shell command 2
  shell_cmd: | 
    ls -la | head -n1
  variables:
    name: '*'
```
Run shell command and save it output to name variable


## Variables

### Set variables
#### Examples
```yaml
- name: set variables
  variables:
    name: bor
```

### Set variables from commands response
Setting variables from responses is the same for different commands.  
It is possible set variables from commands responses, using [gjson](github.com/tidwall/gjson) 
#### Set variables from request command
```yaml
- name: setup variables
  variables: |
    db2: "dbname"
- name: db2 for user SHOULD be creating
  path: /clusters/{{$CLUSTER_ID}}/databases
  method: POST
  request: |
    {
      "name": "{{$db2}}"
    }
  response: |
    {
      "name":"{{$db2}}",
      "status": "CREATING"
    }
  variables:
    db2_id: id
```

#### Set variables from database command example
```yaml
  db_query: |
    select 1 as a;
  db_response: '[{"a":1}]'
```

## Compare response

### Comparison params
- `allowArrayExtraItems` allow extra items on array comparison
- `ignoreArraysOrdering` ignore array ordering on array comparison
#### Examples 
##### Request
```yaml
- name: check request 2
  method: GET
  path: /{{$path}}
  fullResponse: |
    {
      "body": {"age": 28,"name":"{{$name}}", "items":[1, 2, 3]}, 
      "status": 200
    }
  comparisonParams:
    ignoreArraysOrdering: true
    allowArrayExtraItems: true
```
The test will pass successfully with the following responses.

- Using `allowArrayExtraItems` parameter
```json
{
  "body": {"age": 28,"name":"{{$name}}", "items":[1, 2, 3, 4, 5, 6]}, 
  "status": 200
}
```
- Using `ignoreArraysOrdering` parameter
```json
{
  "body": {"age": 28,"name":"{{$name}}", "items":[3, 2, 1]}, 
  "status": 200
}
```

### Modificators

#### Regexp
When comparing responses, you can use regular expressions with the `$matchRegexp` modifier.
##### Examples
- master can be any 
```json
    [
      {
        "name":"{{$db1}}",
        "status": "OK",
        "connections": {
          "master": "$matchRegexp(^.+$)"
        }
      }
    ]
```
- version should start from 15
```json
 {"version": "$matchRegexp(^15.+$)"} 
```
#### Custom modificators
- `any` 
  Example 
  ```json
    {
      "body":{"name":"$any", "items":[1, 2, 3, 4]}, 
    }
   ```
- `notEmpty`  
  Example 
  ```json
    {
      "body":{"age": "$num","name":"$any", "items":[1, 2, 3, 4]}, 
      "status": "$oneOf(300, 200)"
    }
   ```
- `num`  
  Example 
  ```json
    {
      "body":{"age": "$num","items":[1, 2, 3, 4]}, 
    }
   ```
- `oneOf`  
  Example 
  ```json
    {
      "status": "$oneOf(300, 200)"
    }
   ```

## Tests flow

### Test steps
Test steps can be combined using the `steps` directive. 
#### Example
```yaml
- name: database list SHOULD be empty for user
  steps:
    - name: database list SHOULD be empty for user in api
      path: /clusters/{{$CLUSTER_ID}}/users/{{$user_id}}/databases
      method: GET
      response: |
        []
      responseStatus: 200
    - name: database list SHOULD be empty for user in db
      db_query: |
        select count(*) from pg_catalog.pg_user where usename='{{$user}}'
      db_response: '[{"count":1}]'
  poll:
    duration: 14s
```

### Test definition
#### Tags
Every test can contain test definition with tags. It is possible to filter tests by tags.

##### Tags example
```yaml
- definition:
    tags:
        - base
        - agent-base
```

### Conditions
Test steps or the entire test can be skipped according to conditions.

#### Definition example
```yaml
- definition:
    tags: ["base"]
    condition: 'notEmpty("{{$HOME}}")'
```
#### Test step example
```yaml
- name: step should be skipped
  condition: 'notEmpty("{{$HOME}}")'
  variables:
    password: qwerty
```


### Polling

#### Example with response regexp
Within 100 seconds, a request will be sent to the 'poll' handler every second. 
The polling will continue until the response matches the regular expression in `respone_regexp` field.
```yaml
- name: check poll handler with regexp success
  method: GET
  path: /poll
  response: |
    {"age": 31,"name":"Tommy"}
  responseStatus: 200
  poll:
    response_regexp: ".+Tom.+"
    duration: 100s
    interval: 1s
```

#### Example with response object
Within 100 seconds, a request will be sent to the 'poll' handler every second. 
The polling will continue until the response matches response object.
```yaml
- name: check poll handler with regexp success
  method: GET
  path: /poll
  response: |
    {"age": 31,"name":"Tommy"}
  responseStatus: 200
  poll:
    response: | 
      {"name":"Tom"} 
    duration: 100s
    interval: 1s
```

#### Poll properties
- `duration` 
- `interval`
- `response_regexp`
- `response`
