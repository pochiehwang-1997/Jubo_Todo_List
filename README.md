# Jubo Todo List

This is a simple Restful API using Golang and MongoDB with some packages.
1. go-chi
2. mongoDB driver

## To run this api:


## To use this todo list:
1. Open a browser

- View home page: [GET] localhost:9000
- View all todos: [GET] localhost:9000/todos
- View one specific todo: [GET] localhost:9000/todos/{id}   *id could be found in all todos
- Create one todo: [POST] localhost:9000/todos body:{Title(required),description}
- Update one todo: [PUT] localhost:9000/todos/{id} body:{Title(required),description,completed}
- Delete one todo: [DELETE] localhost:9000/todos/{id}


