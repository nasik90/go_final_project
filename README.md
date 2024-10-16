Проект посвящен созданию планировщика. Планировщик позволяет создавать, удалять, получать задачи. Так же есть возможность поиска по названию, комментарию и дате задачи. У задачи есть дата выполнения задачи, а так же механизм повторения задач: ежегодное повторение, перенос на указанное количество дней, повторение в указанные дни недели, в указанные дни месяца и месяцы года. 

Все задания со * выполнены.

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

Пример команды для создания docker конейнера `docker build --tag my-todo-app:v1 .`, для запуска `docker run -d -p 7540:7540 my-todo-app:v1`.

Для запуска планировщика необходимо перейти по адресу http://localhost:7540/

Запуск тестов: `go test ./tests`