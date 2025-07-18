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

func generateCase(rng *rand.Rand) (string, int, int) {
	n := rng.Intn(100) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '('
		} else {
			b[i] = ')'
		}
	}
	s := string(b)
	l, c := longestRegular(s)
	return s + "\n", l, c
}

func longestRegular(s string) (int, int) {
	stack := []int{-1}
	best := 0
	count := 1
	for i, ch := range s {
		if ch == '(' {
			stack = append(stack, i)
		} else {
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				stack = append(stack, i)
			} else {
				length := i - stack[len(stack)-1]
				if length == best {
					count++
				} else if length > best {
					best = length
					count = 1
				}
			}
		}
	}
	if best == 0 {
		return 0, 1
	}
	return best, count
}

func runCase(bin, input string, expLen, expCnt int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var gotL, gotC int
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &gotL, &gotC); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	if gotL != expLen || gotC != expCnt {
		return fmt.Errorf("expected %d %d got %d %d", expLen, expCnt, gotL, gotC)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, l, c := generateCase(rng)
		if err := runCase(bin, in, l, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
