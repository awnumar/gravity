# pocket

<img align="right" src="https://cdn.rawgit.com/libeclipse/pocket/master/pocket.svg" height="166">

[![license](https://img.shields.io/github/license/libeclipse/pocket.svg)](https://raw.githubusercontent.com/libeclipse/pocket/master/LICENSE) [![Build Status](https://travis-ci.org/libeclipse/pocket.svg?branch=master)](https://travis-ci.org/libeclipse/pocket) [![Go Report Card](https://goreportcard.com/badge/github.com/libeclipse/pocket)](https://goreportcard.com/report/github.com/libeclipse/pocket)

***Note: Still in alpha stages. Should not (yet) be used seriously.***

Protect highly sensitive information completely, even in the case of your password being breached.

## Features

***Not all of these features have yet been implemented.***

* ***Multi-layer security*** - The password alone isn't enough to compromise your secrets.
* ***Hidden data length*** - Secrets are padded and split over multiple entries in a way that makes it impossible to ascertain which ones are linked together. This makes it impossible for an attacker to even determine the length of your data, whether it's a 2 GB file or a simple string.
* ***Multiple password support*** - You're free to use different passwords for different entries, and no one would ever know that you did.
* ***Decoy entries*** - A random number of randomly generated decoys will randomly be added to the secrets store and won't be differentiable from real entries. This will make it plausible to claim that `n` of the entries are real and the rest are decoys, where `n >= 0`.
* ***Deniability*** - Since *pocket* will not stop you from using different passwords, it is possible to add some of your own decoys. In the event of [rubber-hose-cryptanalysis](https://en.wikipedia.org/wiki/Rubber-hose_cryptanalysis), you can give up the password/identifiers for these decoy entries and claim that the rest of them are random decoys added by the program.
* ***Hidden entry identifiers*** - The entry identifiers are hashed so that an attacker cannot even tell what type of data is stored. There have been many cases where users have encrypted their data but file names have still given them away. In *pocket*, this is mitigated.
* ***Cleared logs and metadata*** - Any occurrence of *pocket* will be cleared from your bash history, and metadata in the secrets file that would reveal any dates/times will set to (and kept at) `January 1, 1970`. These measures will prevent anyone from correlating the logs to any entries, and will also hide the fact that you've used the application at all: further backing up the claim that some/all of the entries are decoys.

## Installation

Simply run:

`~ >> go get github.com/libeclipse/pocket`

This will fetch, compile, and install *pocket* automatically. An added bonus is that it should now be in your PATH so you can call the program from anywhere with a simple:

`~ >> pocket`

## Credits

- [@dotcppfile](https://twitter.com/dotcppfile)
- [@mnzt](https://github.com/mnzt)
