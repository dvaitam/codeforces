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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(s string) string {
	n := len(s)
	moves := (n - 11) / 2
	count8 := 0
	limit := n - 10
	for i := 0; i < limit && i < len(s); i++ {
		if s[i] == '8' {
			count8++
		}
	}
	if count8 > moves {
		return "YES"
	}
	return "NO"
}

func checkCase(bin string, s string) error {
	input := fmt.Sprintf("%d\n%s\n", len(s), s)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != expected(s) {
		return fmt.Errorf("expected %s got %s", expected(s), out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []string{"8", "12345678901", "88888888888"}
	for len(tests) < 100 {
		n := rng.Intn(9)*2 + 11 // odd length >=11
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteByte(byte('0' + rng.Intn(10)))
		}
		tests = append(tests, sb.String())
	}
	for i, s := range tests {
		if err := checkCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
