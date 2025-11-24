# tl (tasks list)

*tl* is simple pet-project for managing tasks with CLI client.
It uses *postgresql*, *gorm*, *echo* and *goose* for migrations.

## Run with docker compose

```shell
git clone https://github.com/visualect/tl
cd tl
cp .env.example .env 
docker compose up -d && go build -o tl cmd/cli/main.go && ./tl -help
```

## Usage

```
  -add
        add task
  -delete int
        delete item
  -list
        list all tasks
  -login string
        log in into account
  -logout
        log out from an account
  -me
        print current logged in user
  -signup string
        sign up an account
  -toggle int
        toggle completion toggle
```

## Example

 ```shell
 #  sign up first
 ./tl -signup your_login
 enter password:
 confirm password:
 user your_login successfully signed up
 
 #  then login
 ./tl -login your_login
 login success
 
 ./tl -me
 you are logged in as your_login
 
 #  add task
 ./tl -add try this application
 task 'try this application' added
 ./tl -add make this task complete
 task 'make this task complete' added
 
 #  list tasks
 ./tl -list
 
 [ ]    1. try this application
 [ ]    2. make this task complete
 
 #  toggle complete
 ./tl -toggle 2
 toggle success
 ./tl -list
 
 [ ]    1. try this application
 [X]    2. make this task complete
 
 
 #  delete task
 ./tl -delete 2
 task deleted
 ./tl -list
 
 [ ]    1. try this application
 
 ```

## Authorizaton

After logging in, a JSON file appears in root app directory (by default *.auth.json*)

## Clean up

Stop docker and remove containers:
`docker compose down`

Delete images:
`docker rmi tl-backend postgres`

Delete volumes:
`docker volume rm tl_tasks `

 
 




