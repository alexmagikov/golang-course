# RepoStat

Микросервисное приложение для получения информации о GitHub-репозиториях.

---

## Быстрый старт

```bash
git clone https://github.com/alexmagikov/golang-course
cd golang-course/task3
make up
```

## Использование

| Сервис             | URL                          |
|--------------------|------------------------------|
| API Gateway        | http://localhost:28080       |
| Swagger UI         | http://localhost:28080/swagger |

**Ping:**
```bash
curl http://localhost:28080/api/ping
```

**Информация о репозитории:**
```bash
curl "http://localhost:28080/api/repositories/info?url=https://github.com/golang/go"
```
