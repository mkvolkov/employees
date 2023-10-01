# employees

Сервер, получающий данные в формате JSON или XML и помещающий
данные в MySQL.

Для работы программы необходим установленный и запущенный MySQL,
в котором создана база данных employees и таблица employees
с соответствующими полями.

Примеры XML хранятся в файле req.xml

Сервер хранит данные двух пользователей user1 и admin2. Первый
имеет права только на чтение (GET), второй - на всё (GET и POST).

Запуск сервера:

```
go run main.go
```

Примеры curl:

Нанять сотрудника (добавить сотрудника в БД):

```
curl -X POST http://localhost:8080/hire -i -u admin2:adminpass4
     -H 'Content-Type: application/json'
     -d '{"name":"Ivan Ivanov","phone":"89179878787","gender":"male","age":37,"email":"iivanov@gmail.com","address":"Moscow Main Street 2"}'
```

Соответствующий пример XML:

```
curl -X POST http://localhost:8080/hire -i -u admin2:adminpass4
     -H 'Content-Type: application/json'
     -d '<?xml version="1.0" encoding="UTF-8"?>
            <data>
                <employee>
                    <name>Grigory Grigoriev</name>
                    <phone>89174512733</phone>
                    <gender>male</gender>
                    <age>32</age>
                    <email>ggrigoriev@gmail.com</email>
                    <address>Moscow 17 Second Street</address>
                </employee>
            </data>'
```

Уволить сотрудника (удалить сотрудника из БД по индексу:)

```
curl -X POST http://localhost:8080/fire -i -u admin2:adminpass4
     -H 'Content-Type: application/json'
     -d '{"id":18}'
```

Уволить сотрудника (пример XML):

```
curl -X POST http://localhost:8080/fire -i -u admin2:adminpass4
     -H 'Content-Type: application/json'
     -d '<?xml version="1.0" encoding="UTF-8"?>
            <data>
                <empl_id>
                    <id>16</id>
                </empl_id>
            </data>
```

Получить количество дней отпуска по индексу:

```
curl -X GET http://localhost:8080/getv -i -u user1:userpass3
     -H 'Content-Type: application/json'
     -d '{"id":5}'
```

Найти сотрудников по имени (частичное совпадение):

```
curl -X GET http://localhost:8080/find -i -u user1:userpass3
     -H 'Content-Type: application/json'
     -d '{"name":"Mikhail"}'
```