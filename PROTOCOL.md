# Technical Information

***Note: This protocol is not (yet) set in stone. It may be amended at any time.***

## :: Definitions

### Inputs

**I<sub>key</sub>** - _A strong password._

**I<sub>p</sub>** - _The user-inputted data that we will be protecting._

**I<sub>id</sub>** - _A string that identifies **I<sub>p</sub>** so it can be located on retrieval._

### Variables

**V<sub>l</sub>** - _The fixed length of plaintext per entry, defined as 1024 bytes._

**X<sub>n</sub>** - _The index of an entry. For example if len(**I<sub>p</sub>**) = 2048, there will be two entries with **X<sub>n</sub>** values of `0` and `1` respectively._

### Derivations

**K<sub>m</sub>** - _A master-key derived from both **I<sub>key</sub>** and **I<sub>id</sub>**._

**K<sub>id</sub>** - _A key derived from both **I<sub>key</sub>** and **I<sub>id</sub>**, that is used to derive **X<sub>V<sub>n</sub></sub>**._

### Outputs

**Z<sub>X<sub>n</sub></sub>** - _A derived identifier to locate a specific entry. Derived from **K<sub>id</sub>** and the respective **X<sub>n</sub>** values._

**C<sub>X<sub>n</sub></sub>** - _A piece of ciphertext with index **X<sub>n</sub>**._

## :: Modus Operandi

In all of the following, **I<sub>p</sub>** is a 1536 byte plaintext.

### Adding an entry

1. Generate **K<sub>m</sub>** - Pass **I<sub>key</sub> || I<sub>id</sub>** to Scrypt (no salt).

2. Generate **K<sub>id</sub>** - Pass **I<sub>id</sub> || I<sub>key</sub>** to Scrypt (no salt).

3. Split **I<sub>p</sub>** into pieces of length 1024. With our **I<sub>p</sub>**, we will get one piece of length 1024 and another piece of length 512.

4. Pad each piece to 1025 bytes using byte padding. [https://en.wikipedia.org/wiki/Padding_(cryptography)#Byte_padding]

5. Encrypt each padded piece separately using XSalsa20 and Poly1305 with the key **K<sub>m</sub>**. In our case, this would give us two values: **C<sub>X<sub>0</sub></sub>** and **C<sub>X<sub>1</sub></sub>**.

6. Generate **Z<sub>X<sub>n</sub></sub>** values for the pieces of ciphertext by computing **sha256(K<sub>id</sub> || X<sub>n</sub>)** for each piece. In our case, we'd compute **sha256(K<sub>id</sub> || 0)** and **sha256(K<sub>id</sub> || 1)**.

7. Add the pairs **Z<sub>X<sub>0</sub></sub>**:**C<sub>X<sub>0</sub></sub>** and **Z<sub>X<sub>1</sub></sub>**:**C<sub>X<sub>1</sub></sub>** to the database.

### Retrieving an entry

1. Generate **K<sub>m</sub>** - Pass **I<sub>key</sub> || I<sub>id</sub>** to Scrypt (no salt).

2. Generate **K<sub>id</sub>** - Pass **I<sub>id</sub> || I<sub>key</sub>** to Scrypt (no salt).

3. Generate **Z<sub>X<sub>0</sub></sub>** by computing **sha256(K<sub>id</sub> || 0)**.

4. Search the database for the key **Z<sub>X<sub>0</sub></sub>** and pull the corresponding value (**C<sub>X<sub>0</sub></sub>**).

5. Keep generating values of **Z<sub>X<sub>n</sub></sub>** and looking for them in the database. Stop when **Z<sub>X<sub>n</sub></sub>** does not exist for the current **X<sub>n</sub>** value. In our case, we'd find two entries with **X<sub>n</sub>** equalling `0` and `1` respectively.

6. Decrypt each **C<sub>X<sub>n</sub></sub>** value that we have.

7. Unpad each decrypted **C<sub>X<sub>n</sub></sub>** value and concatenate the resulting values in order of **X<sub>n</sub>** ascending. In our case, we'd have two pieces of data of lengths 1024 bytes and 512 bytes respectively, so we'd join them in order of **X<sub>0</sub>** || **X<sub>1</sub>**.

8. We now have the original decrypted data. Output it to the user.

### Deleting an entry
