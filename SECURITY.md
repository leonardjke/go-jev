# Security Policy

## Supported versions

This project is pre-1.0. Fixes land on the latest minor release.

## Reporting a vulnerability

Please do not open a public issue for a security problem.

Report it through GitHub's private advisory form:
https://github.com/leonardjke/go-jev/security/advisories/new

Include what you found, how to reproduce it, and what an attacker could do
with it. You can expect an initial response within a few days.

## A note on API keys

The client holds its key in a type that redacts itself through `fmt`,
`String()`, JSON and `encoding.TextMarshaler`, so printing or logging a client
does not disclose the credential. Reaching the real value takes an explicit
`string()` conversion.

If you find a path that leaks the key — a format verb, a marshaller, an error
message, a debug log line — that is a vulnerability in this library and is
worth reporting.
