# Redis Clone Server

A toy implementation of a Redis-like server, written in Go. This project demonstrates how to:

- Parse and respond using the RESP (REdis Serialization Protocol)
- Implement commands for manipulating data (Strings, Hashes, etc.)
- Persist data via an Append-Only File (AOF)

While it’s not a fully-featured Redis replacement, it’s a helpful educational tool for understanding how Redis works under the hood.

## Features

1. RESP Implementation

   - A basic parser for incoming client commands
   - Marshalling and unmarshalling of Redis-like data types (bulk strings, arrays, etc.)

2. In-Memory Commands

   - Support for simple string commands (SET, GET, etc.)
   - Hash commands (HSET, HGET, etc.)
   - Easily extensible with new commands

3. Append-Only File (AOF) Persistence

   - Logs write commands to storage.store
   - Replays them on startup to restore state

4. Simple Concurrency Model

   - Listens on TCP port (default :6379)
   - Spawns a goroutine per client connection
   - Uses a mutex for safe file writes

## Folder Structure

```bash
├── LICENSE
├── README.md
├── cmd
│ └── server
│ ├── main.go
│ └── main_test.go
├── commands
│ ├── hashes.go
│ ├── hashes_test.go
│ ├── main.go
│ ├── string.go
│ └── string_test.go
├── go.mod
├── go.sum
├── resp
│ ├── helper.go
│ ├── helper_test.go
│ ├── marshal.go
│ ├── marshal_test.go
│ ├── reader.go
│ ├── reader_test.go
│ └── types.go
├── storage
│ ├── aof.go
│ └── aof_test.go
└── storage.store
```

- cmd/server: Entry point for the application (main binary).
- commands: Contains Redis-like command logic (Strings, Hashes, etc.).
- resp: Implements RESP protocol parsing and marshalling.
- storage: Manages Append-Only File (AOF) creation, writing, reading, and syncing.
- storage.store: The default AOF file used by the server.

## Getting Started

### Prerequisites

- Go 1.18+ (or a reasonably recent version)

### Installation & Running

1. Clone the repository (or download the source):

```bash
git clone https://github.com/<your-username>/redis-clone.git
cd redis-clone
```

2. Build the server:

```bash
cd cmd/server
go build -o redis-clone-server
```

3. Run the server:

```bash
./redis-clone-server
```

By default, it listens on TCP port 6379.

### Connecting to the Server

You can connect with any Redis client (e.g., the official redis-cli) by specifying the host and port:

```bash
redis-cli -p 6379
```

Then use commands like:

```bash
SET mykey "Hello, world!"
GET mykey
```

## Usage Examples

Once connected via redis-cli or any RESP-compatible client:

- String commands

```bash
SET name "Go Developer"
GET name
```

- Hash commands

```bash
HSET user:1000 username "alice"
HGET user:1000 username
```

More commands can be found in the `commands` package.

## Persistence

The server writes incoming write-commands (like SET, HSET) to an append-only file (storage.store).
Upon startup, it replays the commands from storage.store to rebuild in-memory state.

**Note**: The file can grow indefinitely. For serious usage, you would implement a rewrite or snapshot mechanism.

## Development

- Run all tests:

```bash
go test ./...
```

- Run test coverage:

```bash
go test -cover ./...
```

This will run tests in each package (commands, resp, storage, etc.).

## Contributing

Feel free to open issues or submit pull requests if you find a bug or want to add a new feature. This project is primarily for educational purposes, but improvements are always welcome!

## License

This project is licensed under the terms of the MIT License. Please see the LICENSE file for full details.
