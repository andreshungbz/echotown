# Echo Town 📢

`echotown` is a simple TCP echo server written in Go. It can handle multiple client connections, does input validation, and has a personality mode and command protocol for custom responses. Server connections/disconnections and client messages are properly logged.

## Video Demonstration

Click on this [YouTube](https://youtu.be/L2FzdG2_gZw?si=2wAJf43wLx-L_pjX) link to see a demonstration of the server in action.

## Running The Server

### Prerequisites

- [Go](https://go.dev/) installed on your system
- [Netcat](https://netcat.sourceforge.net/) utility
- [Make](https://www.gnu.org/software/make/) build tool

### Build Steps

1. Clone the repository.

```
git clone https://github.com/andreshungbz/echotown.git
```

2. Change directory to the project folder.

```
cd echotown
```

3. Build the project using `make`

```
make
```

### Server Steps

1. Start the Echo Town server with the following command. The server will use port 4000 by default.

```
./bin/echotown
```

> [!TIP]
> Specify a different port with the `-port` flag. For example, `./bin/echotown -port 4040`.

### Client Steps

1. With the server running in a terminal, open another terminal window to use Netcat to connect to the host and port. If the server is running on the same machine using the default port, you can enter the following command:

```
nc localhost 4000
```

## Examining Logs

A `log` folder with server and client logs is created in the directory where the server was run. If you were following the previous steps, this will be in the project root folder. They can be examined with any text editor.

## Tests

Because terminals often have a limit of 1024 characters that can be entered/pasted, tests for input overflow and bad input can be examined in the `internal/server/server_test.go` file. To verify all tests are passing, run the following command in the project root folder:

```
go test ./...
```

> [!NOTE]
> The server tests use port 4001, so make sure the port is not in use when running the tests.

## Cleanup

Using the server and running the tests will create `log` folders with log files throughout the project directory. To remove all of them, as well as the binary, run the following command in the project root folder:

```
make clean
```

## Development Experience

### Most Educationally Enriching Functionality

The most educationally enriching functionality was the logger implementation. Seeing how versatile the `log` standard library package is was cool. Being able to easily determine the log output to standard output, to a file, or both simultaneously was very useful. I got to review file handling concepts I was introduced to in Programming 2, such as file modes and Linux file permissions (I still can't rapidly convert from the octal system, but I recognize to which group the permissions should map). Using my `logger` package helped organize my project and enabled easy extensibility when I needed to implement client message logging.

### Functionality That Required the Most Research

The functionality that required the most research was the input validation and handling. The first thing I had to figure out was what kind of inputs constituted bad input to determine the scope to handle. After learning that terminals often have a 1024-character input limit, I decided that testing the input through a Go test would be the best way. Much time was spent looking at the documentation pages for the `net` and `bufio` standard library packages to test effectively by programmatically writing the input. Ultimately, writing those tests early was worth it, as they became helpful while implementing other features or refactoring code.
