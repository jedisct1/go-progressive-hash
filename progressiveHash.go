package progressiveHash

import (
	"errors"

	"golang.org/x/crypto/blake2b"
)

// Hash - Hash `in` with seedisation `seed`, an initial number of rounds `initialRounds`,
// and increase the work factor `progressiveLength` times to eventually produce an `outLength` bit-long digest.
func Hash(in []byte, seed []byte, initialRounds uint64, progressiveLength int, outLength int) ([]byte, error) {
	return hash(in, seed, initialRounds, progressiveLength, outLength, nil)
}

// Verify - Verify a previously computed hash `expected` using the input `in`, a seedisation `seed`,
// increasing the work factor `progressiveLength` times to match the initial hash length `outLength` bits.
func Verify(in []byte, seed []byte, initialRounds uint64, progressiveLength int, outLength int, expected []byte) error {
	_, err := hash(in, seed, initialRounds, progressiveLength, outLength, &expected)
	return err
}

func hash(in []byte, seed []byte, initialRounds uint64, progressiveLength int, outLength int, expected *[]byte) ([]byte, error) {
	if progressiveLength < 0 || outLength < 1 {
		return nil, errors.New("invalid output length")
	}
	if len(seed) > 127 {
		return nil, errors.New("seed too long")
	}
	if outLength < progressiveLength {
		outLength = progressiveLength
	}
	xin := make([]byte, 128+len(in))
	copy(xin, seed)
	copy(xin[128:], in)
	h := blake2b.Sum512(xin)
	outLengthBytes := (outLength + 7) / 8
	if outLengthBytes > len(h) {
		return nil, errors.New("output too long")
	}
	var out []byte
	if expected != nil {
		if len(*expected) != outLengthBytes {
			return out, errors.New("expected length not matching the given parameters")
		}
	} else {
		out = make([]byte, outLengthBytes)
	}
	xrounds := initialRounds
	for i := uint(0); i < uint(progressiveLength); i++ {
		for j := uint64(0); j < xrounds; j++ {
			h = blake2b.Sum512(h[:])
		}
		if expected != nil {
			if ((*expected)[i>>3]>>(i&7))&1 != uint8(h[0]&1) {
				return out, errors.New("mismatch")
			}
		} else {
			out[i>>3] |= uint8(h[0]&1) << (i & 7)
		}
		xrounds *= 2
	}
	if expected != nil {
		ck := uint8(0)
		for i := uint(progressiveLength); i < uint(outLength); i++ {
			ck |= ((*expected)[i>>3]>>(i&7))&1 ^ uint8((h[i>>3]>>(i&7))&1)
		}
		if ck != 0 {
			return out, errors.New("mismatch")
		}
	} else {
		for i := uint(progressiveLength); i < uint(outLength); i++ {
			out[i>>3] |= uint8((h[i>>3]>>(i&7))&1) << (i & 7)
		}
	}
	return out, nil
}
