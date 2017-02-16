<p align="center">
  <img src="https://cdn.rawgit.com/libeclipse/pocket/documentation/prettify-readme/images/pocket.svg" height="130" />
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

Whether you want to encrypt your super-secret files, store your super-secret passwords, save some super-secret strings, log a super-secret diary entry, or have something to look at in wonder; pocket has you covered.

## Security Properties

* ***Multi-layer security*** - The password alone isn't enough to compromise your data.
* ***Hidden data length*** - Secrets are padded and split over multiple entries in a way that makes it impossible to ascertain which ones are linked together. This (used with the decoy entries) effectively conceals the length of your data, whether it's a 2 GB file or a simple string.
* ***Multiple password support*** - You're free to use different passwords for different entries, and no one (except you) would ever know that you did.
* ***Decoy entries*** - You can add decoy data that is **not** be differentiable from real data that you store. This lets you claim that some (or all) of the entries in the database aren't real and therefore that you're unable to give up keys for them.
* ***Deniability*** - Since multiple-passwords are a thing, it's possible to add some of your own decoys under a different password to your normal one. In the event of [rubber-hose-cryptanalysis](https://en.wikipedia.org/wiki/Rubber-hose_cryptanalysis), you can give up the passwords/identifiers for these fake entries and claim that the rest of the data is composed of random decoys added by the program. But why stop there? Add as many layers of decoys as you want! Forced to reveal your passwords? Give them one set. They don't believe that's all of them? Give up another set... And another... And another... (*I'm only half kidding.*)
* ***Hidden everything*** - Pocket makes sure that an attacker will not be able to ascertain the length, type, or content of any data, and also prevents the inference of things like the number of secrets stored.

For a full overview of the protocol, click [here](/PROTOCOL.md).

## Installation

Simply run:

`~ >> go get github.com/libeclipse/pocket`

This will fetch, compile, and install *pocket* automatically. An added bonus is that it should now be in your PATH so you can call the program from anywhere with a simple:

`~ >> pocket`
