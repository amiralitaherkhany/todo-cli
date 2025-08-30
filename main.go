package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type Task struct {
	ID         int
	Title      string
	DueDate    string
	CategoryID int
	IsDone     bool
	UserID     int
}

type Category struct {
	ID     int
	Title  string
	Color  string
	UserID int
}

var (
	userStorage     []User
	categoryStorage []Category
	taskStorage     []Task

	authenticatedUser *User
	serializationMode string
)

const (
	userStoragePath               = "user.txt"
	ManDarAvardiSerializationMode = "mandaravardi"
	JsonSerializationMode         = "json"
)

func main() {

	serializeMode := flag.String("serialize-mode", ManDarAvardiSerializationMode, "serialization mode to write data to file")
	command := flag.String("command", "no-command", "command to run")
	flag.Parse()

	// load user storage from file
	loadUserStorageFromFile(*serializeMode)

	fmt.Println("Hello to TODO app")

	switch *serializeMode {
	case ManDarAvardiSerializationMode:
		serializationMode = ManDarAvardiSerializationMode
	default:
		serializationMode = JsonSerializationMode
	}
	// if there is a user record with corresponding data allow the user to continue

	for {
		runCommand(*command)

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("please enter another command")
		scanner.Scan()
		*command = scanner.Text()
	}
}

func runCommand(command string) {
	if command != "register-user" && command != "exit" && authenticatedUser == nil {
		login()

		if authenticatedUser == nil {
			return
		}
	}

	switch command {
	case "create-task":
		createTask()
	case "create-category":
		createCategory()
	case "register-user":
		registerUser()
	case "list-task":
		listTask()
	case "login":
		login()
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("command is not valid", command)
	}
}

func createTask() {

	scanner := bufio.NewScanner(os.Stdin)
	var title, duedate, category string

	fmt.Println("please enter the task title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the task category id")
	scanner.Scan()
	category = scanner.Text()

	categoryID, err := strconv.Atoi(category)
	if err != nil {
		fmt.Printf("category-id is not valid integer, %v\n", err)

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
		fmt.Printf("category-id is not found\n")

		return
	}

	fmt.Println("please enter the task due date")
	scanner.Scan()
	duedate = scanner.Text()

	// validation
	// category validate

	task := Task{
		ID:         len(taskStorage) + 1,
		Title:      title,
		DueDate:    duedate,
		CategoryID: categoryID,
		IsDone:     false,
		UserID:     authenticatedUser.ID,
	}

	taskStorage = append(taskStorage, task)
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
	fmt.Println("category", title, color)

	c := Category{
		ID:     len(categoryStorage) + 1,
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}

	categoryStorage = append(categoryStorage, c)
}

func registerUser() {
	scanner := bufio.NewScanner(os.Stdin)
	var name, email, password string
	var id int
	fmt.Println("please enter the name")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("please enter the email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter the password")
	scanner.Scan()
	password = scanner.Text()

	id = len(userStorage) + 1

	user := User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: hashThePassword(password),
	}

	fmt.Println("user:", id, email, password)

	userStorage = append(userStorage, user)

	writeUserToFile(user)
}

func hashThePassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func login() {
	fmt.Println("login process")
	scanner := bufio.NewScanner(os.Stdin)
	var email, password string

	fmt.Println("please enter email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter the password")
	scanner.Scan()
	password = scanner.Text()

	for _, user := range userStorage {
		if user.Email == email && user.Password == hashThePassword(password) {
			authenticatedUser = &user
			break
		}
	}

	if authenticatedUser == nil {
		fmt.Println("the email or password is not correct")
	}
}

func listTask() {
	for _, task := range taskStorage {
		if task.UserID == authenticatedUser.ID {
			fmt.Println(task)
		}
	}
}

func loadUserStorageFromFile(serializationMode string) {
	file, err := os.Open(userStoragePath)
	if err != nil {
		fmt.Println("can't open the file", err)
	}

	var data = make([]byte, 1024)
	_, oErr := file.Read(data)
	if oErr != nil {
		fmt.Println("can't read from the file", oErr)

		return
	}

	var dataStr = string(data)

	userSlice := strings.Split(dataStr, "\n")
	fmt.Println("len userSlice", len(userSlice), serializationMode)
	for _, u := range userSlice {
		var userStruct = User{}

		switch serializationMode {
		case ManDarAvardiSerializationMode:
			var dErr error
			userStruct, dErr = deserilizeFromManDaravardi(u)
			if dErr != nil {
				fmt.Println("can't deserialize user record to user struct", dErr)

				return
			}
		case JsonSerializationMode:
			if u[0] != '{' && u[len(u)-1] != '}' {
				continue
			}

			uErr := json.Unmarshal([]byte(u), &userStruct)
			if uErr != nil {
				fmt.Println("can't deserialize user record to user struct with json mode", uErr)

				return
			}
		default:
			fmt.Println("invalid serialization mode")

			return
		}

		userStorage = append(userStorage, userStruct)
	}
}

func writeUserToFile(user User) {
	var file *os.File

	file, err := os.OpenFile(userStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("can't create or open file", err)

		return
	}
	defer file.Close()

	var data []byte
	// serialize the user struct/object
	if serializationMode == ManDarAvardiSerializationMode {
		data = []byte(fmt.Sprintf("id: %d, name: %s, email: %s, password: %s\n", user.ID, user.Name,
			user.Email, user.Password))
	} else if serializationMode == JsonSerializationMode {
		//json

		var jErr error
		data, jErr = json.Marshal(user)
		if jErr != nil {
			fmt.Println("can't marshal user struct to json", jErr)

			return
		}

		data = append(data, []byte("\n")...)

	} else {
		fmt.Println("invalid serialization mode")

		return
	}

	numberOfWrittenBytes, wErr := file.Write(data)
	if wErr != nil {
		fmt.Printf("can't write to the file %v\n", wErr)

		return
	}

	fmt.Println("numberOfWrittenBytes", numberOfWrittenBytes)
}

func deserilizeFromManDaravardi(userStr string) (User, error) {

	if userStr == "" {
		return User{}, errors.New("user string is empty")
	}

	var user = User{}

	userFields := strings.Split(userStr, ",")
	for _, field := range userFields {
		values := strings.Split(field, ": ")
		if len(values) != 2 {
			fmt.Println("field is not valid, skipping...", len(values))

			continue
		}
		fieldName := strings.ReplaceAll(values[0], " ", "")
		fieldValue := values[1]

		switch fieldName {
		case "id":
			id, err := strconv.Atoi(fieldValue)
			if err != nil {
				return User{}, errors.New("strconv error")
			}
			user.ID = id
		case "name":
			user.Name = fieldValue
		case "email":
			user.Email = fieldValue
		case "password":
			user.Password = fieldValue
		}

	}

	return user, nil
}
