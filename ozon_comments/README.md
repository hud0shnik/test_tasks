# Тестовое задание для Ozon

Реализовать систему для добавления и чтения постов и комментариев с использованием GraphQL, аналогичную комментариям к постам на популярных платформах, таких как Хабр или Reddit.

Характеристики системы постов:
- Можно просмотреть список постов.
- Можно просмотреть пост и комментарии под ним.
- Пользователь, написавший пост, может запретить оставление комментариев к своему посту.

Характеристики системы комментариев к постам:
- Комментарии организованы иерархически, позволяя вложенность без ограничений.
- Длина текста комментария ограничена до, например, 2000 символов.
- Система пагинации для получения списка комментариев.

(*) Дополнительные требования для реализации через GraphQL Subscriptions:
Комментарии к постам должны доставляться асинхронно, т.е. клиенты, подписанные на определенный пост, должны получать уведомления о новых комментариях без необходимости повторного запроса.

Требования к реализации:
- Система должна быть написана на языке Go.
- Использование Docker для распространения сервиса в виде Docker-образа.
- Хранение данных может быть как в памяти (in-memory), так и в PostgreSQL. Выбор хранилища должен быть определяемым параметром при запуске сервиса.
- Покрытие реализованного функционала unit-тестами.

Критерии оценки:
- Как хранятся комментарии и как организована таблица в базе данных/in-memory, включая механизм пагинации.
- Качество и чистота кода, структура проекта и распределение файлов по пакетам.
- Обработка ошибок в различных сценариях использования.
- Удобство и логичность использования системы комментариев.
- Эффективность работы системы при множественном одновременном использовании, сравнимая с популярными сервисами, такими как Хабр.
- В реализации учитываются возможные проблемы с производительностью, такие как проблемы с n+1 запросами и большая вложенность комментариев.

# Как запускать
Локально:
```
go run cmd/main.go
```
Через докер:
```
docker compose --env-file ./.env up -d
```
# Где хранить данные
Для выбора места для хранения постов и комментариев используется переменная окружения STORAGE:
1. STORAGE="in-memory" - сервис будет хранить данные в памяти
2. STORAGE="db" - сервис будет хранить данные в БД (PostgreSQL)
Переменные окружения можно записать в файл .env в корне проекта

# API
Для обмена информацией сервис использует GraphQL. Файлы спецификаций находятся в папке **graph/**. 
API реализовывает следующий функционал:
### Мутация CreatePost(post: PostInput!):Post!
Для создания постов.
```graphql
mutation{
	CreatePost(post:{
    author: "admin",
    header:"head",
    content:"content"
    commentsAllowed: true
  }){
    id
    author
    header
    content
    commentsAllowed
    createdAt
  }
}
```
### Запрос GetAllPosts(page: Int, pageSize: Int): [Post!]!
Для получения всех постов. Используется пагинация для более удобного вывода (page - номер страницы, pageSize - размер страниц).
```graphql
query {
  GetAllPosts(page:2, pageSize:3) {
    id
    author
    header
    content
    createdAt  
  }
}
```
### Запрос GetPostById(id: Int): PostOutput!
Для получения поста по айди.
```graphql
query {
  GetPostById(id:15) {
    id
    author
    header
    content
    createdAt
    commentsAllowed
  }
}
```
### Мутация CreateComment(input: CommentIntput!): Comment!
Для создания комментария к посту.
```graphql
mutation{
  CreateComment(input:{
    post: 8,
    author: "admin",
    content: "content of comment"
  }){
    id
    post
    author
    content
    createdAt
  }
}
```
### Запрос GetAllCommentsByPost(id: Int, page: Int, pageSize: Int): [Comment!]!
Для получения всех комментариев к посту по айди. Используется пагинация для более удобного вывода (page - номер страницы, pageSize - размер страниц).
```graphql
query {
  GetAllCommentsByPost(id:15, page:1, pageSize:3) {
    id
    author
    content
    createdAt  
  }
}
```
### Запрос GetAllRepliesByComment(id: Int): [Comment!]!
Для получения всех ответов на комментарий по айди.
```graphql
query {
  GetAllRepliesByComment(id:3){
    id
    post
    author
    content
    createdAt
    replyTo
  }
}
```
