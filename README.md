# Progressive hashing

Traditional password stretching functions such as PBKDF2, Scrypt, Bcrypt and Argon2 compute a fixed-length digest from a variable length, low-entropy input. These functions are intentionally slow, possibly with large memory requirements, in order to mitigate dictionary attacks.

The flipside is that in order to verify that an input matches a previously computed hash, the whole computation has to be performed again. This requires significant resources, even if the input doesn't happen to match the expected hash.

This package implements a proof of concept of a _progressive hash function_, a PBKDF2 variant that can detect a mismatching hash for a given input _while doing the computation_, and thus return early instead of performing the full computation.

The first output bits of this function are computed as follows:

```text
H(x): BLAKE2B-512(x)
G(R, x): R recursive iterations of H(x)
o: output
N: number of progressive bits in the output
M: total output length
seed: application-specific seed
B(b, x): bit b of x
R0: initial number of iterations
C := 2

h := H(seed || pad || in)

r := R0
h = G(r, h)
B(0, o) := B(0, h)

r := r * C
h = G(r, h)
B(1, o) := B(0, h)

...

r := r * C
h = G(r, h)
B(n, o) := B(0, h)
```

Individual bits of the prefix are computed sequentially and the work factor doubles after every bit.

Finally, the remaining `M-N` bits are copied from bits of `h` at the same position.

Computing a full new hash requires a full computation.

However, the verification function can compute the first bit, stop early if it doesn't match, compute the next one only if it does and repeat the process until `N` bits have been accumulated.

`M` should be large enough to avoid collisions. The maximum output length is 512 bits.

On the other hand, `N` should be short enough to produce collisions during the fast part of the computation.

# Usage

```go
seed := []byte("test application")
in := []byte("input string")

h, err := progessiveHash.Hash(in, seed, 50000, 8, 256)
if err != nil {
	panic(err)
}

err = progessiveHash.Verify(in, seed, 50000, 8, 256, h)
if err != nil {
	panic(err)
}

// This is unlikely to require a full computation
badIn := []byte("wrong input")
err = progessiveHash.Verify(badIn, seed, 50000, 8, 256, h)
if err != nil {
	panic("verification shouldn't have passed")
}
```
