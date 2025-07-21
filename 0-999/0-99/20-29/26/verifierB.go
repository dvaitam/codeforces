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

func expected(s string) int {
	open := 0
	ans := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			open++
		case ')':
			if open > 0 {
				open--
				ans += 2
			}
		}
	}
	return ans
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

func generateBracketSequence(length int, rng *rand.Rand) string {
	chars := make([]byte, length)
	for i := 0; i < length; i++ {
		if rng.Intn(2) == 0 {
			chars[i] = '('
		} else {
			chars[i] = ')'
		}
	}
	return string(chars)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		length := rng.Intn(1000) + 1
		s := generateBracketSequence(length, rng)
		
		input := s + "\n"
		
		expectedOut := expected(s)
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