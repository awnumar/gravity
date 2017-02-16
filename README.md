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
* ***Hidden data length*** - Data is split across multiple entries, effectively concealing the length.
* ***Multiple password support*** - You're free to use different passwords for different entries, and no one (except you) would ever know that you did.
* ***Decoy entries*** - You can add decoy data that is **not** be differentiable from real data that you store. This lets you claim that some (or all) of the entries in the database aren't real and therefore that you're unable to give up keys for them.
* ***Deniability*** - Multiple-password support combined with decoys basically gives you deniable encryption. Simply add a few entries under a different password to your normal one and if you're ever forced to disclose your keys, give up these and claim that the rest of the data is composed of random decoys.
* ***Hidden everything*** - Pocket makes sure that an attacker will not be able to ascertain the length, type, or content of any data, and also prevents the inference of things like the number of secrets stored.

For a full overview of the protocol, click [here](/PROTOCOL.md).

## Installation

Simply run:

`~ >> go get github.com/libeclipse/pocket`

This will fetch, compile, and install *pocket* automatically. An added bonus is that it should now be in your PATH so you can call the program from anywhere with a simple:

`~ >> pocket`
