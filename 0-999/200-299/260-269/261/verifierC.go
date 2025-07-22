package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const MAX = 64

var comb [MAX][MAX]uint64

func initComb() {
	for i := 0; i < MAX; i++ {
		comb[i][0] = 1
		for j := 1; j <= i; j++ {
			comb[i][j] = comb[i-1][j-1] + comb[i-1][j]
		}
	}
}

func count(n uint64, t uint64) uint64 {
	if t == 0 || t&(t-1) != 0 {
		return 0
	}
	k := bits.TrailingZeros64(t)
	Ln := 64 - bits.LeadingZeros64(n)
	var ans uint64
	for L := k + 1; L < Ln; L++ {
		ans += comb[L-1][k]
	}
	zeros := 0
	for i := Ln - 2; i >= 0; i-- {
		if (n>>uint(i))&1 == 1 {
			remZeros := k - (zeros + 1)
			remBits := i
			if remZeros >= 0 && remZeros <= remBits {
				ans += comb[remBits][remZeros]
			}
		}
		if (n>>uint(i))&1 == 0 {
			zeros++
		}
	}
	if zeros == k {
		ans++
	}
	return ans
}

func genCaseC(rng *rand.Rand) (string, string) {
	k := rng.Intn(6)
	t := uint64(1) << k
	n := uint64(rng.Int63n(1<<20) + int64(t))
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, t)
	exp := fmt.Sprintf("%d\n", count(n, t))
	return sb.String(), exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", strings.TrimSpace(expected), got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	initComb()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseC(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
