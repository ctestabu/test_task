
```
psql -U postgres
CREATE DATABASE test_db;
```
```
psql -U postgres -d test_db -f schema_new.sql
```

```
export DATABASE_URL=postgres://postgres:password@127.0.0.1:5432/test_db?sslmode=disable
```


```
curl -d '{"login":"alice","password":"secret"}' -H "Content-Type: application/json" http://localhost:8080/api/auth
```
```
curl -X POST -H 'Authorization: Bearer <token>' -d 'Hello, Alice!' http://localhost:8080/api/upload-asset/hello
```

```
curl -X POST -H 'Authorization: Bearer <token>' -d 'Hello, Alice!' http://localhost:8080/api/asset/hello
```

```
curl -H 'Authorization: Bearer <token>' http://localhost:8080/api/list-assets
```

```curl -X DELETE -H 'Authorization: Bearer <token>' http://localhost:8080/api/delete-asset/hello4
```

### Что можно улучшить в схеме БД?

Добавить индекс для ускорения запросов.
Хранение паролей с безопасным хешированием.
Хранение файлов отдельно от БД.
Если удалить пользователя, его сессии и файлы останутся в БД.