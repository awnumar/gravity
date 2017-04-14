<p align="center">
  <img src="https://cdn.rawgit.com/libeclipse/pocket/master/dissident.svg" height="140" />
  <h3 align="center">Dissident (beta)</h3>
  <p align="center">An encryption program that prioritises deniability.</p>
  <p align="center">
    <a href="https://travis-ci.org/libeclipse/dissident"><img src="https://travis-ci.org/libeclipse/dissident.svg?branch=master"></a>
    <a href="https://ci.appveyor.com/project/libeclipse/dissident/branch/master"><img src="https://ci.appveyor.com/api/projects/status/9v38wh14fa6klc7v/branch/master?svg=true"></a>
    <a href="https://dependencyci.com/github/libeclipse/dissident"><img src="https://dependencyci.com/github/libeclipse/dissident/badge"></a>
    <a href="https://goreportcard.com/report/github.com/libeclipse/dissident"><img src="https://goreportcard.com/badge/github.com/libeclipse/dissident"></a>
  </p>
</p>

---

Plausible deniability is defined as *a condition in which a subject can safely and believably deny knowledge of any particular truth that may exist so as to shield the subject from any responsibility associated with the knowledge of such truth*.

We think that's a beautiful idea, and so, that is what Dissident gives you.

## What it does

**Dissident hides the length, content, and existence of data...**

To achieve this, Dissident stores data in fixed-sized entries, each indifferentiable from any other. You are able to generate and add as many random entries as you like, and so---since no one would ever know if you did---you can claim that any data that is in the database is composed of decoys.

All data can be decoys, so an adversary cannot be sure that there is any real data.

**...even if the master password is breached.**

An attacker would also need access to the secure *identifier* for the ciphertext that they want.

Both the master password together with the identifier are needed to locate the correct ciphertext amongst the entries and also to derive the encryption key to unlock the data.

**And of course: deniable encryption.**

You are able to use any master password for an entry, and under each master password you can store as many ciphertexts as you'd like.

So if you pick a second *decoy* master password and store a bunch of legitimate-looking entries under it, then---if you are ever forced to disclose your keys---you can give up this decoy master password and its identifiers. The rest of the data in the database is simply composed of random decoys added by Dissident. (\*wink wink\*)

The complete protocol can be found [here](PROTOCOL).

## Installation

Simply run:

`$ go get github.com/libeclipse/dissident`

This will fetch, compile, and install Dissident automatically. If you have `$GOPATH` in your PATH, you should be able to run it with a simple:

`$ dissident`

## Responsible disclosure

If you are aware of a security bug, notifying us privately is in the interest of all users. We can then discuss it post-mortem.

To do this, please send a PGP encrypted message to my [email](mailto:awn@cryptolosophy.io). My PGP public-key is available on my [keybase](https://keybase.io/awn).

My PGP public-key fingerprint is:

> 5469 F4B9 688C 3FEE E105 0CA3 FAEE B039 F313 3EA8

To import it directly into GPG, run `$ curl https://keybase.io/awn/pgp_keys.asc | gpg --import`.
