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

func expected(n int) int {
	spf := make([]int, n+1)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			for j := i; j <= n; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
	
	count := 0
	for i := 2; i <= n; i++ {
		x := i
		prev := 0
		distinct := 0
		for x > 1 {
			p := spf[x]
			if p != prev {
				distinct++
				prev = p
			}
			x /= p
		}
		if distinct == 2 {
			count++
		}
	}
	return count
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		n := rng.Intn(3000) + 1
		
		input := fmt.Sprintf("%d\n", n)
		
		expectedOut := expected(n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		
		gotInt, parseErr := strconv.Atoi(got)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output: %v\ninput:\n%s", i+1, parseErr, input)
			os.Exit(1)
		}
		
		if gotInt != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, expectedOut, gotInt, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}