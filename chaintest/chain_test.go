package chains

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"math/big"
	"math/rand"
	"testing"

	"golang.org/x/crypto/blake2b"
)

const (
	defaultSteps  = 4 * 24 * 30 // CENs every 15 mins, 30 days = 2880
	defaultSteps2 = 6 * 24 * 10 // CENs every 10 mins, 14 days = 1440
)

func DefaultModulus() *big.Int {
	t0 := new(big.Int)
	t1 := new(big.Int)
	t2 := new(big.Int)
	return new(big.Int).Add(t1.Sub(t0.Exp(big.NewInt(2), big.NewInt(256), nil), t2.Mul(t1.Exp(big.NewInt(2), big.NewInt(32), nil), big.NewInt(351))), big.NewInt(1))
}

// ForwardMiMC is fast
func ForwardMiMC(nsteps int, modulus *big.Int, roundConstants []*big.Int) *big.Int {
	t0 := new(big.Int)
	t1 := new(big.Int)
	t2 := new(big.Int)

	// Forward MiMC
	r := rand.New(rand.NewSource(99))
	input := big.NewInt(0).Rand(r, modulus)
	trace := new(big.Int).Set(input)
	for i := 1; i < nsteps; i++ {
		trace.Mod(t2.Add(t1.Mul(trace, t0.Mul(trace, trace)), roundConstants[i%len(roundConstants)]), modulus)
	}
	return new(big.Int).Set(trace)
}

// ReverseMiMC is slow
func ReverseMiMC(nsteps int, modulus *big.Int, roundConstants []*big.Int, input *big.Int, output *big.Int) bool {
	t1 := new(big.Int)
	t2 := new(big.Int)
	rtrace := new(big.Int).Set(output)
	littleFermatExpt := new(big.Int).Div(t2.Sub(t1.Mul(big.NewInt(2), modulus), big.NewInt(1)), big.NewInt(3))
	for i := nsteps - 1; i > 0; i-- {
		rtrace.Exp(t2.Sub(rtrace, roundConstants[i%len(roundConstants)]), littleFermatExpt, modulus)
	}
	return rtrace.Cmp(input) == 0
}

func GenRandomBytes(size int) (blk []byte, err error) {
	blk = make([]byte, size)
	_, err = rand.Read(blk)
	return
}

// Computehash returns the hash of its inputs
func Computehash(hasher hash.Hash, data ...[]byte) []byte {
	for _, b := range data {
		_, err := hasher.Write(b)
		if err != nil {
			panic(1)
		}
	}
	return hasher.Sum(nil)
}

func HashChain(nsteps int, hasher hash.Hash, b []byte) {
	for i := 1; i < nsteps; i++ {
		b = Computehash(hasher, b)
	}
}

func AESChain(nsteps int, cipher cipher.Block, b []byte) {
	cur := make([]byte, len(b))
	next := make([]byte, len(b))
	copy(cur[:], b[:])
	for i := 0; i < nsteps; i++ {
		cipher.Encrypt(next, cur)
		copy(cur[:], next[:])
		//fmt.Printf("%x\n", cur)
	}
}

func BenchmarkForwardMiMC(b *testing.B) {
	modulus := DefaultModulus()

	// 64 constants
	SEVEN := big.NewInt(7)
	FORTYTWO := big.NewInt(42)
	t2 := new(big.Int)
	roundConstants := make([]*big.Int, 64)
	for i := int64(0); i < 64; i++ {
		roundConstants[i] = new(big.Int).Xor(t2.Exp(big.NewInt(i), SEVEN, nil), FORTYTWO)
	}

	// run the MiMC function b.N times
	for n := 0; n < b.N; n++ {
		ForwardMiMC(defaultSteps, modulus, roundConstants)
	}
}

func BenchmarkSHA256(b *testing.B) {
	// run a SHA256 hash chain b.N times
	for n := 0; n < b.N; n++ {
		b0, _ := GenRandomBytes(32)
		HashChain(defaultSteps, sha256.New(), b0)
	}
}

func BenchmarkSHA1(b *testing.B) {
	// run a SHA1 chain b.N times
	for n := 0; n < b.N; n++ {
		b0, _ := GenRandomBytes(32)
		HashChain(defaultSteps, sha1.New(), b0)
	}
}

func BenchmarkBlake2b(b *testing.B) {
	// run a Blake2b chain b.N times
	for n := 0; n < b.N; n++ {
		b0, _ := GenRandomBytes(32)
		hasher, _ := blake2b.New256(b0)
		HashChain(defaultSteps, hasher, b0)
	}
}

func BenchmarkAES(b *testing.B) {
	// run the MiMC function b.N times
	for n := 0; n < b.N; n++ {
		b0, _ := GenRandomBytes(32)
		hasher, _ := aes.NewCipher(b0)
		AESChain(defaultSteps, hasher, b0)
	}
}

func TestAESChain(t *testing.T) {
	b0, _ := GenRandomBytes(32)
	hasher, _ := aes.NewCipher(b0)
	AESChain(defaultSteps, hasher, b0)
}
