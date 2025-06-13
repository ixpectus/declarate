# Функциональные требования к библиотеке Declarate

## Общее описание

**Declarate** - это библиотека для декларативного тестирования API, CLI и других систем с использованием простого синтаксиса YAML. Основная идея - описывать тесты в понятном и лаконичном виде, фокусируясь на расширяемости и гибкости.

**Вдохновлена**: [gonkey](https://github.com/lamoda/gonkey)

## Основные принципы

### 1. Декларативность
- Тесты описываются в YAML-файлах с простым и понятным синтаксисом
- Минимальное количество кода для описания сложных сценариев
- Читаемость и понятность тестов для нетехнических специалистов

### 2. Модульность и расширяемость
- Архитектура на основе команд (commands) с единым интерфейсом
- Возможность легкого добавления новых типов команд
- Интерфейс `contract.Doer` для всех исполняемых действий
- Интерфейс `contract.CommandBuilder` для создания команд

### 3. Гибкость конфигурации
- Поддержка двух форматов конфигурации: краткого и расширенного
- Возможность использования переменных и их подстановки
- Поддержка условных выражений и циклов

## Архитектурные требования

### Интерфейс Doer
Все команды должны реализовывать интерфейс:
```go
type Doer interface {
    Do() error                    // Выполнение команды
    ResponseBody() *string        // Получение тела ответа
    IsValid() error              // Валидация конфигурации
    GetConfig() interface{}      // Получение конфигурации
    Check() error               // Проверка результата
    SetVars(vv Vars)           // Установка переменных
    SetReport(r ReportAttachement) // Установка отчета
}
```

### Интерфейс CommandBuilder
Для создания команд из YAML конфигурации:
```go
type CommandBuilder interface {
    Build(unmarshal func(interface{}) error) (Doer, error)
}
```

## Функциональные требования

### 1. Команды (Commands)

#### 1.1 HTTP Request Command
**Назначение**: Выполнение HTTP-запросов и сравнение ответов

**Краткий формат**:
```yaml
- name: check request
  method: GET
  path: /user/{{$user_id}}
  response: '{"name": "{{$name}}"}'
  responseStatus: 200
```

**Расширенный формат**:
```yaml
- name: check request
  request:
    method: GET
    path: /user/{{$user_id}}
    headers:
      Authorization: "Bearer {{$token}}"
    body: '{"data": "value"}'
    response: '{"name": "{{$name}}"}'
    responseStatus: 200
    comparisonParams:
      ignoreArraysOrdering: true
```

**Требования**:
- Поддержка всех HTTP методов (GET, POST, PUT, DELETE, PATCH, etc.)
- Настройка заголовков и cookies
- Проверка статуса ответа
- Сравнение тела ответа с ожидаемым
- Поддержка различных параметров сравнения

#### 1.2 Database Command
**Назначение**: Выполнение SQL-запросов и проверка результатов

**Формат**:
```yaml
- name: check database
  db_conn: "{{$connection_string}}"
  db_query: "SELECT id, name FROM users WHERE id = {{$user_id}}"
  db_response: '[{"id": 1, "name": "John"}]'
```

**Требования**:
- Поддержка PostgreSQL (и возможность расширения)
- Выполнение SELECT, INSERT, UPDATE, DELETE запросов
- Сравнение результатов с ожидаемыми значениями
- Поддержка параметризованных запросов
- Возможность установки переменных из результатов запроса

#### 1.3 Shell Command
**Назначение**: Выполнение shell-команд

**Формат**:
```yaml
- name: run shell command
  shell_cmd: "ls -la | head -n1"
  shell_response: "total 42"
```

**Расширенный формат**:
```yaml
- name: run shell command
  shell:
    cmd: "ls -la | head -n1"
    response: "total 42"
    comparisonParams:
      allowArrayExtraItems: true
```

**Требования**:
- Выполнение произвольных shell-команд
- Захват stdout команды
- Сравнение вывода с ожидаемым результатом
- Поддержка переменных в командах

#### 1.4 Script Command
**Назначение**: Выполнение скриптов

**Формат**:
```yaml
- name: run script
  script_path: "./scripts/setup.sh"
  script_response: "Setup completed"
  script_nowait: false
```

**Расширенный формат**:
```yaml
- name: run script
  script:
    path: "./scripts/setup.sh"
    response: "Setup completed"
    nowait: false
```

**Требования**:
- Выполнение скриптов по указанному пути
- Асинхронное выполнение (nowait)
- Захват вывода скрипта
- Проверка результата выполнения

#### 1.5 Variables Command
**Назначение**: Установка и управление переменными

**Формат**:
```yaml
- name: set variables
  variables:
    user_name: "John Doe"
    user_id: 123
    token: "{{.response.access_token}}"
```

**Требования**:
- Установка обычных переменных
- Установка персистентных переменных
- Извлечение значений из ответов команд (gjson)
- Применение переменных в шаблонах

#### 1.6 Echo Command
**Назначение**: Вывод отладочной информации

**Формат**:
```yaml
- name: debug output
  echo_message: "Processing user: {{$user_name}}"
  echo_response: "Processing user: John Doe"
```

**Требования**:
- Вывод сообщений с подстановкой переменных
- Проверка корректности вывода
- Использование для отладки и логирования

### 2. Система переменных

#### 2.1 Интерфейс Vars
```go
type Vars interface {
    Set(k, val string) error           // Установка переменной
    SetAll(m map[string]string) (map[string]string, error) // Установка множества переменных
    Get(k string) string               // Получение переменной
    Apply(text string) string          // Применение переменных к тексту
    SetPersistent(k, val string) error // Установка персистентной переменной
}
```

#### 2.2 Типы переменных
- **Обычные переменные**: существуют в рамках одного теста
- **Персистентные переменные**: сохраняются между запусками тестов
- **Переменные из ответов**: извлекаются из результатов команд

#### 2.3 Шаблонизация
- Синтаксис: `{{$variable_name}}`
- Поддержка функций: `{{.response.field_name}}`
- Автоматическая подстановка во всех текстовых полях

### 3. Система сравнения ответов

#### 3.1 Параметры сравнения (CompareParams)
```yaml
comparisonParams:
  ignoreArraysOrdering: true      # Игнорировать порядок в массивах
  allowArrayExtraItems: true      # Разрешить дополнительные элементы в массивах
  ignoreValues: ["id", "timestamp"] # Игнорировать указанные поля
```

#### 3.2 Поддерживаемые форматы
- JSON - основной формат для API
- Plain text - для shell команд и скриптов
- Регулярные выражения

### 4. Структура тестов

#### 4.1 Файловая организация
- Тесты организуются в YAML-файлы
- Поддержка вложенных директорий
- Возможность запуска отдельных файлов или директорий

#### 4.2 Последовательность выполнения
```yaml
- name: setup
  variables:
    user_name: "Test User"

- name: create user
  method: POST
  path: /users
  request: '{"name": "{{$user_name}}"}'
  variables:
    user_id: "id"  # Извлечь id из ответа

- name: verify user
  method: GET
  path: /users/{{$user_id}}
  response: '{"name": "{{$user_name}}"}'
```

#### 4.3 Вложенные шаги
```yaml
- name: complex scenario
  steps:
    - name: step 1
      method: GET
      path: /data
    - name: step 2
      method: POST
      path: /process
```

### 5. Продвинутые возможности

#### 5.1 Polling (опрос с интервалом)
```yaml
- name: wait for completion
  method: GET
  path: /status/{{$job_id}}
  response: '{"status": "completed"}'
  poll:
    duration: 30s
    interval: 2s
    response_body_regexp: '"status":\s*"completed"'
```

#### 5.2 Условия выполнения
```yaml
- name: conditional step
  method: GET
  path: /data
  condition: "{{$env}} == 'production'"
```

#### 5.3 Теги и фильтрация
```yaml
- name: basic test
  tags: ["basic", "smoke"]
  method: GET
  path: /health

- name: advanced test  
  tags: ["advanced", "integration"]
  method: POST
  path: /complex-operation
```

**Команды запуска**:
```bash
./declarate -dir ./tests -tags "basic"        # Только базовые тесты
./declarate -dir ./tests -tags "smoke,basic"  # Smoke и базовые тесты
./declarate -dir ./tests -tests "user,auth"   # Определенные файлы/директории
```

### 6. Интеграция и отчетность

#### 6.1 Allure Reports
- Автоматическое создание отчетов в формате Allure
- Прикрепление артефактов (запросы, ответы)
- Группировка по тестам и наборам

#### 6.2 Output система
```go
type Output interface {
    Log(message Message)
    SetReport(r Report)
}
```

#### 6.3 Test Wrapper
```go
type TestWrapper interface {
    BeforeTest(file string, conf *RunConfig, lvl int)
    AfterTest(conf *RunConfig, result Result)
    BeforeTestStep(file string, conf *RunConfig, lvl int)
    AfterTestStep(conf *RunConfig, result Result, isPolling bool)
}
```

### 7. CLI интерфейс

#### 7.1 Основные команды
```bash
./declarate -dir <path>                 # Запуск тестов из директории
./declarate -dir <path> -dryRun         # Показать список тестов без выполнения
./declarate -dir <path> -tags <tags>    # Запуск тестов с определенными тегами
./declarate -dir <path> -tests <names>  # Запуск определенных тестов
./declarate -dir <path> -s <skip>       # Пропустить определенные тесты
```

#### 7.2 Параметры конфигурации
- `-clearPersistent` - очистить персистентные переменные
- `-withProgressBar` - показать прогресс-бар
- `-defaultDBConn` - строка подключения к БД по умолчанию
- `-defaultHost` - хост по умолчанию для HTTP-запросов

### 8. Требования к производительности

- Параллельное выполнение независимых тестов
- Эффективная работа с большими объемами данных
- Минимальное потребление памяти
- Быстрый запуск и инициализация

### 9. Требования к качеству

#### 9.1 Тестируемость
- Покрытие unit-тестами всех основных компонентов
- Интеграционные тесты для проверки взаимодействия
- Моки для внешних зависимостей

#### 9.2 Документация
- Подробная документация API
- Примеры использования для каждого типа команд
- Руководство по расширению функциональности

#### 9.3 Обратная совместимость
- Поддержка миграции из gonkey
- Стабильный API между версиями
- Deprecation warnings для устаревшей функциональности

### 10. Технические требования

#### 10.1 Зависимости
- Go 1.20+
- Поддержка основных ОС (Linux, macOS, Windows)
- Минимальные внешние зависимости

#### 10.2 Безопасность
- Безопасная обработка переменных окружения
- Валидация входных данных
- Предотвращение injection-атак

#### 10.3 Конфигурация
- Поддержка переменных окружения
- Файлы конфигурации
- CLI параметры с приоритетом

Эти требования обеспечивают создание мощной, гибкой и расширяемой системы декларативного тестирования, подходящей для различных сценариев использования - от простых smoke-тестов до сложных интеграционных проверок.
