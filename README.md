<p align="center">
  <img src="https://cdn.rawgit.com/libeclipse/tranquil/master/tranquil.svg" height="200" />
  <h3 align="center">Tranquil (beta)</h3>
  <p align="center">Confidentiality includes deniability.</p>
  <p align="center">
    <a href="https://travis-ci.org/libeclipse/tranquil"><img src="https://travis-ci.org/libeclipse/tranquil.svg?branch=master"></a>
    <a href="https://ci.appveyor.com/project/libeclipse/tranquil/branch/master"><img src="https://ci.appveyor.com/api/projects/status/cm3cc244ct0yt92s/branch/master?svg=true"></a>
    <a href="https://dependencyci.com/github/libeclipse/tranquil"><img src="https://dependencyci.com/github/libeclipse/tranquil/badge"></a>
    <a href="https://goreportcard.com/report/github.com/libeclipse/tranquil"><img src="https://goreportcard.com/badge/github.com/libeclipse/tranquil"></a>
  </p>
</p>

---

Whether you want to encrypt your super-secret files, store your super-secret passwords, save some super-secret strings, log a super-secret diary entry, or have something to look at in wonder -- pocket has you covered.

## How it works

On a high-level, Pocket does some [magic](/PROTOCOL) to store your data in such a way that nobody can get the length, type, or content of it; even if they have the right password. They won't even be sure that it actually exists! (Plausible deniability is a wonderful thing.)

The data is all stored in a single database, side-by-side with some optional decoy entries. Along with the multiple-password support, this allows for proper deniable encryption. Just add some legit-looking entries under an alternate password, throw in a few thousand decoys, and there you have it. If you're ever compelled to give up your keys, simply give up the alternate password and identifiers.

## Installation

Simply run:

`$ go get github.com/libeclipse/pocket`

This will fetch, compile, and install Pocket automatically. If you have `$GOPATH` in your PATH, you should be able to run Pocket with a simple:

`$ pocket`

## Responsible disclosure

If you are aware of a security bug, notifying us privately is in the interest of all users. We can then discuss it post-mortem.

To do this, please send a PGP encrypted message to my [email](mailto:awn@cryptolosophy.io). My PGP public-key is available on my [keybase](https://keybase.io/awn).

My PGP public-key fingerprint is:

> 5469 F4B9 688C 3FEE E105 0CA3 FAEE B039 F313 3EA8

To import it directly into GPG, run `$ curl https://keybase.io/awn/pgp_keys.asc | gpg --import`.
