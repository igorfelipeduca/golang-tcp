# Golang TCP Server

This project implements a simple TCP server using Go.

## Project Structure

- `main.go`: The main file that contains the server implementation.
- `go.mod`: The module file that contains the dependencies.

## How to Run

1. Clone the repository
2. Run `go mod tidy` to download the dependencies
3. Run `go run main.go` to start the server

## How to Use

1. Install telnet (macOS: `brew install telnet`)
2. Run `telnet localhost 8080`
3. Type a message and press enter
4. The server will respond with the message "Hey %s. I got your message: %s"

## License

This project is open-sourced under the MIT License - see the LICENSE file for details.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
