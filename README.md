# NFO Notifications

An application to scrape the event log page on the [NFO Control Panel](https://www.nfoservers.com/control/) and post new
events to a Discord webhook

## Building

Ensure you have [go installed](https://golang.org/doc/install) before continuing.

You can build the application using the following command:

```
$ go build
```

This will create an executable (`NFONotifications` on Unix systems and `NFONotifications.exe` on Windows) you can use to
run the application.

## Usage

The application is configured using command-line flags:

| Flag | Description                               |
| ---- | ----------------------------------------- |
| `-e` | Your account email                        |
| `-p` | Your account password                     |
| `-t` | A valid authentication token              |
| `-h` | The server identifier to fetch events for   |
| `-w` | The Discord webhook URL to post events to |

To successfully authenticate yourself, an existing authentication token or an email and password must be specified. If
both are provided, the application will prioritise using the authentication token.

The server identifier and webhook URL must be provided for the application to run.

### Example Usage

#### Using an authentication token

```
$ ./NFONotifications -t xxxxxxxxxxxxxxxx -h myserver -w https://discord.com/api/webhooks/000000000000000000/xxxxxxxxxxxxxxxx
```

#### Using an email & password

```
$ ./NFONotifications -e hello@example.com -p MyNFOAccount! -h myserver -w https://discord.com/api/webhooks/000000000000000000/xxxxxxxxxxxxxxxx
```

## Contributions

Merge Requests are welcome, there are places where error handling could be improved.