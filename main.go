package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Client struct {
	Name   string
	Conn   net.Conn
	Joined time.Time
}

type Message struct {
	Timestamp time.Time
	Name      string
	Text      string
}

var (
	clients  = make(map[*Client]bool)
	mutex    = sync.Mutex{}
	messages []Message
	logFile  *os.File
)

func init() {
	var err error
	logFile, err = os.OpenFile("chat.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
}

func main() {
	port := "8989"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	fmt.Printf("Listening on the port :%s\n", port)
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer ln.Close()

	for {
		if len(clients) >= 10 {
			log.Printf("Server is Full")
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("Failed to accept connection: %v", err)
			} else {
				conn.Close()
			}
			continue
		}

		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	client := &Client{
		Conn:   conn,
		Joined: time.Now(),
	}

	// Ask for the client's name
	fmt.Fprintf(conn, "Welcome to TCP-Chat!\n%s\n[ENTER YOUR NAME]: ", getAsciiLogo())
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("Failed to read client name: %v", err)
		return
	}
	client.Name = strings.TrimSpace(name)

	// Validate the client's name
	for isUsernameTaken(client.Name) {
		fmt.Fprintf(conn, "Username %s is already taken. Please choose another.\n[ENTER ANOTHER NAME]: ", client.Name)
		name, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("Failed to read client name: %v", err)
			return
		}
		client.Name = strings.TrimSpace(name)
	}

	// Add the client to the list of active clients
	mutex.Lock()
	if len(clients) >= 10 {
		log.Printf("Maximum number of clients (10) reached, rejecting new connection")
		conn.Close()
		mutex.Unlock()
		return
	}
	clients[client] = true
	mutex.Unlock()

	// Notify other clients about the new connection
	broadcastMessage(Message{
		Timestamp: time.Now(),
		Name:      client.Name,
		Text:      "has joined our chat...",
	})

	// Send previous messages to the new client
	for _, msg := range messages {
		if msg.Name != client.Name || msg.Text != "has joined our chat..." {
			sendMessageToClient(conn, msg)
		}
	}

	// Handle client messages
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		msg = strings.TrimSpace(msg)

		if msg == "" {
			continue // Skip empty messages
		}

		if strings.HasPrefix(msg, "/name ") {
			newName := strings.TrimSpace(strings.TrimPrefix(msg, "/name "))
			if len(newName) < 3 || len(newName) > 20 {
				fmt.Fprintf(conn, "Invalid username length. Must be between 3 and 20 characters.\n")
				continue
			}
			if isUsernameTaken(newName) {
				fmt.Fprintf(conn, "Username %s is already taken. Please choose another.\n", newName)
				continue
			}
			oldName := client.Name
			client.Name = newName
			broadcastMessage(Message{
				Timestamp: time.Now(),
				Name:      oldName,
				Text:      "changed their name to " + newName,
			})
			continue
		}

		broadcastMessage(Message{
			Timestamp: time.Now(),
			Name:      client.Name,
			Text:      msg,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from connection: %v\n", err)
	}

	mutex.Lock()
	delete(clients, client)
	mutex.Unlock()

	broadcastMessage(Message{
		Timestamp: time.Now(),
		Name:      client.Name,
		Text:      "has left our chat...",
	})
}

func broadcastMessage(msg Message) {
	mutex.Lock()
	defer mutex.Unlock()

	// Avoid duplicate messages by checking the last message
	if len(messages) > 0 && messages[len(messages)-1].Name == msg.Name && messages[len(messages)-1].Text == msg.Text {
		return
	}

	messages = append(messages, msg)
	logMessage(msg) // Save the message to the log file
	for client := range clients {
		sendMessageToClient(client.Conn, msg)
	}
}

func sendMessageToClient(conn net.Conn, msg Message) {
	if strings.Contains(msg.Text, "joined our chat") || strings.Contains(msg.Text, "left our chat") {
		fmt.Fprintf(conn, "%s %s\n", msg.Name, msg.Text)
	} else {
		fmt.Fprintf(conn, "[%s][%s]: %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"), msg.Name, msg.Text)
	}
}

func logMessage(msg Message) {
	log.Printf("[%s][%s]: %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"), msg.Name, msg.Text)
}

func isUsernameTaken(name string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	for client := range clients {
		if client.Name == name {
			return true
		}
	}
	return false
}

func getAsciiLogo() string {
	return "\n" +
		"         _nnnn_\n" +
		"        dGGGGMMb\n" +
		"       @p~qp~~qMb\n" +
		"       M|@||@) M|\n" +
		"       @,----.JM|\n" +
		"      JS^\\__/  qKL\n" +
		"     dZP        qKRb\n" +
		"    dZP          qKKb\n" +
		"   fZP            SMMb\n" +
		"   HZM            MMMM\n" +
		"   FqM            MMMM\n" +
		" __| \".        |\\dS\"qML\n" +
		" |    `.       | `' \\Zq\n" +
		"_)      \\.___.,|     .'\n" +
		"\\____   )MMMMMP|   .'\n" +
		"     `-'       `--'\n" +
		""
}
