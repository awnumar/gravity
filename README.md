***Note: Still in alpha stages. Should not (yet) be used seriously.***

---

<p align="center">
  <img src="https://cdn.rawgit.com/libeclipse/pocket/master/pocket.svg" height="130" />
  <h3 align="center">Pocket</h3>
  <p align="center">The guardian of super-secret things.</p>
  <p align="center">
    <a href="https://travis-ci.org/libeclipse/pocket"><img src="https://travis-ci.org/libeclipse/pocket.svg?branch=master"></a>
    <a href="https://ci.appveyor.com/project/libeclipse/pocket/branch/master"><img src="https://ci.appveyor.com/api/projects/status/s2enb60sa9asjg87/branch/master?svg=true"></a>
    <a href="https://dependencyci.com/github/libeclipse/pocket"><img src="https://dependencyci.com/github/libeclipse/pocket/badge"></a>
    <a href="https://goreportcard.com/report/github.com/libeclipse/pocket"><img src="https://goreportcard.com/badge/github.com/libeclipse/pocket"></a>
  </p>
</p>

---

Whether you want to encrypt your super-secret files, store your super-secret passwords, save some super-secret strings, log a super-secret diary entry, or have something to look at in wonder -- pocket has you covered.

## How it works

On a high-level, Pocket does some [magic](/PROTOCOL.md) to store your data in such a way that nobody can get the length, type, or content of it; even if they have the right password. They won't even be sure that it actually exists! (Plausible deniability is a wonderful thing.)

The data is all stored in a single database, side-by-side with some optional decoy entries. Along with the multiple-password support, this allows for proper deniable encryption. Just add some legit-looking entries under an alternate password, throw in a few thousand decoys, and there you have it.

Pocket uses XSalsa20 with a Poly1305 MAC for encryption and authentication -- this is implemented with the [NaCl](https://godoc.org/golang.org/x/crypto/nacl/secretbox) package. For key-derivation we use [Scrypt](https://godoc.org/golang.org/x/crypto/scrypt), and for hashing we use [BLAKE2b](https://godoc.org/golang.org/x/crypto/blake2b). Both of these are implemented natively in Golang's crypto library.

## Installation

Simply run:

`$ go get github.com/libeclipse/pocket`

This will fetch, compile, and install *pocket* automatically. An added bonus is that it should now be in your PATH so you can call the program from anywhere with a simple:

`$ pocket`

## Responsible disclosure

If you are aware of a security bug, notifying us privately is in the interest of all users. We can then discuss it post-mortem.

To do this, please send a PGP encrypted message to my [email](mailto:libeclipse@gmail.com). My PGP public-key is available on my [keybase](https://keybase.io/awn).

My PGP public-key fingerprint is:

> 5469 F4B9 688C 3FEE E105 0CA3 FAEE B039 F313 3EA8

To import it directly into GPG, run `$ curl https://keybase.io/awn/pgp_keys.asc | gpg --import`.
