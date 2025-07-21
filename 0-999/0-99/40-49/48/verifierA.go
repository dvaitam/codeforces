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

func runBinary(bin, input string) (string, error) {
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

func expected(f, m, s string) string {
	beats := map[string]string{
		"rock":     "scissors",
		"scissors": "paper",
		"paper":    "rock",
	}
	gestures := []string{f, m, s}
	idx := -1
	common := ""
	switch {
	case f == m && f != s:
		idx = 2
		common = f
	case f == s && f != m:
		idx = 1
		common = f
	case m == s && m != f:
		idx = 0
		common = m
	default:
		return "?"
	}
	if beats[gestures[idx]] == common {
		switch idx {
		case 0:
			return "F"
		case 1:
			return "M"
		case 2:
			return "S"
		}
	}
	return "?"
}

func generateCase(r *rand.Rand) (string, string) {
	opts := []string{"rock", "paper", "scissors"}
	f := opts[r.Intn(3)]
	m := opts[r.Intn(3)]
	s := opts[r.Intn(3)]
	input := fmt.Sprintf("%s\n%s\n%s\n", f, m, s)
	return input, expected(f, m, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(r)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
