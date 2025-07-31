package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	userStoreFilePath = "users.txt"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type Task struct {
	ID, UserID, CategoryID int
	IsDone                 bool
	Title, DueDate         string
}

type Category struct {
	ID, UserID int
	Title      string
	Color      string
}

var userStorage []User
var taskStorage []Task
var categoryStorage []Category
var authenticatedUser *User

func main() {

	loadUserStorage()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("welcome to todo app")
	command := flag.String("command", "no-command", "command to run")

	flag.Parse()

	for {
		runCommand(*command)
		fmt.Println("waiting for another command...")
		scanner.Scan()
		*command = scanner.Text()
	}

}

func loadUserStorage() {
	file, oErr := os.Open(userStoreFilePath)
	if oErr != nil {
		fmt.Println("cant open file", oErr)

		return
	}
	b := make([]byte, 1000)
	_, rErr := file.Read(b)
	if rErr != nil {
		fmt.Println("cant read file", rErr)

		return
	}
	strData := string(b)
	records := strings.Split(strData, "\n")
	for i, r := range records {
		if i == len(records)-1 {
			continue
		}
		fmt.Println(i, r)
	}

}
func runCommand(command string) {
	if command != "register-user" && command != "exit" && authenticatedUser == nil {
		fmt.Println("You must login first")
		loginUser()
		if authenticatedUser == nil {
			fmt.Println("email or password is incorrect!!")

			return
		}
	}
	switch command {
	case "create-task":
		createTask()
	case "list-task":
		listTask()
	case "create-category":
		createCategory()
	case "register-user":
		registerUser()
	case "login-user":
		loginUser()
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("command is not valid", command)
	}
}

func listTask() {
	for _, task := range taskStorage {
		if task.UserID == authenticatedUser.ID {
			fmt.Println(task)
		}
	}
}
func createTask() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("please enter the task title")
	scanner.Scan()
	title := scanner.Text()

	fmt.Println("please enter the task dueDate")
	scanner.Scan()
	dueDate := scanner.Text()

	fmt.Println("please enter the task category id")
	scanner.Scan()
	categoryID, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("category id is not a valid integer:", err)

		return
	}

	isFound := false
	for _, c := range categoryStorage {
		if c.ID == categoryID && c.UserID == authenticatedUser.ID {
			isFound = true

			break
		}
	}
	if !isFound {
		fmt.Println("category not found:")

		return
	}

	task := Task{
		ID:         len(taskStorage) + 1,
		Title:      title,
		CategoryID: categoryID,
		DueDate:    dueDate,
		IsDone:     false,
		UserID:     authenticatedUser.ID,
	}

	taskStorage = append(taskStorage, task)
	fmt.Printf("task created: %+v\n", task)
}
func createCategory() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, color string

	fmt.Println("please enter the category title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the category color")
	scanner.Scan()
	color = scanner.Text()

	category := Category{
		ID:     len(categoryStorage) + 1,
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}
	categoryStorage = append(categoryStorage, category)
	fmt.Println("category created:", category)
}

func registerUser() {
	scanner := bufio.NewScanner(os.Stdin)
	var email, name, password string

	fmt.Println("please enter user name")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("please enter user email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter user password")
	scanner.Scan()
	password = scanner.Text()

	user := User{
		ID:       len(userStorage) + 1,
		Name:     name,
		Email:    email,
		Password: password,
	}
	userStorage = append(userStorage, user)

	file, oErr := os.OpenFile(userStoreFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if oErr != nil {
		fmt.Println("an internal error occurred:", oErr)

		return
	}
	defer file.Close()

	data := fmt.Sprintf("%+v\n", user)

	_, wErr := file.WriteString(data)
	if wErr != nil {
		fmt.Println("an internal error occurred", wErr)

		return
	}

	fmt.Printf("user created: %+v\n", user)
}
func loginUser() {
	var email, password string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("please enter your email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter your password")
	scanner.Scan()
	password = scanner.Text()

	for _, user := range userStorage {
		if email == user.Email && password == user.Password {
			authenticatedUser = &user
			fmt.Println("you are logged in!")

			break
		}
	}

}
