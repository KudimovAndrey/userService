# UserService
___

## Краткое описание
___
Небольшой HTTP-сервис, который принимает входящие соединения с JSON-данными и обрабатывает их.

Прокси слушает http://localhost:80 и перенаправляет запросы на сервер http://localhost:3400 в случаях отказа сервера или ожидания более чем 30 секунд запрос будет перенаправлен на резервный сервер http://localhost:3500

В качестве хранилища информации выступает postgreSQL. Для этого необходимо создать файл "linkFromDB.txt" с настройкой подключения к базе данных в корне проекта. Пример файла:
```sh
postgres://postgres:{password}@localhost:5432/{name_db}
```

## Пример

Запускаем main.go файл с флагом "-host localhost:3400", а также в качестве резервного сервера запускаем "-host localhost:3500"

Terminal(1)
```sh
go run main.go -host localhost:3400
```
Terminal(2)
```sh
go run main.go -host localhost:3500
```

Для создания пользователя необходима отправить post-запрос с указанием полей: имя, возраст и массив друзей.
В качестве ответа получаем ID созданого пользователя.
```sh
Request:
  curl -i -X POST -H "Content-Type:application/json" -d "{\"name\":\"andrey\",\"age\":22,\"friends\":[]}" http://localhost:80/
Response:
  User was created
  user_id:50
```

Для того чтобы сделать друзей из двух пользователей необходимо отправить post-запрос с указанием ID пользователя, который инициировал запрос на дружбу, и ID пользователя, который примет инициатора в друзья.
```sh
Request:
  curl -i -X POST -H "Content-Type:application/json" -d "{\"sourceID\":50,\"targetID\":51}" http://localhost:80/makeFriends
Response:
  andrey и Nikita теперь друзья
```

Для удаления пользователя необходимо отправить delete-запрос с указанием ID пользователя, которого нужно удалить.
```sh
Request:
  curl -i -X DELETE -H "Content-Type:application/json" -d "{\"TargetID\":52}" http://localhost:80/
Response:
  A User with the name was deleted:Artem
  User_id:52
```
Для получения всех друзей пользователя, необходимо отправить get-запрос c указанием  ID  пользователя.
```sh
Request:
  curl http://localhost:80/51
Response:
  User_id:50
  Name:andrey
  Age:22
  Friends:[51]
```
Для обновления возраста пользователя необходимо отправить put-запрос с указанием ID пользователя и нового возраста.
```sh
Request:
  curl -i -X PUT -H "Content-Type:application/json" -d "{\"age\":3}" http://localhost:80/51
Response:
  User's age has been successfully updated
```