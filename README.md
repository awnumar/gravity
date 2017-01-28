# pocket

<img align="right" src="https://cdn.rawgit.com/libeclipse/pocket/master/pocket.svg" height="164">

[![license](https://img.shields.io/github/license/libeclipse/pocket.svg)](https://raw.githubusercontent.com/libeclipse/pocket/master/LICENSE) [![Build Status](https://travis-ci.org/libeclipse/pocket.svg?branch=master)](https://travis-ci.org/libeclipse/pocket) [![Go Report Card](https://goreportcard.com/badge/github.com/libeclipse/pocket)](https://goreportcard.com/report/github.com/libeclipse/pocket)

***Note: Still in alpha stages. Should not (yet) be used seriously.***

Protect highly sensitive information, even in the case of your password being breached.

## Features

***Not all of these features have yet been implemented.***

* ***Multi-layer security*** - The password alone isn't enough to compromise your secrets.
* ***Hidden data length*** - Secrets are padded and split over multiple entries in a way that makes it impossible to ascertain which ones are linked together. This effectively conceals the length of your data, whether it's a 2 GB file or a simple string.
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

## Usage

### Modes

*Pocket* can be launched in the following modes:

##### :: `pocket help`

This mode outputs a basic help message outlining the different modes and arguments. For further and more in-depth usage information, this project-page (or the README file) should instead be consulted.

##### :: `pocket add [-c N,r,p int]`

The *add* mode is used to add new secrets to the store.

You'll be prompted to enter a password and an identifier. Both of those things together are used to derive the encryption key that protects your secrets, so a strong identifier is recommended alongside a strong password.

For the identifier, you should aim to use a phrase like `l33t encryption key for them thingz init` instead of something like `encryption key` which could easily be guessed. There's also nothing stopping you from using random values for both fields, assuming that you can remember them.

Speaking of not stopping you from doing things, you're also free to use different passwords for different entries. Aside from increasing security, this also has the side effect of allowing deniable encryption. Simply add a few legit-looking secrets with a decoy key and if you're ever forced to disclose your keys, just give up the decoys. The program adds its own decoys so you can claim that the other encrypted entries are just that: decoys.

##### :: `pocket get [-c N,r,p int]`

The *get* mode is used for retrieving secrets from the store.

You'll be prompted to enter a password and an identifier. The program will then derive the secure identifier[s] and encryption key, and then retrieve, decrypt, and output the plaintext.

##### :: `pocket forget [-c N,r,p int]`

The *forget* mode is used for removing secrets from the store.

You'll just need to enter the identifier for the entry and *pocket* will derive the secure identifier, locate the entry, and remove it from the store.

You won't be asked for a confirmation, so when you run forget, make sure that you mean it.

### Options

##### :: `[-c N,r,p int]`

This option is to specify custom scrypt cost-parameters that will be used to derive the identifier[s] and encryption key. If you don't understand what this is, leave it alone. The default values are `18,8,1`.

Things to note:

* The N parameter is as a power of two. So, for example if you want `N = 2^20`, you'd enter `N = 20`.
* If you set a custom cost-factor for an entry, you'll have to specify that same cost-factor every time you run `get` and also when you run `forget`.

## Credits

- [@Alipha](https://github.com/alipha/) - Thought of the idea of how we'd link data together across multiple entries.
- [@dotcppfile](https://twitter.com/dotcppfile) - Brainstormed ideas and suggested countless improvements.
- [@mnzt](https://github.com/mnzt) - Hangs around and reviews stuff. :)
