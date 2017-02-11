# Protocol

***Note: This document is not (yet) set in stone. It may be amended at any time.***

## Inputs

`master_password` - _A strong password._

`plaintext[n]` - _The user-inputted data that we will be protecting, where `n` refers to the index of the slice of plaintext after splitting it into chunks._

`identifier` - _A string that is used to locate the correct ciphertext on retrieval._

## Definitions and Derivations

### :: `master_key`

`master_key = Scrypt(master_password || identifier)`

The `master_key` is 32 bytes long and is what is used to actually encrypt the plaintext itself.

### :: `root_identifier`

`root_identifier = Scrypt(identifier || master_password)`

The `root_identifier` 32 bytes long and is used to derive individual `derived_identifier[n]` values.

### :: `ciphertext[n]`

`ciphertext[n] = XSalsa20Poly1305(master_key, plaintext[n])`

`ciphertext[n]` refers to the result of encrypting `plaintext[n]` with `master_key`.

### :: `derived_identifier[n]`

`derived_identifier[n] = sha256(root_identifier || n)`

The `derived_identifier[n]` is 32 bytes long and is what is actually stored in the database alongside chunks of the ciphertext. The reason for it is so that we are able to store ciphertexts across multiple entries in the database without leaking information about which entries are linked or how many entries compose the data.

`n` refers to the index of the slice of the ciphertext that we're deriving the identifier for: `derived_identifier[n]` corresponds to `ciphertext[n]`.

## Modus Operandi

### :: Adding an entry

1. Split `plaintext` into chunks of length 1024 bytes. The last chunk will have a length of `len(plaintext) mod 1024`.

2. For each `n`, pad `plaintext[n]` to 1025 bytes using byte-padding: a variant of bit-padding<sup>[0]</sup> but with whole bytes instead of bits.

3. For each `n`, derive `ciphertext[n]` and `derived_identifier[n]`.

4. Save every `derived_identifier[n]`:`ciphertext[n]` pair in the database.

### :: Retrieving an entry

1. Compute `derived_identifier[0]` and look up the corresponding `ciphertext[0]` in the database.

2. Keep computing `derived_identifier[n+1]` and looking for it in the database. Stop when the key does not exist.

3. Decrypt each `ciphertext[n]` that we have to get corresponding `plaintext[n]` values.

4. Unpad each `plaintext[n]` and concatenate the resulting values in order of `n` ascending. This will give us `plaintext`.

### :: Deleting an entry

1. Compute `derived_identifier[0]` and remove it from the database.

2. Keep computing `derived_identifier[n+1]` and removing it from the database. Stop when the key does not exist.

## References

[0] A, Menezes., P, van Oorschot., S, Vanstone. (1996, October 16). Handbook of Applied Cryptography: Algorithm 9.30. Retrieved from http://cacr.uwaterloo.ca/hac/about/chap9.pdf#page=15
