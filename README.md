# Jubo Todo List

This is a simple Restful API following Model-View-Controller Architecture (MVC), using Golang and MongoDB with some packages.
1. go-chi
2. MongoDB driver

## To run this api:
1. Download MongoDB and Golang
2. Run mongod
3. Clone this project
```
git clone https://github.com/pochiehwang-1997/Jubo_Todo_List.git
```
4. Get all dependencies
```
go get ./...
```
5. Run project
```
go run main.go
```



## To use this todo list:
1. Open a browser

- View home page: [GET] localhost:9000
- View all todos: [GET] localhost:9000/todos
- View one specific todo: [GET] localhost:9000/todos/{id}   *id could be found in all todos
- Create one todo: [POST] localhost:9000/todos body:{Title(required),description}
- Update one todo: [PUT] localhost:9000/todos/{id} body:{Title(required),description,completed}
- Delete one todo: [DELETE] localhost:9000/todos/{id}


