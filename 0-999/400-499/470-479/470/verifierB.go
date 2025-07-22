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

func solveCase(s string) string {
	if strings.Contains(s, "666") {
		return "YES"
	}
	return "NO"
}

func generateCase(rng *rand.Rand) string {
	l := rng.Intn(100) + 1
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('0' + rng.Intn(10))
	}
	if l >= 3 && rng.Intn(2) == 0 {
		pos := rng.Intn(l - 2)
		b[pos] = '6'
		b[pos+1] = '6'
		b[pos+2] = '6'
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s := generateCase(rng)
		input := fmt.Sprintf("%s\n", s)
		expect := solveCase(s)
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
