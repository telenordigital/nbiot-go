# horde-go
Horde-Go provides a Go client for the [REST API](https://api.nbiot.telenor.io) for
[Telenor NB-IoT](https://nbiot.engineering).

## Configuration

The configuration file is located at `${HOME}/.horde`. The file is a simple
list of key/value pairs. Additional values are ignored. Comments must start
with a `#`:

    #
    # This is the URL of the Horde REST API. The default value is
    # https://api.nbiot.telenor.io and can usually be omitted.
    address=https://api.nbiot.telenor.io

    #
    # This is the API token. Create new token by logging in to the Horde
    # front-end at https://nbiot.engineering and create a new token there.
    token=<your api token goes here>


The configuration file settings can be overridden by setting the environment
variables `HORDE_ADDRESS` and `HORDE_TOKEN`. If you only use environment variables
the configuration file can be ignored.

Use the `NewWithAddr` function to bypass the default configuration file and
environment variables when you want to configure the client programmatically.
