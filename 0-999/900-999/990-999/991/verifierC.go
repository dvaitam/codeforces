package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func good(n, k int64) bool {
	candies := n
	eaten := int64(0)
	for candies > 0 {
		if candies < k {
			eaten += candies
			break
		}
		eaten += k
		candies -= k
		petya := candies / 10
		candies -= petya
	}
	return eaten*2 >= n
}

func expectedC(n int64) int64 {
	l, r := int64(1), n
	for l < r {
		mid := (l + r) / 2
		if good(n, mid) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}

func generateCaseC(rng *rand.Rand) (string, int64) {
	n := rng.Int63n(1_000_000_000_000) + 1 // up to 1e12
	input := fmt.Sprintf("%d\n", n)
	expected := expectedC(n)
	return input, expected
}

func runCaseC(bin, input string, expected int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseC(rng)
		if err := runCaseC(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
