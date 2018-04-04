package progressiveHash

import "testing"

func testHash(t *testing.T) {
	personalisation := []byte("test application")
	in := []byte("input string")

	h, err := Hash(in, personalisation, 50000, 8, 256)
	if err != nil {
		panic(err)
	}

	err = Verify(in, personalisation, 50000, 8, 256, h)
	if err != nil {
		panic(err)
	}

	err = Verify(in, personalisation, 10000, 8, 256, h)
	if err != nil {
		panic("verification shouldn't have passed")
	}
}
