# Cuckoo

There are many variations of the Cuckoo Cycle algorithm, this is the variation
specifically for Aeternity. That being said, the only functional difference
is the `ROTL` on the `hasher.XorLanes()` response (current versions of cuckoo
do not do this). The header construction varies a bit too, Aeternity uses
a `uint64` nonce instead of `uint32`.