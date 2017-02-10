# Technical Information

***Note: This protocol is not (yet) set in stone. It may be amended at any time.***

## :: Definitions

### Inputs

**I<sub>key</sub>** - *A strong password.*

**I<sub>p</sub>** - *The user-inputted data that we will be protecting.*

**I<sub>id</sub>** - *A string that identifies I<sub>p</sub> so it can be located on retrieval.*

### Variables

**V<sub>l</sub>** - *The fixed length of plaintext per entry, defined as 1024 bytes.*

**V<sub>n</sub>** - *The index of a slice of data. For example if len(I<sub>p</sub>) = 2048, there will be two saved entries with V<sub>n</sub> equal to `0` and `1` respectively.*

### Derivations

**K<sub>m</sub>** - *A master-key derived from both I<sub>key</sub> and I<sub>id</sub>.*

**K<sub>id</sub>** - *A key derived from both I<sub>key</sub> and I<sub>id</sub>, that is used to derive X<sub>V<sub>n</sub></sub>.*

**X<sub>V<sub>n</sub></sub>** - *A derived identifier to locate a specific entry. Derived from K<sub>id</sub> and V<sub>n</sub>.*

## :: Modus Operandi
