# Docgen

Docgen is a protobuf generator, generating markdown documentation. It is designed to be used with gunk, but can work separately, outside of gunk.

For now, it goes through the same gunk package used
in `gunk generate` command and generates a documentation markdown file `all.md`.

## Installation

Use the following command to install docgen:

```sh
$ go get -u github.com/gunk/docgen
```

This will place `docgen` in your `$GOBIN`

## Usage

Following describes usage with gunk. Outside of gunk, you will need to somehow plug in docgen to protoc.

In your `.gunkconfig` add the following:

```ini
[generate]
command=docgen
out=examples/util/v1/
```

### Code examples

To generate code examples, add the following to the `.gunkconfig` docgen section:

```ini
lang=go
```

Then add your `*.go` files near your gunk files. The examples files must be
named according to the gunk method you want to showcase.

Example:

```go
// UpdateAccount updates an account.
UpdateAccount(UpdateAccountRequest)

// DeleteAccount deletes an account.
DeleteAccount(DeleteAccountRequest)
```

You should have `update_account.go` and `delete_account.go`.

## Tests

Tests are a separate module in `tests`, because they require gunk; while the main module does not.
