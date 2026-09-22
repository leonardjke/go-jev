# Contributing

Thanks for taking the time. Issues and pull requests are both welcome.

## Getting set up

```sh
git clone https://github.com/leonardjke/go-jev
cd go-jev
go test ./...
```

Go 1.24 or later. There is nothing else to install — the only dependency is
`testify`, and CI runs the same commands the `Makefile` does.

## Before opening a pull request

```sh
make        # gofmt, go vet, go test -race
```

CI additionally checks that `go mod tidy` leaves no diff and runs
`golangci-lint` with the config in `.golangci.yml`.

Keep the change focused: one concern per pull request, and a commit message
that says what changed and why.

## Tests

Unit tests run offline and should stay that way — the transport is exercised
against an `httptest` server. Add a case alongside the existing ones in
`request_test.go` or `redaction_test.go`.

End-to-end tests are behind the `e2e` build tag and talk to the live API. They
skip unless `API_KEY` is set:

```sh
cp .env.example .env    # fill in API_KEY
API_KEY=... API_URL=... go test -tags e2e -v ./...
```

Never commit a real key. `.env` is gitignored; `.env.example` is the one that
gets checked in, and it stays empty.

## Public API

The examples in `example_test.go` and `data/example_test.go` are compiled and
run by `go test`, and their `// Output:` blocks are asserted. If you change
what goes on the wire, those are the tests that will tell you.

Anything exported gets a doc comment. New instruction context keys generally
want a typed `With*` helper on `Question`, mirroring the one on
`data.Instructions`; `Set` stays the escape hatch for the rest.

The API is pre-1.0, so breaking changes are possible — but note them in
`CHANGELOG.md` under `Unreleased` so the next release picks them up.

## Licence

Contributions are accepted under the Apache 2.0 licence that covers this
project.
