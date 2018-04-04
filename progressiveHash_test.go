package progressiveHash

import "testing"

func testHash(t *testing.T) {
	seed := []byte("test application")
	in := []byte("input string")

	h, err := Hash(in, seed, 50000, 8, 256)
	if err != nil {
		panic(err)
	}

	err = Verify(in, seed, 50000, 8, 256, h)
	if err != nil {
		panic(err)
	}

	err = Verify(in, seed, 10000, 8, 256, h)
	if err != nil {
		panic("verification shouldn't have passed")
	}
}
