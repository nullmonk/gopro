# gopro
Gopro is a simple tool and API to raw decode protobuf blobs.

If you've ever wanted to quickly hack up a golang program which looks at a protobuf blob, consider using gopro.

Gopro is basic, it does not attempt any fancy decoding, it just provides the bytes and the
wiretype in a basic object.

## Usage

Commandline usage:
```bash
go run cmd/gopro.go binary.blob
go run cmd/gopro.go < binary.blob
```

Golang usage:
```Go

func Foo() {
    ...
    msg, err := gopro.Decode(buf)
	if err != nil {
		return
	}

    // Do anything with the message here
	...

    // Print the message
    gopro.DumpMessage(msg)
}
```
