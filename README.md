<p align="center">
  <img src="https://cdn.rawgit.com/libeclipse/pocket/documentation/prettify-readme/images/pocket.svg" height="150" />
  <img src="https://cdn.rawgit.com/libeclipse/pocket/documentation/prettify-readme/images/pocket-text.svg" />
  <p align="center">The guardian of super-secret things.</p>
  <p align="center">
    <a href="https://travis-ci.org/libeclipse/pocket"><img src="https://travis-ci.org/libeclipse/pocket.svg?branch=master"></a>
    <a href="https://ci.appveyor.com/project/libeclipse/pocket/branch/master"><img src="https://ci.appveyor.com/api/projects/status/s2enb60sa9asjg87/branch/master?svg=true"></a>
    <a href="https://dependencyci.com/github/libeclipse/pocket"><img src="https://dependencyci.com/github/libeclipse/pocket/badge"></a>
    <a href="https://goreportcard.com/report/github.com/libeclipse/pocket"><img src="https://goreportcard.com/badge/github.com/libeclipse/pocket"></a>
  </p>
</p>

***Note: Still in alpha stages. Should not (yet) be used seriously.***

Protect super-secret ***things*** even in the case of your password being breached.

**You could use *pocket* to:**

* Encrypt your *super-secret* files.
* Store your *super-secret* passwords.
* Save some *super-secret* strings.
* Log a *super-secret* diary entry.
* Just look at in wonder when you're bored.

## Security Properties

***Not all of these have yet been implemented.***

* ***Multi-layer security*** - The password alone isn't enough to compromise your data.
* ***Hidden data length*** - Secrets are padded and split over multiple entries in a way that makes it impossible to ascertain which ones are linked together. This (used with the decoy entries) effectively conceals the length of your data, whether it's a 2 GB file or a simple string.
* ***Multiple password support*** - You're free to use different passwords for different entries, and no one (except you) would ever know that you did.
* ***Decoy entries*** - You can add decoy data that is **not** be differentiable from real data that you store. This lets you claim that some (or all) of the entries in the database aren't real and therefore that you're unable to give up keys for them.
* ***Deniability*** - Since multiple-passwords are a thing, it's possible to add some of your own decoys under a different password to your normal one. In the event of [rubber-hose-cryptanalysis](https://en.wikipedia.org/wiki/Rubber-hose_cryptanalysis), you can give up the passwords/identifiers for these fake entries and claim that the rest of the data is composed of random decoys added by the program. But, you might be asking, *why stop there*? Add as many layers of decoys as you want! Forced to reveal your passwords? Give them one set. They don't believe that's all of them? Give up another set... And another... And another... Just keep dropping them bombs. (*I'm only half joking.*)
* ***Hidden everything*** - *Pocket* makes sure that an attacker will not be able to ascertain the length, type, or content of any data, and also prevents the inference of things like the number of secrets stored.

For a full overview of the protocol, click [here](/PROTOCOL.md).

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

For the identifier, you should aim to use a phrase like `l33t encryption key for them thingz init` instead of something like `encryption key` which could easily be guessed. However, there's nothing stopping you from using random values for both fields, assuming that you can remember them.

You're also free to use different passwords for different entries. Aside from increasing security, this also has the side effect of allowing deniable encryption. Simply add a few legit-looking secrets with a decoy key and if you're ever forced to disclose your keys, just give up the decoys. The program adds its own decoys so you can claim that the other encrypted entries are just that: decoys.

##### :: `pocket get [-c N,r,p int]`

The *get* mode is used for retrieving secrets from the store.

You'll be prompted to enter a password and an identifier. The program will derive the secure identifier[s] and encryption key, and then retrieve, decrypt, and output the plaintext.

##### :: `pocket forget [-c N,r,p int]`

The *forget* mode is used for removing secrets from the store.

You'll just need to enter the identifier for the entry and *pocket* will derive the secure identifier, locate the entry, and remove it from the store.

You won't be asked for a confirmation, so when you run forget, make sure that you mean it.

### Options

##### :: `[-c N,r,p int]`

This option is to specify custom scrypt cost-parameters that will be used to derive the identifier[s] and encryption key. If you don't understand what this is, leave it alone. The default values are `18,8,1`.

Things to note:

* The N parameter should be specified as the power of two. So for example if you want `N = 2^20`, you'd enter `N = 20`.
* If you set a custom cost-factor for an entry, you'll have to specify that same cost-factor every time you run `get` and also when you run `forget`.

## Credits

- [@Alipha](https://github.com/alipha/) - Thought of the idea of how we'd link data together across multiple entries.
- [@dotcppfile](https://twitter.com/dotcppfile) - Brainstormed ideas and suggested countless improvements.
- [@mnzt](https://github.com/mnzt) - Hangs around and reviews stuff. :)
