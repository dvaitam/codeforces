package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(s, t string) string {
	if t == "a" {
		return "1"
	}
	if strings.Contains(t, "a") {
		return "-1"
	}
	ans := int64(1)
	for i := 0; i < len(s); i++ {
		ans *= 2
	}
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(44))
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 100; i++ {
		sLen := rng.Intn(20) + 1
		s := strings.Repeat("a", sLen)
		tLen := rng.Intn(20) + 1
		var sb strings.Builder
		for j := 0; j < tLen; j++ {
			sb.WriteByte(letters[rng.Intn(26)])
		}
		t := sb.String()
		input := fmt.Sprintf("1\n%s\n%s\n", s, t)
		exp := expected(s, t)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
