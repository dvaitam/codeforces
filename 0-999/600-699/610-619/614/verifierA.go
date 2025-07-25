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

func expectedOutput(l, r, k int64) string {
	if k == 1 {
		if l <= 1 && 1 <= r {
			return "1\n"
		}
		return "-1\n"
	}
	var res []string
	for x := int64(1); x <= r; {
		if x >= l {
			res = append(res, fmt.Sprintf("%d", x))
		}
		if x > r/k {
			break
		}
		x *= k
	}
	if len(res) == 0 {
		return "-1\n"
	}
	return strings.Join(res, " ") + "\n"
}

func generateCase(rng *rand.Rand) (string, string) {
	l := rng.Int63n(1_000_000_000_000_000) + 1
	r := l + rng.Int63n(1_000_000)
	var k int64
	if rng.Intn(10) == 0 {
		k = 1
	} else {
		k = rng.Int63n(1_000_000_000-1) + 2
	}
	input := fmt.Sprintf("%d %d %d\n", l, r, k)
	expect := expectedOutput(l, r, k)
	return input, expect
}

func runCase(bin string, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(expected)
	if got != want {
		return fmt.Errorf("expected %q got %q", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
