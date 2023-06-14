package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"encoding/base64"
	

	"github.com/awnumar/memguard"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	envFile         = ".env"
	encryptedEnvKey = "your-encryption-key" // Cambia esto con tu propia clave de cifrado
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
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	input = strings.TrimSpace(input)
	return input
}

func checkCommandAvailability(command string) {
	_, err := exec.LookPath(command)
	if err != nil {
		fmt.Printf("Command '%s' not found.\n", command)
		os.Exit(1)
	}
}

func encryptEnvData(data []byte, key *[32]byte) ([]byte, error) {
	defer memguard.WipeBytes(data)

	nonce := new([24]byte)
	_, err := rand.Read(nonce[:])
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	encrypted := secretbox.Seal(nonce[:], data, nonce, key)
	return encrypted, nil
}

func decryptEnvData(encrypted []byte, key *[32]byte) ([]byte, error) {
	nonce := new([24]byte)
	copy(nonce[:], encrypted[:24])

	decrypted, ok := secretbox.Open(nil, encrypted[24:], nonce, key)
	if !ok {
		return nil, fmt.Errorf("failed to decrypt data")
	}

	return decrypted, nil
}

func writeEncryptedEnvFile(namespace, podName, dbUser, dbPassword string, key *[32]byte) {
	content := fmt.Sprintf(`NAMESPACE=%s
POD_NAME=%s
DB_USER=%s
DB_PASSWORD=%s`, namespace, podName, dbUser, dbPassword)

	encrypted, err := encryptEnvData([]byte(content), key)
	if err != nil {
		fmt.Println("Error encrypting .env file:", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(envFile, encrypted, 0644)
	if err != nil {
		fmt.Println("Error writing .env file:", err)
		os.Exit(1)
	}
}

func readEncryptedEnvFile(key *[32]byte) (string, string, string, string) {
	encrypted, err := ioutil.ReadFile(envFile)
	if err != nil {
		fmt.Println("Error reading .env file:", err)
		os.Exit(1)
	}

	decrypted, err := decryptEnvData(encrypted, key)
	if err != nil {
		fmt.Println("Error decrypting .env file:", err)
		os.Exit(1)
	}

	lines := strings.Split(string(decrypted), "\n")
	var namespace, podName, dbUser, dbPassword string

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "NAMESPACE":
				namespace = value
			case "POD_NAME":
				podName = value
			case "DB_USER":
				dbUser = value
			case "DB_PASSWORD":
				dbPassword = value
			}
		}
	}

	return namespace, podName, dbUser, dbPassword
}

func generateEncryptionKey() *[32]byte {
	key := new([32]byte)
	_, err := rand.Read(key[:])
	if err != nil {
		fmt.Println("Error generating encryption key:", err)
		panic(err)
	}
	return key
}

func writeEncryptionKeyToFile(key *[32]byte, keyFile string) {
	keyBase64 := base64.StdEncoding.EncodeToString(key[:])
	err := ioutil.WriteFile(keyFile, []byte(keyBase64), 0600)
	if err != nil {
		fmt.Println("Error writing encryption key file:", err)
		os.Exit(1)
	}
}

func readEncryptionKeyFromFile(keyFile string) *[32]byte {
	keyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		fmt.Println("Error reading encryption key file:", err)
		os.Exit(1)
	}

	keyBase64 := strings.TrimSpace(string(keyData))
	keySlice, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		fmt.Println("Error parsing encryption key:", err)
		os.Exit(1)
	}

	var key [32]byte
	copy(key[:], keySlice)

	return &key
}

func printHelp() {
	fmt.Println(`Usage: ./pg_cli [options] [command]

Options:
  -h, --help   Show help information
  config       Configure database connection
  generate_key Generate and save encryption key

Commands:
  databases                List all databases
  create_db [db_name]      Create a new database
  drop_db [db_name]        Drop an existing database
  tables [db_name]         List all tables in a database
  table_info [db_name] [table_name]
                           Show information about a table
  select [db_name] [table_name]
                           Select all records from a table
  insert [db_name] [table_name] [values]
                           Insert a new record into a table
  update [db_name] [table_name] [set_clause] [condition]
                           Update records in a table
  delete [db_name] [table_name] [condition]
                           Delete records from a table`)
}

func main() {
	checkCommandAvailability("kubectl")

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	arg := os.Args[1]
	if arg == "-h" || arg == "--help" {
		printHelp()
		os.Exit(0)
	}

	if arg == "config" {
		namespace := readInput("Namespace: ")
		podName := readInput("Pod Name: ")
		dbUser := readInput("Database User: ")
		dbPassword := readInput("Database Password: ")

		key := readEncryptionKeyFromFile(".encryption_key")
		writeEncryptedEnvFile(namespace, podName, dbUser, dbPassword, key)

		fmt.Println("Configuration saved.")
		os.Exit(0)
	}

	if arg == "generate_key" {
		key := generateEncryptionKey()
		writeEncryptionKeyToFile(key, ".encryption_key")

		fmt.Println("Encryption key generated and saved.")
		os.Exit(0)
	}

	key := readEncryptionKeyFromFile(".encryption_key")

	namespace, podName, dbUser, dbPassword := readEncryptedEnvFile(key)

	_ = dbPassword // Evita la advertencia de compilación de variable no utilizada

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
