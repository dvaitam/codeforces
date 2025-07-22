package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func trailingNines(x int64) int {
	cnt := 0
	for x%10 == 9 {
		cnt++
		x /= 10
	}
	return cnt
}

func expectedAnswerB(p, d int64) int64 {
	best := p
	bestK := trailingNines(p)
	pow10 := make([]int64, 19)
	pow10[0] = 1
	for i := 1; i < 19; i++ {
		pow10[i] = pow10[i-1] * 10
	}
	for k := bestK + 1; k < 19; k++ {
		mod := pow10[k]
		if p < mod {
			break
		}
		cand := (p/mod)*mod - 1
		if cand >= 0 && p-cand <= d {
			best = cand
			bestK = k
		}
	}
	return best
}

func generateCaseB(rng *rand.Rand) (int64, int64) {
	p := rng.Int63n(1_000_000_000_000) + 1 // up to 1e12
	d := rng.Int63n(p)
	return p, d
}

func runCaseB(bin string, p, d int64) error {
	input := fmt.Sprintf("%d %d\n", p, d)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprint(expectedAnswerB(p, d))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// deterministic test
	if err := runCaseB(bin, 99, 10); err != nil {
		fmt.Fprintln(os.Stderr, "deterministic case failed:", err)
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {
		p, d := generateCaseB(rng)
		if err := runCaseB(bin, p, d); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d\n", i+1, err, p, d)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
