# Convoy

A super simple message queue. A message is just a string.

Intended to be used with [gown](https://github.com/lfknudsen/gown).

## Setup

If you have `make` and [install](https://github.com/lfknudsen/install),
you can compile the executables very easily:
```
make
```

Otherwise, to compile the message queue:
```
go build -C cmd/convoy -o ../..
```
To compile the sample client implementation:
```
go build -C cmd/convoy-client -o ../..
```

## Usage

Start the server:

```
convoy [(--port|-p) <port number>]
```

Using the client implementation, you can add a message to the queue with
```
convoy-client [(--port|-p) <port number>] <message(s)...>
```

