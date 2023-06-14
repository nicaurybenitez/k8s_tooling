package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func execCommand(command string) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}

func readInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return input
}

func main() {
	// Solicitar los datos de conexión al usuario
	namespace := readInput("Namespace: ")
	podName := readInput("Pod Name: ")
	dbUser := readInput("Database User: ")
	dbPassword := readInput("Database Password: ")

	if len(namespace) == 0 || len(podName) == 0 || len(dbUser) == 0 || len(dbPassword) == 0 {
		fmt.Println("Please provide all required inputs.")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./pg_cli [command]")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "databases":
		execCommand(fmt.Sprintf("kubectl exec -it -n %s %s -- psql -U %s -c \"\\l\"", namespace, podName, dbUser))
	case "create_db":
		if len(args) < 1 {
			fmt.Println("Usage: ./pg_cli create_db [db_name]")
			os.Exit(1)
		}
		execCommand(fmt.Sprintf("kubectl exec -it -n %s %s -- psql -U %s -c \"CREATE DATABASE %s\"", namespace, podName, dbUser, args[0]))
	case "drop_db":
		if len(args) < 1 {
			fmt.Println("Usage: ./pg_cli drop_db [db_name]")
			os.Exit(1)
		}
		execCommand(fmt.Sprintf("kubectl exec -it -n %s %s -- psql -U %s -c \"DROP DATABASE %s\"", namespace, podName, dbUser, args[0]))
	case "tables":
		if len(args) < 1 {
			fmt.Println("Usage: ./pg_cli tables [db_name]")
			os.Exit(1)
		}
		execCommand(fmt.Sprintf("kubectl exec -it -n %s %s -- psql -U %s -c \"\\dt %s.*\"", namespace, podName, dbUser, args[0]))
	case "table_info":
		if len(args) < 2 {
			fmt.Println("Usage: ./pg_cli table_info [db_name] [table_name]")
			os.Exit(1)
		}
		execCommand(fmt.Sprintf("kubectl exec -it -n %s %s -- psql -U %s -c \"\\d %s.%s\"", namespace, podName, dbUser, args[0], args[1]))
	case "select":
		if len(args) < 2 {
			fmt.Println("Usage: ./pg_cli select [db_name] [table_name]")
			os.Exit(1)
		}
		execCommand(fmt.Sprintf("kubectl exec -it -n %s %s -- psql -U %s -c \"SELECT * FROM %s.%s\"", namespace, podName, dbUser, args[0], args[1]))
	case "insert":
		if len(args) < 3 {
			fmt.Println("Usage: ./pg_cli insert [db_name] [table_name] [values]")
			os.Exit(1)
		}
		execCommand(fmt.Sprintf("kubectl exec -it -n %s %s -- psql -U %s -c \"INSERT INTO %s.%s VALUES %s\"", namespace, podName, dbUser, args[0], args[1], args[2]))
	case "update":
		if len(args) < 4 {
			fmt.Println("Usage: ./pg_cli update [db_name] [table_name] [set_clause] [condition]")
			os.Exit(1)
		}
		execCommand(fmt.Sprintf("kubectl exec -it -n %s %s -- psql -U %s -c \"UPDATE %s.%s SET %s WHERE %s\"", namespace, podName, dbUser, args[0], args[1], args[2], args[3]))
	case "delete":
		if len(args) < 3 {
			fmt.Println("Usage: ./pg_cli delete [db_name] [table_name] [condition]")
			os.Exit(1)
		}
		execCommand(fmt.Sprintf("kubectl exec -it -n %s %s -- psql -U %s -c \"DELETE FROM %s.%s WHERE %s\"", namespace, podName, dbUser, args[0], args[1], args[2]))
	default:
		fmt.Println("Command not found.")
		os.Exit(1)
	}
}

