[![codecov](https://codecov.io/github/ajugalushkin/goph-keeper/branch/iter3/graph/badge.svg?token=AM7EOPJ3D6)](https://codecov.io/github/ajugalushkin/goph-keeper)
[![Lint tests](https://github.com/ajugalushkin/goph-keeper/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/ajugalushkin/goph-keeper/actions/workflows/golangci-lint.yml)
[![Unit tests](https://github.com/ajugalushkin/goph-keeper/actions/workflows/action-codecov.yml/badge.svg)](https://github.com/ajugalushkin/goph-keeper/actions/workflows/action-codecov.yml)
------
# Менеджер паролей GophKeeper

GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю безопасно хранить
логины, пароли, данные банковских карт, произвольные текстовые и бинарные данные.

Сервер поддерживает следующий функционал:
* регистрация, аутентификация и авторизация пользователей;
* хранение приватных данных пользователей;
* синхронизация данных между несколькими авторизованными клиентами одного владельца;
* передача приватных данных владельцу по запросу.

Клиент реализует следующую бизнес-логику:
* регистрация, аутентификация и авторизация пользователей на удалённом сервере;
* доступ к приватным данным по запросу.

## Настройка и запуск сервера

Перед запуском сервера необходимо создать конфигурационный файл с настройками
и задать его с помощью флага --config или через переменную окружения SERVER_CONFIG. 
Пример конфигурационного файла:

```
# server-config.yaml

env: "dev"  #dev, prod
grpc:
  address: ":8080"
  timeout: 1h
token:
  ttl: 10h
  secret: secret_key
storage:
  path: postgres://praktikum:pass@postgres:5432/goph_keeper
minio:
  endpoint: minio:9000
  username: praktikum
  password: pass1234567
  ssl: false
  bucket: vault
```

После создания файла, необходимо указать его в файле .env, в папке docker

```
# .env
SERVER_CONFIG=/config/server-config.yaml
TOKEN_SECRET=secret_key
```

После настройки запуск сервера осуществляется командной:

```
make up
```

## Настройка и запуск клиента

Перед запуском клиента необходимо создать конфигурационный файл с настройками
и задать его через флаг --config или через переменную окружения CLIENT_CONFIG. Пример конфигурационного файла:

```
# client-config.yaml

env: "dev"  #dev, prod
client:
  address: ":8080"
  timeout: 1h
  retries: 3
```

Пример настройки клиента через переменные окружения:

```
CLIENT_CONFIG=/config/client-config.yaml
```

## Процедуры регистрации, аутентификации, авторизации

При регистрации пользователя необходимо указать адрес электронной почты и пароль.
Пример команды регистрации:

```
./goph-keeper-cli auth register -e user@mail.ru -p 123456
```

В случае успешного выполнения запроса регистрации нового пользователя,
необходимо выполнить команду login для получения токена доступа:

```
./goph-keeper-cli auth login -e user@mail.ru -p 123456
```

## Хранение приватных данных пользователя

После записи полученного при регистрации токена доступа в файл token.txt
становятся доступны команды для управления приватными данными пользователя.

### Добавление приватных данных

1. Пример команды, сохраняющей данные банковской карты:

```
./goph-keeper-cli keep create card \
  --name visa \
  --number 1111222233334444 \
  --date 12/22 \
  --holder Alexandr \
  --code 512
```

2. Пример команды создания пары логин пароль:

```
./goph-keeper-cli keep create credentials \
  --name yandex-mail \
  --login user@yandex.ru \
  --password 12345678
```

3. Пример команды для сохранения текстовой информации

```
./goph-keeper-cli keep create text \
  --name pushkin \
  --data "medny vsadnik"
```

3. Пример команды для сохранения бинарных данных

```
./goph-keeper-cli keep create bin \
  --name code \
  --file_path test_video.mp4
```

### Получение данных

Для получения приватных данных необходимо указать название секрета, пример:

```
./goph-keeper-cli keep get --name visa
```

Также можно вывести список всех приватных данных пользователя:

```
./goph-keeper-cli keep list
```

### Редактирование и удаление данных

Пример редактирования данных о банковской карте:

```
./goph-keeper-cli keep update card \
  --name mastercard \
  --number 5555222277775555 \
  --date 01/36 \
  --holder "Alexandr Nevskiy" \
  --code 789
```

Пример команды удаления данных:

```
./goph-keeper-cli keep delete --name visa
```