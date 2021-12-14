# Autolykos2

This is a very rough draft that seems to work off of the
single test vector I've found. There are still some things
I'm not quite sure about, like the nonce being `uint64` so
this is definitely not ready yet.

Since this has no good test vectors and I haven't used it
in production, I'll have to either use a stratum proxy
off an existing pool or write a basic pool to generate
a good set of test vectors.

I'm fairly surprised this works since I really don't understand
Autolykos very well. It needs to be cleaned up but I'll have to 
spend time understanding the algorithm before that happens. 

Locally it takes about 750µs, with 90% of that being `Blake2b256`
time. There really isn't much room for optimization, but it could
probably be dropped down by 5% or so. It's probably not worth it
since 750µs is 10-20x faster than any DAG based algorithm.