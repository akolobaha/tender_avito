## Тестовое задание в Авито

### Сборка
```
    docker build -t tender .
```

### Запуск
```
    docker run -p 8080:8080 tender
```

### ping
```
    curl --location 'localhost:8080/api/ping'
```