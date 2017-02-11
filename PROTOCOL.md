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

1. Generate **K<sub>m</sub>** - Pass **I<sub>key</sub> || I<sub>id</sub>** to Scrypt (no salt).

2. Generate **K<sub>id</sub>** - Pass **I<sub>id</sub> || I<sub>key</sub>** to Scrypt (no salt).

3. Split **I<sub>p</sub>** into pieces of length 1024. With our **I<sub>p</sub>**, we will get one piece of length 1024 and another piece of length 512.

4. Pad each piece to 1025 bytes using byte padding. [https://en.wikipedia.org/wiki/Padding_(cryptography)#Byte_padding]

5. Encrypt each padded piece separately using XSalsa20 and Poly1305 with the key **K<sub>m</sub>**. In our case, this would give us two values: **C<sub>X<sub>0</sub></sub>** and **C<sub>X<sub>1</sub></sub>**.

6. Generate **Z<sub>X<sub>n</sub></sub>** values for the pieces of ciphertext by computing **sha256(K<sub>id</sub> || X<sub>n</sub>)** for each piece. In our case, we'd compute **sha256(K<sub>id</sub> || 0)** and **sha256(K<sub>id</sub> || 1)**.

7. Add the pairs **Z<sub>X<sub>0</sub></sub>**:**C<sub>X<sub>0</sub></sub>** and **Z<sub>X<sub>1</sub></sub>**:**C<sub>X<sub>1</sub></sub>** to the database.

### :: Retrieving an entry

1. Generate **K<sub>m</sub>** - Pass **I<sub>key</sub> || I<sub>id</sub>** to Scrypt (no salt).

2. Generate **K<sub>id</sub>** - Pass **I<sub>id</sub> || I<sub>key</sub>** to Scrypt (no salt).

3. Generate **Z<sub>X<sub>0</sub></sub>** by computing **sha256(K<sub>id</sub> || 0)**.

4. Search the database for the key **Z<sub>X<sub>0</sub></sub>** and pull the corresponding value (**C<sub>X<sub>0</sub></sub>**).

5. Keep generating values of **Z<sub>X<sub>n</sub></sub>** and looking for them in the database. Stop when **Z<sub>X<sub>n</sub></sub>** does not exist for the current **X<sub>n</sub>** value. In our case, we'd find two entries with **X<sub>n</sub>** equalling `0` and `1` respectively.

6. Decrypt each **C<sub>X<sub>n</sub></sub>** value that we have.

7. Unpad each decrypted **C<sub>X<sub>n</sub></sub>** value and concatenate the resulting values in order of **X<sub>n</sub>** ascending. In our case, we'd have two pieces of data of lengths 1024 bytes and 512 bytes respectively, so we'd join them in order of **X<sub>0</sub>** || **X<sub>1</sub>**.

8. We now have the original decrypted data. Output it to the user.

### :: Deleting an entry

1. Generate **K<sub>id</sub>** - Pass **I<sub>id</sub> || I<sub>key</sub>** to Scrypt (no salt).

2. Generate **Z<sub>X<sub>0</sub></sub>** by computing **sha256(K<sub>id</sub> || 0)**.

3. Search the database for the key **Z<sub>X<sub>0</sub></sub>** and remove it.

4. Keep generating values of **Z<sub>X<sub>n</sub></sub>**, looking for them in the database and removing them. Stop when **Z<sub>X<sub>n</sub></sub>** does not exist for the current **X<sub>n</sub>** value. In our case, we'd find and remove two entries with **X<sub>n</sub>** equalling `0` and `1` respectively.
