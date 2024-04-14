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
  request: |
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


## Comparison params
