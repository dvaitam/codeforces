package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(a, b int64) int64 {
	count := int64(0)
	for l := 2; l <= 60; l++ {
		allOnes := (int64(1) << uint(l)) - 1
		for i := 0; i < l-1; i++ {
			num := allOnes - (int64(1) << uint(i))
			if num >= a && num <= b {
				count++
			}
		}
	}
	return count
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(2))
	for i := 0; i < 100; i++ {
		var a, b int64
		a = rng.Int63n(1<<60) + 1
		b = a + rng.Int63n(1<<20)
		if rng.Intn(4) == 0 {
			a = 1
			b = (1 << 60) - 1
		}
		if a > b {
			a, b = b, a
		}
		input := fmt.Sprintf("%d %d\n", a, b)
		want := expected(a, b)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil || got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", i+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
