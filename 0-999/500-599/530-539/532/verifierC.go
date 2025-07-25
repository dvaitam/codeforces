package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []string {
	r := rand.New(rand.NewSource(3))
	tests := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		xp := r.Intn(20)
		yp := r.Intn(20)
		xv := r.Intn(20)
		yv := r.Intn(20)
		tests = append(tests, fmt.Sprintf("%d %d %d %d\n", xp, yp, xv, yv))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	cand := os.Args[1]
	official := "./officialC"
	if err := exec.Command("go", "build", "-o", official, "532C.go").Run(); err != nil {
		fmt.Println("failed to build official solution:", err)
		os.Exit(1)
	}
	defer os.Remove(official)
	tests := generateTests()
	for i, tc := range tests {
		exp, eerr := runBinary(official, tc)
		got, gerr := runBinary(cand, tc)
		if eerr != nil {
			fmt.Printf("official solution failed on test %d: %v\n", i+1, eerr)
			os.Exit(1)
		}
		if gerr != nil {
			fmt.Printf("candidate failed on test %d: %v\n", i+1, gerr)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
