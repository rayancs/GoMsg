package main

// import (
// 	"bufio"
// 	"fmt"
// 	"net"
// 	"sync"
// )

// var clients = make(map[*net.Conn]bool)
// var mu sync.Mutex // Mutex to protect access to the clients map

// func main() {
// 	server := StartTcp(8080)

// 	client := JoinTcp()

// }
// func JoinTcp() *net.TCPConn {
// 	var address string
// 	fmt.Print("Enter a TCP address (e.g., 192.168.1.10:8080): ")

// 	// Read the input from the user
// 	_, err := fmt.Scan(&address) // Reads a single word (up to whitespace)
// 	if err != nil {
// 		fmt.Println("Error reading input:", err)
// 		return nil
// 	}
// 	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
// 	if err != nil {
// 		fmt.Println("Error reading input:", err)
// 		return nil
// 	}
// 	conn, err := net.DialTCP("tcp", nil, tcpAddress)
// 	if err != nil {
// 		fmt.Println("Error reading input:", err)
// 		return nil
// 	}
// 	return conn
// }
// func StartTcp(port int) *net.Conn {
// 	localIp, err := GetLocalIp()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return nil
// 	}
// 	formatedAddress := fmt.Sprintf("%s:%d", localIp, port)
// 	fmt.Println(formatedAddress)
// 	tcpAddress, err := net.ResolveTCPAddr("tcp", formatedAddress)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return nil
// 	}
// 	listner, err := net.ListenTCP("tcp", tcpAddress)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return nil
// 	}
// 	conn, err := listner.Accept()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return nil
// 	}
// 	return &conn

// }
// func GetLocalIp() (string, error) {
// 	// will return all network interfaces
// 	networkInterfaces, err := net.Interfaces()
// 	if err != nil {
// 		return "", err
// 	}
// 	// loop therough the list of interfaces

// 	for _, Inet := range networkInterfaces {
// 		/*
// 			Check The Following Falgs
// 			btw Go's net pacakge provides you with flags
// 			check if the interface u are looping at i is not null
// 			** using go net package **
// 			@flag
// 				net.FlagUp
// 			@desc
// 				to check if interface / ethernet / wifi etc is active or not
// 			@usage
// 				net.FlagUp == 1 // active
// 			@flag
// 				net.FlagLoopback
// 			@desc
// 				to check if its a internal network

// 		*/
// 		// check inactive and internals
// 		if Inet.Flags&net.FlagUp == 0 || Inet.Flags&net.FlagLoopback != 0 {
// 			continue
// 		}
// 		// if not check the addresses of the interface
// 		addresses, err := Inet.Addrs()
// 		if err != nil {
// 			return "", err
// 		}
// 		// parse usinga swith in golang
// 		for _, addr := range addresses {
// 			var ip net.IP
// 			switch v := addr.(type) {
// 			case *net.IPNet:
// 				/*
// 					@desc
// 						ip with subnet
// 				*/
// 				ip = v.IP
// 			case *net.IPAddr:
// 				/*
// 					@desc
// 						ip without  subnet
// 				*/
// 				ip = v.IP
// 			}
// 			// validate var ip
// 			if ip != nil && ip.To4() != nil && !ip.IsLoopback() {
// 				return ip.String(), nil
// 			}

// 		}
// 	}
// 	return "", fmt.Errorf("no non-loopback IPv4 address found")

// }
// func handleConnection(conn net.Conn) {
// 	defer conn.Close() // Ensure connection is closed when done

// 	// Create a reader for the connection
// 	reader := bufio.NewReader(conn)
// 	for {
// 		// Read data from the client
// 		message, err := reader.ReadString('\n')
// 		if err != nil {
// 			fmt.Println("Error reading from client:", err)
// 			break
// 		}
// 		fmt.Print("Received from client:", message)

// 		// Broadcast the message to all connected clients
// 		broadcastMessage(message, conn)
// 	}

// 	// Remove the client from the list when done
// 	mu.Lock()
// 	delete(clients, &conn)
// 	mu.Unlock()
// }

// // Broadcast a message to all clients
// func broadcastMessage(message string, sender net.Conn) {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	for client := range clients {
// 		if client != &sender { // Don't send the message back to the sender
// 			_, err := (*client).Write([]byte(message)) // Send message to client
// 			if err != nil {
// 				fmt.Println("Error sending message to client:", err)
// 				(*client).Close()       // Close the client connection on error
// 				delete(clients, client) // Remove client from map
// 			}
// 		}
// 	}
// }
