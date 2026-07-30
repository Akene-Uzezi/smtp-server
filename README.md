# smtp-server

A minimal, self-hosted SMTP server written in Go. Accepts incoming emails over SMTP and logs them to stdout and file. Designed for development and testing workflows — not for production mail relay.

## What it does

- Listens for SMTP connections on a configurable port
- Supports `HELO`/`EHLO`, `MAIL FROM`, `RCPT TO`, `DATA`, and `QUIT`
- Parses received messages and logs sender, recipient(s), subject, date, and body
- Persists all accepted mail and server events to `logs/`

No mail is actually delivered — this is a **capture/development SMTP sink**.

## Tech stack

- [Go](https://go.dev/) (standard library only — no external dependencies)

## Installation

```bash
git clone <repo-url>
cd smtp-server
go build -o smtp-server ./cmd/
```

## Usage

```bash
./smtp-server
```

Server starts on `:2525` by default. Point a mail client, application, or another SMTP server at that host/port and send a message. Received mail appears in stdout and in `logs/emailLogger.log`.

Example with `telnet` or `nc`:

```
$ nc localhost 2525
220 whatever.com smtp Server Ready
EHLO example.com
200 message ready to accept message
MAIL FROM:<sender@example.com>
250 2.1.0 Sender OK
RCPT TO:<recipient@example.com>
250 2.1.5 Recipient OK
DATA
354 Start mail input; end with <CR><LF>,<CR><LF>
Subject: test

Hello world
.
message accepted for delivery
QUIT
```

## Self-hosting

### Requirements

- Go 1.26+ (or download the pre-built binary)
- Linux/Unix-like environment (file paths use `filepath.Join`)
- Write permission on the working directory (creates `logs/`)

### Configuration

Port is hardcoded in `cmd/main.go:18`. Change it there if needed:

```go
addr := ":2525"
```

No environment variables, config files, databases, or external services are required.

### DNS (if receiving real mail from the internet)

If you plan to route public mail to this server, you will need DNS records on your domain:

| Record | Purpose |
|---|---|
| **MX** | Point your domain to the server's host |
| **SPF** | `v=spf1 a mx -all` (adjust for your sending sources) |
| **DKIM** | Generate keys and publish a TXT record; this server does **not** sign DKIM itself |
| **DMARC** | `v=DMARC1; p=reject; rua=mailto:you@example.com` |

### Firewall

Open TCP port **2525** (or your custom port) inbound.

### Production run

Build and run directly, or use `screen`/`tmux`/`systemd` to keep it alive:

```bash
go build -o smtp-server ./cmd/
./smtp-server
```

Logs are written to:

- `logs/emailLogger.log` — all received mail
- `logs/successFile.log` — server lifecycle events
- `logs/errFile.log` — errors

## Limitations

- No STARTTLS or implicit TLS
- No SMTP AUTH
- No outbound delivery — mail is only logged locally
- Single recipient per `RCPT TO` command (appends all recipients to the envelope)
- No queue, retry, or persistence beyond log files

## License

No license file is included in this repository.
