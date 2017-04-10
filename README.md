<p align="center">
  <img src="tranquil.png" height="140" />
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

Plausible deniability is defined as ***a condition in which a subject can safely and believably deny knowledge of any particular truth that may exist so as to shield the subject from any responsibility associated with the knowledge of such truth***.

Put simply, that is what Tranquil gives you.

## How it works

While other encryption programs concentrate on hiding the content of data, Tranquil also hides the fact that the data even exists.

It is impossible to differentiate between real and random data in the database, and the two can exist in any proportion. This makes it possible for you to add your own decoy data that you would disclose as real if forced to do so.

An attacker cannot ascertain which data is real, so she cannot know for certain if the decrypted entries compose the entirety of the database. Therefore you have plausible deniability and a deniable encryption scheme.

## Installation

Simply run:

`$ go get github.com/libeclipse/tranquil`

This will fetch, compile, and install Tranquil automatically. If you have `$GOPATH` in your PATH, you should be able to run Pocket with a simple:

`$ tranquil`

## Responsible disclosure

If you are aware of a security bug, notifying us privately is in the interest of all users. We can then discuss it post-mortem.

To do this, please send a PGP encrypted message to my [email](mailto:awn@cryptolosophy.io). My PGP public-key is available on my [keybase](https://keybase.io/awn).

My PGP public-key fingerprint is:

> 5469 F4B9 688C 3FEE E105 0CA3 FAEE B039 F313 3EA8

To import it directly into GPG, run `$ curl https://keybase.io/awn/pgp_keys.asc | gpg --import`.
