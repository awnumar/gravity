# pocket

Protect super secret passwords and sketchy snippets - even in the case of your password being leaked.

Each secret is encrypted with its own individual encryption key, derived from the password supplied and a string that identifies the secret. The secret is then stored as a `hashed identifier`: `encrypted secret` pair.

This ensures that even in the event of the password being compromised, an adversary would be unable to decrypt the secret, and would instead have to resort to brute-forcing the identifier.

## Technical Information

### Hash Function

The `Scrypt` key deriviation function is used - with `N = 2^20` for the encryption key and `N = 2^18` for hashing the identifier.

This means that the brute-force attack described above becomes infeasable for any non-trivial identifier.

### Encryption

Rolling your own crypto is bad. That's why Pocket uses the excellent NaCl library's symmetric encryption functions. That's `xSalsa20` with a `Poly1305 MAC` for confidentiality, authenticity, and integrity.

## Installation and Usage

To compile the program, simply run:

`~ >> go build ./pocket.go`

This will create a binary in the current directory called `pocket`.

Optionally for ease of use, you can then create a soft link to a directory that's in your path:

`~ >> ln -s /path/to/pocket /usr/bin/pocket`

This will allow you to call the program from anywhere with a simple:

`~ >> pocket`
