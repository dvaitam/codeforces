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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func computeB(l, r, x, y int64) int {
	if y%x != 0 {
		return 0
	}
	k := y / x
	cnt := 0
	for d := int64(1); d*d <= k; d++ {
		if k%d == 0 {
			m := d
			n := k / d
			if gcd(m, n) == 1 {
				a := x * m
				b := x * n
				if a >= l && a <= r && b >= l && b <= r {
					if m == n {
						cnt++
					} else {
						cnt += 2
					}
				}
			}
		}
	}
	return cnt
}

type testCaseB struct {
	input    string
	expected int
}

func generateCaseB(rng *rand.Rand) testCaseB {
	l := int64(rng.Intn(50) + 1)
	r := l + int64(rng.Intn(50))
	x := int64(rng.Intn(30) + 1)
	y := int64(rng.Intn(60) + 1)
	if x > y {
		x, y = y, x
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", l, r, x, y)
	return testCaseB{input: sb.String(), expected: computeB(l, r, x, y)}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		tc := generateCaseB(rng)
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, tc.input)
			os.Exit(1)
		}
		var val int
		if _, err := fmt.Sscan(out, &val); err != nil || val != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %s\ninput:\n%s", i, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
