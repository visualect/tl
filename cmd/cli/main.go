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
	"golang.org/x/term"
)

var (
	authFilename = ".auth.json"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	loginPtr := flag.String("login", "", "log in into account")
	signUpPtr := flag.String("signup", "", "sign up an account")

	addPtr := flag.Bool("add", false, "add task")
	listPtr := flag.Bool("list", false, "list all tasks")
	completePtr := flag.Int("tc", 0, "toggle completion toggle")
	// deletePtr := flag.Int("delete", 0, "delete item")
	flag.Parse()

	switch {
	case len(*loginPtr) > 0:
		if _, ok := client.IsFileExists(authFilename); ok {
			fmt.Println("you are already logged in")
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

		err = client.SaveFile(authFilename, data)
		if err != nil {
			log.Fatal("failed to save token to local file")
		}
		fmt.Println("login success")
	case len(*signUpPtr) > 0:
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
	case *addPtr:
		newTask := strings.Join(flag.Args(), " ")
		err := client.AddTask(newTask)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("task added")
	case *listPtr:
		// var list []models.Task
		fmt.Println("list")
	case *completePtr > 0:

	default:
		log.Fatal("invalid argument")
	}
}
