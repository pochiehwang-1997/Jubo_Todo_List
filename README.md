# Jubo Todo List

This is a simple Restful API following Model-View-Controller Architecture (MVC), using Golang and MongoDB with some packages.
1. go-chi
2. MongoDB driver

## To run this api:
- Run by Docker
```
docker-compose up --build
```

- Run by terminal
1. Download MongoDB and Golang
2. Run mongod
3. Clone this project
```
git clone https://github.com/pochiehwang-1997/Jubo_Todo_List.git
```
4. Get all dependencies
```
go mod download
```
5. Build
```
go build
```
5. Run project
```
go run main.go
```



## To use this todo list:
1. Open a browser

- View home page: [GET] localhost:8080
- View all todos: [GET] localhost:8080/todos
- View one specific todo: [GET] localhost:8080/todos/{id}   *id could be found in all todos
- Create one todo: [POST] localhost:8080/todos body:{Title(required),description}
- Update one todo: [PUT] localhost:8080/todos/{id} body:{Title(required),description,completed}
- Delete one todo: [DELETE] localhost:8080/todos/{id}


