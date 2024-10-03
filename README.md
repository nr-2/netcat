# Net-Cat 🐱

## Description
Net-Cat is a project that recreates the functionality of NetCat in a Server-Client Architecture. It allows communication between a server and multiple clients over TCP connections, enabling chat functionality similar to group messaging.

## Features
- TCP connection between server and multiple clients (1-to-many relationship).
- Clients are required to provide a name upon connection.
- Control over the maximum number of connections (up to 10 clients).
- Clients can send messages to the chat.
- Empty messages from clients are ignored.
- Messages are timestamped and identified by the sender's name in the format: `[timestamp][client.name]: [client.message]`.
- New clients joining receive all previous chat messages.
- Notification to all clients when a new client joins or leaves the chat.
- Clients receive messages sent by other clients.
- Disconnection of one client does not affect the others.
- Supports default port 8989 if none is specified.

## Installation and Usage


### 1. **Run the server:**
``` 
go run .
 ```

This will start the server listening on the default port 8989.

### 2. **Run the client:**
   ``` 
   nc localhost 8989 
   ```

Replace localhost and 8989 with the appropriate IP address and port number if running on a different machine or port.

### 3. **Usage:**
Upon connecting, you will see a Linux ASCII art logo and be prompted to enter your name.
Enter your desired username.
Start chatting by typing messages. Messages will be broadcasted to all connected clients.
Type `/name [new_name] `to change your username.
Use Ctrl+C to terminate the client or server.

## Requirements
Go programming language environment must be installed.



