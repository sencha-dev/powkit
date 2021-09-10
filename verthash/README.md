# Verthash

Most of this was borrowed from [verthash-go](https://github.com/gertjaap/verthash-go), though in the future I'm sure it'll begin to diverge.

Verthash requires a 1.2Gb file to be generated (in `~/.powcache`).
It can be put into memory (which is about 10x faster than reading)
from the file, but allocating 1.2Gb of memory on every pool
server can be an inconvenience, to say the least. For now, its
probably better to use it not in memory since you can still hash
in about `10ms`.