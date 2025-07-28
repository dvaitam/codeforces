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

type testCase struct {
	rating int
}

func solveCase(tc testCase) string {
	switch {
	case tc.rating >= 1900:
		return "Division 1"
	case tc.rating >= 1600:
		return "Division 2"
	case tc.rating >= 1400:
		return "Division 3"
	default:
		return "Division 4"
	}
}

func genCase(rng *rand.Rand) (string, string) {
	tc := testCase{rating: rng.Intn(10001) - 5000}
	input := fmt.Sprintf("1\n%d\n", tc.rating)
	output := solveCase(tc) + "\n"
	return input, output
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(exp), got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
