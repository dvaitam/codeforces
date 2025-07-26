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

func sumDigits(x int64) int64 {
	var s int64
	for x > 0 {
		s += x % 10
		x /= 10
	}
	return s
}

func expected(n, s int64) int64 {
	low, high := int64(1), n
	for low <= high {
		mid := (low + high) / 2
		if mid-sumDigits(mid) >= s {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	if low > n {
		return 0
	}
	return n - low + 1
}

func runBin(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Int63n(1_000_000_000_000) + 1
	s := rng.Int63n(n) + 1
	input := fmt.Sprintf("%d %d\n", n, s)
	exp := fmt.Sprintf("%d", expected(n, s))
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
