# pocket

[![Build Status](https://travis-ci.org/libeclipse/pocket.svg?branch=master)](https://travis-ci.org/libeclipse/pocket) [![Go Report Card](https://goreportcard.com/badge/github.com/libeclipse/pocket)](https://goreportcard.com/report/github.com/libeclipse/pocket)

Protect super secret passwords and sketchy snippets - even in the case of your password being leaked.

Each secret is encrypted with its own individual encryption key, derived from the password supplied and a string that identifies the secret. The secret is then stored as a `hashed identifier`:`encrypted secret` pair.

This ensures that even in the event of the password being compromised, an adversary would be unable to decrypt the secret, and would instead have to resort to brute-forcing the identifier.

Another side-effect of this arrangement is that it doesn't disallow you from using multiple passwords. So (if you felt like it) you could use separate passwords for different secrets, and no one would ever know.

## Technical Information

### Hash Function

The `Scrypt` key deriviation function is used with `N = 2^20`.

This means that the brute-force attack described above becomes infeasible for any non-trivial identifier.

### Encryption

Rolling your own crypto is bad. That's why Pocket uses the excellent NaCl library's symmetric encryption functions. That's `xSalsa20` with a `Poly1305 MAC` for confidentiality, authenticity, and integrity.

## Installation and Usage

Simply run:

`~ >> go get github.com/libeclipse/pocket`

This will fetch, compile, and install Pocket automatically. An added bonus is that it should now be in your PATH so you can call the program from anywhere with a simple:

`~ >> pocket`
