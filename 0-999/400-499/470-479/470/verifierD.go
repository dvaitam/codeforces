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

func runCandidate(bin, input string) (string, error) {
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

func solveCase(k int, s string) string {
	res := make([]byte, len(s))
	shift := byte(k % 26)
	for i := 0; i < len(s); i++ {
		ch := s[i] - 'A'
		ch = (ch + shift) % 26
		res[i] = 'A' + ch
	}
	return string(res)
}

func generateCase(rng *rand.Rand) (int, string) {
	k := rng.Intn(26)
	l := rng.Intn(10) + 1
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('A' + rng.Intn(26))
	}
	return k, string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		k, s := generateCase(rng)
		input := fmt.Sprintf("%d\n%s\n", k, s)
		expect := solveCase(k, s)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
