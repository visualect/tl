package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/visualect/tl/internal/client"
	"github.com/visualect/tl/internal/dto"
	"github.com/visualect/tl/internal/models"
	"golang.org/x/term"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	loginPtr := flag.String("login", "", "log in into account")
	signUpPtr := flag.String("signup", "", "sign up an account")
	logOutPtr := flag.Bool("logout", false, "log out from an account")

	addPtr := flag.Bool("add", false, "add task")
	listPtr := flag.Bool("list", false, "list all tasks")
	completePtr := flag.String("toggle", "", "toggle completion toggle")
	deletePtr := flag.String("delete", "", "delete item")

	flag.Parse()

	switch {
	case len(*loginPtr) > 0:
		if _, ok := client.IsFileExists(client.AuthFilename); ok {
			fmt.Println("you are currently logged in")
			return
		}
		login := *loginPtr
		fmt.Println("enter password:")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal(err)
		}

		data, err := client.Login(login, string(password))
		if err != nil {
			log.Fatal(err)
		}

		err = client.SaveFile(client.AuthFilename, data)
		if err != nil {
			log.Fatal("failed to save token to local file")
		}
		fmt.Println("login success")
	case len(*signUpPtr) > 0:
		if _, ok := client.IsFileExists(client.AuthFilename); ok {
			fmt.Println("you are currently logged in")
			return
		}
		login := *signUpPtr
		fmt.Println("enter password:")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("confirm password:")
		passwordConfirm, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal(err)
		}

		if string(password) != string(passwordConfirm) {
			log.Fatal("passwords don't match")
		}

		l, err := client.SignUp(login, string(password))
		if err != nil {
			log.Fatal(err)
		}

		var resJSON dto.RegisterResponse
		err = json.NewDecoder(bytes.NewReader(l)).Decode(&resJSON)
		fmt.Printf("user %s successfully signed up\n", resJSON.Login)
	case *logOutPtr:
		if _, ok := client.IsFileExists(client.AuthFilename); !ok {
			fmt.Println("you are already logged out")
			return
		}
		err := client.DeleteFile(client.AuthFilename)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("you are logged out")
	case *addPtr:
		if _, ok := client.IsFileExists(client.AuthFilename); !ok {
			fmt.Println("please, log in first")
			return
		}
		newTask := strings.Join(flag.Args(), " ")
		data, err := client.AddTask(newTask)
		if err != nil {
			log.Fatal(err)
		}
		var t dto.AddTaskRequest
		err = json.NewDecoder(bytes.NewReader(data)).Decode(&t)
		fmt.Printf("task '%s' added\n", t.Task)
	case *listPtr:
		if _, ok := client.IsFileExists(client.AuthFilename); !ok {
			fmt.Println("please, log in first")
			return
		}

		data, err := client.GetTasks()
		if err != nil {
			log.Fatal(err)
		}

		var list []models.Task
		err = json.NewDecoder(bytes.NewReader(data)).Decode(&list)
		if err != nil {
			log.Fatal(err)
		}

		if len(list) == 0 {
			fmt.Println("you list is empty")
		}

		c := map[bool]string{
			true:  "[x]",
			false: "[ ]",
		}

		// TODO: change to show index + 1
		for _, task := range list {
			fmt.Printf("%s\t%d. %s\n", c[task.Completed], task.ID, task.Task)
		}
	case len(*completePtr) > 0:
		if _, ok := client.IsFileExists(client.AuthFilename); !ok {
			fmt.Println("please, log in first")
			return
		}

		_, err := client.ToggleCompleteTask(*completePtr)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("toggle complete success")
	case len(*deletePtr) > 0:
		if _, ok := client.IsFileExists(client.AuthFilename); !ok {
			fmt.Println("please, log in first")
			return
		}

		_, err := client.DeleteTask(*deletePtr)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("task deleted")
	default:
		log.Fatal("invalid argument")
	}
}
