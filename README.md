# webapp-tutorial-backend

## API

- POST user/login
- POST user/refresh

- GET todos
- GET todos/{id}
- POST todos/new
- PUT todos/{id}
- DELETE todos/{id}

### register

- POST user/register

```
$ curl -v --request POST -H "Content-Type: application/json" -d '{"email": "hatsune@miku.com", "password": "PASSWORD"}' http://localhost:8080/user/register

< HTTP/1.1 201 Created
< Content-Type: application/json; charset=UTF-8
< Set-Cookie: refresh_token=REFRESH_TOKEN_IS_HERE
< Vary: Accept-Encoding
< Date: Fri, 30 Apr 2021 05:16:03 GMT
< Content-Length: 262
<
{"email":"hatsune@miku.com","access_token":"ACCESS_TOKE_IS_HERE","created_at":"2021-04-30T14:16:03.807055192+09:00"}
```

### login

```
$ curl -v --request POST -H "Content-Type: application/json" -d '{"email": "hatsune@miku.com", "password": "PASSWORD"}' http://localhost:8080/user/login

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Set-Cookie: refresh_token=REFRESH_TOKEN_IS_HERE
< Vary: Accept-Encoding
< Date: Fri, 30 Apr 2021 05:20:11 GMT
< Content-Length: 184
<
{"access_token":"ACCESS_TOKE_IS_HERE"}
```

### token refresh

```
$ curl -v --request POST -b "refresh_token=$REFRESH_TOKEN" http://localhost:8080/user/refresh

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Set-Cookie: refresh_token=REFRESH_TOKEN_IS_HERE
< Vary: Accept-Encoding
< Date: Fri, 30 Apr 2021 05:16:53 GMT
< Content-Length: 184
<
{"access_token":"ACCESS_TOKE_IS_HERE"}
```

### create TODO

```
$ curl -v --request POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"content": "go home"}' http://localhost:8080/todos

< HTTP/1.1 201 Created
< Content-Type: application/json; charset=UTF-8
< Vary: Accept-Encoding
< Date: Fri, 30 Apr 2021 05:21:04 GMT
< Content-Length: 202
<
{"id":"f233e9a1-01c0-4e43-aca9-089076f21a5d","content":"go home","completed":false,"created_at":"2021-04-30T14:21:04.055762286+09:00","updated_at":"2021-04-30T14:21:04.055762286+09:00","deleted":false}
```


### get all TODO

```
$ curl -v -H "Authorization: Bearer $TOKEN" http://localhost:8080/todos

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Vary: Accept-Encoding
< Date: Fri, 30 Apr 2021 05:22:38 GMT
< Content-Length: 174
<
[{"id":"f233e9a1-01c0-4e43-aca9-089076f21a5d","content":"go home","completed":false,"created_at":"2021-04-30T05:21:04Z","updated_at":"2021-04-30T05:21:04Z","deleted":false}]
```

### get a TODO

```
$ curl -v -H "Authorization: Bearer $TOKEN" http://localhost:8080/todos/f233e9a1-01c0-4e43-aca9-089076f21a5d

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Vary: Accept-Encoding
< Date: Fri, 30 Apr 2021 05:23:21 GMT
< Content-Length: 172
<
{"id":"f233e9a1-01c0-4e43-aca9-089076f21a5d","content":"go home","completed":false,"created_at":"2021-04-30T05:21:04Z","updated_at":"2021-04-30T05:21:04Z","deleted":false}
```

### update TODO

```
$ curl -v --request PUT -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"content": "go home!!", "completed": true, "deleted": false}' http://localhost:8080/todos/f233e9a1-01c0-4e43-aca9-089076f21a5d

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Vary: Accept-Encoding
< Date: Fri, 30 Apr 2021 05:24:42 GMT
< Content-Length: 173
<
{"id":"f233e9a1-01c0-4e43-aca9-089076f21a5d","content":"go home!!","completed":true,"created_at":"2021-04-30T05:21:04Z","updated_at":"2021-04-30T05:21:04Z","deleted":false}
```

### delete TODO

```
$ curl -v --request DELETE -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" http://localhost:8080/todos/f233e9a1-01c0-4e43-aca9-089076f21a5d

< HTTP/1.1 200 OK
< Vary: Accept-Encoding
< Date: Fri, 30 Apr 2021 05:26:18 GMT
< Content-Length: 0
<
```