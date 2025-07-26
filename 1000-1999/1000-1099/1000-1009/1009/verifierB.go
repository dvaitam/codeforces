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

func expectedAnswerB(s string) string {
	cnt1 := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '1' {
			cnt1++
		}
	}
	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != '1' {
			sb.WriteByte(s[i])
		}
	}
	s2 := sb.String()
	idx := strings.IndexByte(s2, '2')
	if idx == -1 {
		idx = len(s2)
	}
	prefix := s2[:idx]
	rest := s2[idx:]
	return prefix + strings.Repeat("1", cnt1) + rest
}

func genCaseB(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := range b {
		t := rng.Intn(3)
		b[i] = byte('0' + t)
	}
	return string(b)
}

func runCaseB(bin, s string) error {
	input := s + "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	exp := expectedAnswerB(s)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
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
	for i := 0; i < 100; i++ {
		s := genCaseB(rng)
		if err := runCaseB(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", i+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
