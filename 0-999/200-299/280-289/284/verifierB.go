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

func expectedResult(s string) int {
	cntI := 0
	cntF := 0
	for _, c := range s {
		if c == 'I' {
			cntI++
		} else if c == 'F' {
			cntF++
		}
	}
	if cntI == 0 {
		return len(s) - cntF
	}
	if cntI == 1 {
		return 1
	}
	return 0
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(100) + 2 // n >= 2
	var sb strings.Builder
	for i := 0; i < n; i++ {
		v := rng.Intn(3)
		switch v {
		case 0:
			sb.WriteByte('A')
		case 1:
			sb.WriteByte('I')
		default:
			sb.WriteByte('F')
		}
	}
	s := sb.String()
	input := fmt.Sprintf("%d\n%s\n", n, s)
	return input, expectedResult(s)
}

func runCase(bin, input string, expected int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
