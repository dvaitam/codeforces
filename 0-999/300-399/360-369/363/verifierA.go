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

func expected(n int) string {
	var digits []int
	if n == 0 {
		digits = append(digits, 0)
	} else {
		for n > 0 {
			digits = append(digits, n%10)
			n /= 10
		}
	}
	var sb strings.Builder
	for i, d := range digits {
		if d >= 5 {
			sb.WriteString("-O|")
		} else {
			sb.WriteString("O-|")
		}
		x := d % 5
		for j := 0; j < x; j++ {
			sb.WriteByte('O')
		}
		sb.WriteByte('-')
		for j := 0; j < 4-x; j++ {
			sb.WriteByte('O')
		}
		if i+1 < len(digits) {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) int {
	return rng.Intn(1000000000)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := generateCase(rng)
		input := fmt.Sprintf("%d\n", n)
		want := expected(n)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\n\ngot\n%s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
