# nbiot-go
[![GoDoc](https://godoc.org/github.com/telenordigital/nbiot-go?status.svg)](https://godoc.org/github.com/telenordigital/nbiot-go)
[![Travis-CI](https://api.travis-ci.com/telenordigital/nbiot-go.svg)](https://travis-ci.com/telenordigital/nbiot-go)

NBIoT-Go provides a Go client for the [REST API](https://api.nbiot.telenor.io) for
[Telenor NB-IoT](https://nbiot.engineering).

## Important note

Version 2.1.1 of the service at nbiot.engineering introduces a breaking change
for the go client library. Versions prior to 2.1.1 won't work with v0.6.0 of
this library.

## Configuration

The configuration file is usually located at `${HOME}/.telenor-nbiot`,
but the library will start at the current directory and scan each
parent directory until it finds a config file. The file is a simple
list of key/value pairs. Additional values are ignored. Comments must
start with a `#`:

    #
    # This is the URL of the Telenor NB-IoT REST API. The default value is
    # https://api.nbiot.telenor.io and can usually be omitted.
    address=https://api.nbiot.telenor.io

    #
    # This is the API token. Create new token by logging in to the Telenor NB-IoT
    # front-end at https://nbiot.engineering and create a new token there.
    token=<your api token goes here>


The configuration file settings can be overridden by setting the environment
variables `TELENOR_NBIOT_ADDRESS` and `TELENOR_NBIOT_TOKEN`. If you only use environment variables
the configuration file can be ignored.

Use the `NewWithAddr` function to bypass the default configuration file and
environment variables when you want to configure the client programmatically.

## Updating resources

The various `Client.Update*` methods work via HTTP PATCH, which means they will only modify or set fields, not delete them.  There are special `Client.Delete*Tag` methods for deleting tags.

## Example

```go
package main

import (
	"io"
	"log"

	"github.com/telenordigital/nbiot-go"
)

func main() {
	client, err := nbiot.New()
	if err != nil {
		log.Fatal("Error creating client:", err)
	}

	collection, err := client.CreateCollection(nbiot.Collection{
		Tags: map[string]string{
			"name": "example collection",
		},
	})
	if err != nil {
		log.Fatal("Error creating collection: ", err)
	}

	device, err := client.CreateDevice(collection.ID, nbiot.Device{
		IMSI: "0345892703458",
		IMEI: "1487252347803",
		Tags: map[string]string{
			"name": "example device",
		},
	})
	if err != nil {
		log.Fatal("Error creating device: ", err)
	}

	stream, err := client.DeviceOutputStream(collection.ID, device.ID)
	if err != nil {
		log.Fatal("Error creating stream: ", err)
	}
	defer stream.Close()

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error receiving data: ", err)
		}

		log.Print("received payload: ", data.Payload)
	}
}
```
