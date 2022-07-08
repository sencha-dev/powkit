# Cuckoo

There are many variations of the Cuckoo Cycle algorithm - here only the ones
that are needed are implemented. `Cuckatoo32` probably should be implemented
as well for Grin. The differences for each variation are as follows (along with
header generation):

  - Aeternity uses Cuckoo29 with a legacy version of the `sipnode` hasher -  
  the only functional difference is the `ROTL` on the `hasher.XorLanes()` 
  response (current versions of cuckoo do not do this)
```go
func generateHeader(hash []byte, nonce uint64) []byte {
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, nonce)
	hashEncoded := []byte(base64.StdEncoding.EncodeToString(hash))
	nonceEncoded := []byte(base64.StdEncoding.EncodeToString(nonceBytes))
	header := append(hashEncoded, append(nonceEncoded, make([]byte, 24)...)...)

	return header
}
```
  - Cortex uses Cuckaroo30 with a single difference - the `sipblock` hasher
uses `siphash48` instead of `siphash24`.
```go
func generateHeader(hash []byte, nonce uint64) []byte {
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, nonce)
	header := append(hash, nonceBytes...)

	return header
}
  ```
