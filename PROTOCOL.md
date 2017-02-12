# Protocol

***Note: This document is not (yet) set in stone. It may be amended at any time.***

## Inputs

`master_password` - _A strong password._

`plaintext` - _The user-inputted data that we will be protecting. When `plaintext` is split, individual chunks will be referred to as `plaintext[n]`, where `n` is the index of the chunk._

`identifier` - *A string that identifies the plaintext that we are storing. It should not be guessable but does not necessarily have to be as strong as `master_password`.*

## Derivations

#### :: `root_key`

> `root_key = Scrypt(master_password || identifier)`

This is 64 bytes long and is what is used to derive `master_key` and `root_identifier`.

The default Scrypt parameters are `N = 2^18`, `r = 16`, `p = 1`.

#### :: `master_key`

> `master_key = root_key[0:32]`

This is 32 bytes long and is what is used as the actual encryption key for all `plaintext[n]`.

#### :: `ciphertext[n]`

> `ciphertext[n] = XSalsa20Poly1305(master_key, plaintext[n])`

`ciphertext[n]` refers to the result of encrypting `plaintext[n]` with `master_key`.

`XSalsa20Poly1305()` is implemented using `NaCl:SecretBox`.

#### :: `root_identifier`

> `root_identifier = root_key[32:64]`

A 32 byte value that is used to derive `derived_identifier[n]`.

#### :: `derived_identifier[n]`

> `derived_identifier[n] = sha256(root_identifier || n)`

A 32 byte value that is stored in the database alongside chunks of the ciphertext. The reason we use this instead of `root_identifier` is so that we are able to store ciphertexts across multiple entries by simply incrementing `n` for every chunk of plaintext. This prevents leakage of information about which entries are linked or how many entries compose `plaintext`.

`n` refers to the index of the chunk of ciphertext that we're deriving the identifier for: `derived_identifier[n]` corresponds to `ciphertext[n]`.

## Modus Operandi

### :: Adding an entry

1. Split `plaintext` into chunks of length 1024 bytes. The last chunk will have a length of `len(plaintext) mod 1024`.

2. For each `n`, pad `plaintext[n]` to 1025 bytes.

3. For each `n`, compute `ciphertext[n]` and `derived_identifier[n]`.

4. Save every `derived_identifier[n]` : `ciphertext[n]` pair to the database.

### :: Retrieving an entry

1. Compute `derived_identifier[0]` and search for it in the database to get `ciphertext[0]`.

2. Keep computing `derived_identifier[n+1]` and looking for it in the database. Stop when nothing is found.

3. Decrypt each `ciphertext[n]` to get a set of `plaintext[n]` values.

4. Unpad each `plaintext[n]` and concatenate the resulting values in order of `n` ascending. This will give us `plaintext`.

### :: Deleting an entry

1. Compute `derived_identifier[0]` and remove it from the database.

2. Keep computing `derived_identifier[n+1]` and removing it from the database. Stop when the key does not exist.

## Miscellaneous

### :: Decoys

The user will have the option to add a certain amount of decoy data. In order to minimise any assumptions that an adversary can make, the number of decoy entries added should not be a predictable number like 10000.

1. Generate three, random, 32 byte values for `master_password`, `plaintext` and `identifier` respectively.

2. Pad `plaintext` to 1025 bytes.

3. Store the `sha256(identifier)` : `XSalsa20Poly1305(master_password, padded_plaintext)` pair in the database.

4. Repeat steps 1 - 3 until a sufficient number of decoys have been added.

Something to note is that the user does not necessarily have to make use of this feature. Rather, simply the fact that it exists allows the user to claim that some or all of the entries in the database are decoys.

### :: Padding

The padding scheme that is used is byte-padding: a variant of bit-padding<sup>[0]</sup> but with whole bytes instead of bits. The reason for this is because it doesn't require the length of padding to be encoded into the padding itself, thereby doing away with problems that arise when `len(padding)` does not fit inside a single byte.

## References

[0] A, Menezes., P, van Oorschot., S, Vanstone. (1996, October 16). Handbook of Applied Cryptography: Algorithm 9.30. Retrieved from http://cacr.uwaterloo.ca/hac/about/chap9.pdf#page=15
