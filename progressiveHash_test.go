package progressiveHash

import "testing"

func TestHash(t *testing.T) {
	seed := []byte("test application")
	in := []byte("input string")

	h, err := Hash(in, seed, 50000, 8, 256)
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}

	err = Verify(in, seed, 50000, 8, 256, h)
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}

	err = Verify(in, seed, 10000, 8, 256, h)
	if err == nil {
		t.Fatal("verification shouldn't have passed")
	}
}

func TestHashNonByteAlignedOutput(t *testing.T) {
	seed := []byte("test application")
	in := []byte("input string")

	h, err := Hash(in, seed, 1, 5, 9)
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}
	if len(h) != 2 {
		t.Fatalf("unexpected hash length: got %d, want 2", len(h))
	}

	err = Verify(in, seed, 1, 5, 9, h)
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
}
