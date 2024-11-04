package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

// Global variable to hold all connected clients
var clients = make(map[*net.Conn]bool)
var mu sync.Mutex // Mutex to protect access to the clients map

func main() {
	fmt.Print("Do you want to start a server (s) or a client (c)? ")

	// Read user input to decide mode
	var mode string
	_, err := fmt.Scan(&mode)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	switch mode {
	case "s":
		go StartTcp(8080) // Start the TCP server
		select {}         // Keep the main goroutine alive
	case "c":
		client := JoinTcp() // Start the TCP client
		if client != nil {
			handleClient(client) // Handle client communication
		}
	default:
		fmt.Println("Invalid option. Please choose 's' for server or 'c' for client.")
	}
}

func StartTcp(port int) {
	localIp, err := GetLocalIp()
	if err != nil {
		fmt.Println("Error getting local IP:", err)
		return
	}

	formattedAddress := fmt.Sprintf("%s:%d", localIp, port)
	fmt.Println("Starting server on", formattedAddress)

	tcpAddress, err := net.ResolveTCPAddr("tcp", formattedAddress)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close() // Ensure listener is closed when done

	// Accept incoming connections in a loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue // Continue to accept other connections
		}

		fmt.Println("Client connected:", conn.RemoteAddr().String())

		mu.Lock()
		clients[&conn] = true // Add the new client to the map
		mu.Unlock()

		// Handle the connection in a separate goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close() // Ensure connection is closed when done

	// Create a reader for the connection
	reader := bufio.NewReader(conn)
	for {
		// Read data from the client
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from client:", err)
			break
		}
		fmt.Print("Received from client:", message)

		// Broadcast the message to all connected clients
		broadcastMessage(message, conn)
	}

	// Remove the client from the list when done
	mu.Lock()
	delete(clients, &conn)
	mu.Unlock()
}

func broadcastMessage(message string, sender net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for client := range clients {
		if client != &sender { // Don't send the message back to the sender
			_, err := (*client).Write([]byte(message)) // Send message to client
			if err != nil {
				fmt.Println("Error sending message to client:", err)
				(*client).Close()       // Close the client connection on error
				delete(clients, client) // Remove client from map
			}
		}
	}
}

func JoinTcp() *net.TCPConn {
	var address string
	fmt.Print("Enter a TCP address (e.g., 192.168.1.10:8080): ")

	// Read the input from the user
	_, err := fmt.Scan(&address)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return nil
	}

	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return nil
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return nil
	}

	return conn
}

// Handle reading from and writing to the client
func handleClient(conn *net.TCPConn) {
	defer conn.Close() // Ensure connection is closed when done

	// Create a writer for the connection
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(os.Stdin) // Use standard input for user messages

	// Read input from the user and send it to the server
	go func() {
		for {
			fmt.Print("Enter message to send (type 'exit' to quit): ")
			input, _ := reader.ReadString('\n')
			if input == "exit\n" {
				break
			}

			// Send the message to the server
			_, err := writer.WriteString(input)
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
			writer.Flush() // Flush the buffer to ensure data is sent
		}
	}()

	// Read responses from the server
	serverReader := bufio.NewReader(conn)
	for {
		response, err := serverReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading response from server:", err)
			return
		}
		fmt.Print("Received from server:", response)
	}
}

// GetLocalIp retrieves the non-loopback IP address of the current device
func GetLocalIp() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no non-loopback IPv4 address found")
}
