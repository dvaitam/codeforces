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

type testCase struct {
	n int
	a int
	b int
	s string
}

func solveCase(tc testCase) int {
	a := tc.a
	b := tc.b
	prev := 0
	ans := 0
	for i := 0; i < tc.n; i++ {
		if tc.s[i] == '*' {
			prev = 0
			continue
		}
		if prev == 1 {
			if b > 0 {
				b--
				ans++
				prev = 2
			} else {
				prev = 0
			}
		} else if prev == 2 {
			if a > 0 {
				a--
				ans++
				prev = 1
			} else {
				prev = 0
			}
		} else {
			if a >= b {
				if a > 0 {
					a--
					ans++
					prev = 1
				} else if b > 0 {
					b--
					ans++
					prev = 2
				} else {
					prev = 0
				}
			} else {
				if b > 0 {
					b--
					ans++
					prev = 2
				} else if a > 0 {
					a--
					ans++
					prev = 1
				} else {
					prev = 0
				}
			}
		}
	}
	return ans
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	a := rng.Intn(n + 1)
	b := rng.Intn(n + 1)
	bs := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			bs[i] = '.'
		} else {
			bs[i] = '*'
		}
	}
	s := string(bs)
	input := fmt.Sprintf("%d %d %d\n%s\n", n, a, b, s)
	tc := testCase{n, a, b, s}
	out := solveCase(tc)
	expected := fmt.Sprintf("%d\n", out)
	return input, expected
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
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
