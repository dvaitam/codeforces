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

func solveE(s string) int {
	cur := 0
	maxPeople := 0
	for _, c := range s {
		if c == '+' {
			cur++
			if cur > maxPeople {
				maxPeople = cur
			}
		} else if c == '-' {
			if cur > 0 {
				cur--
			} else {
				maxPeople++
			}
		}
	}
	return maxPeople
}

func generateCase(rng *rand.Rand) (string, string) {
	l := rng.Intn(50) + 1
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '+'
		} else {
			b[i] = '-'
		}
	}
	s := string(b)
	return s + "\n", fmt.Sprint(solveE(s))
}

func runCase(bin, input, expected string) error {
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
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
