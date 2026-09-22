# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-09-22

First release.

### Added

- `jev.New` client with functional options: `WithModel`, `WithEndpoint`,
  `WithClient` and `WithLogger`.
- `Jev.Request` — judges a piece of state against any number of questions and
  returns one answer per question, keyed the way it was asked.
- Three question constructors: `Noul` (yes/no), `Choice` (one key out of
  described answers) and `Score` (one point on an ordered scale).
- Instruction context through chained `WithField`, `WithExtractedValue`,
  `WithCompare`, `WithFocus`, `WithNote` and `WithCriteria`, plus `Set` as the
  escape hatch for keys with no typed helper. Unrecognized keys round-trip
  untouched.
- API keys held in a self-redacting type, so printing, logging or marshalling
  a client never discloses the credential.
- `*APIError` carrying `StatusCode` and `Body` for non-2xx responses;
  request assembly validated before anything is sent.
- The `data` package, usable on its own for servers implementing this API or
  for tests that build a payload by hand, with typed criteria via
  `QuestionCriteria`.

[Unreleased]: https://github.com/leonardjke/go-jev/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/leonardjke/go-jev/releases/tag/v0.1.0
