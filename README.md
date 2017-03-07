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

On a high-level, Pocket does some [magic](/PROTOCOL.md) to store your data in such a way that nobody can get the length, type, or content of it. They won't even be sure that it exists at all! Plausible deniability is a wonderful thing.

The data is all stored in a single database, side-by-side with some optional decoy entries. Along with the multiple-password support, this allows for proper deniable encryption. Just add some legit-looking entries under an alternate master_password, throw in a few thousand decoys, and there you have it.

And if someone does manage to find out your password, it alone isn't enough to even locate the real entries amongst the decoys, never mind decrypt them.

## Installation

Simply run:

`$ go get github.com/libeclipse/pocket`

This will fetch, compile, and install *pocket* automatically. An added bonus is that it should now be in your PATH so you can call the program from anywhere with a simple:

`$ pocket`

## Reporting security issues

If you are aware of a security bug, notifying us privately is in the interest of all users. We can then discuss it post-mortem.

To do this, please send a PGP encrypted message to my [email](libeclipse@gmail.com). My PGP public-key is available [here](https://keybase.io/awn) [`5469 F4B9 688C 3FEE E105 0CA3 FAEE B039 F313 3EA8`].

To import it directly into GPG, run `curl https://keybase.io/awn/pgp_keys.asc | gpg --import`.
