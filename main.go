package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Something goes wrong on your platform:(")
	}
}

func showHelp() {
	println("1: show   - show list of todo")
	println("2: clear  - clear list of todo")
	println("2: add    - add new todo")
	println("2: exit   - exit of app")
}

func clearTodo() {
	if err := os.Truncate("todo.txt", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	} else {
		fmt.Println("todo list was deleted")
	}
}

func showTodo() {
	file, err := os.Open("todo.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	index := 0

	for scanner.Scan() {
		index++
		fmt.Printf("%d: %s\n", index, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func add(text string) {
	f, err := os.OpenFile("todo.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	} else {
		fmt.Println("new todo was added")
	}
}

func refreshScreen() {
	CallClear()
	fmt.Println(">>> TODO")
}

func checkFile() {
	f, err := os.OpenFile("todo.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	f.Close()
}

func main() {
	checkFile()
	refreshScreen()
	defer CallClear()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("command: ")
		// wait for user command
		command, _ := reader.ReadString('\n')

		command = strings.TrimSuffix(command, "\r\n")

		if command == "exit" {
			break
		}

		switch command {
		case "help":
			refreshScreen()
			showHelp()
			fmt.Scanln()
			refreshScreen()
			break
		case "show":
			refreshScreen()
			showTodo()
			fmt.Scanln()
			refreshScreen()
			break
		case "clear":
			refreshScreen()
			clearTodo()
			fmt.Scanln()
			refreshScreen()
			break
		case "add":
			refreshScreen()
			fmt.Print("new todo: ")
			newTodo, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			add(newTodo)
			fmt.Scanln()
			refreshScreen()
			break
		default:
			refreshScreen()
			break
		}
	}
}
