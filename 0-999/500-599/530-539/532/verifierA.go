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

func runBinary(bin string, input string) (string, error) {
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
	r := rand.New(rand.NewSource(1))
	tests := make([]string, 0, 100)
	for t := 0; t < 100; t++ {
		n := r.Intn(5) + 1
		var b strings.Builder
		fmt.Fprintf(&b, "%d\n", n)
		for i := 1; i <= n; i++ {
			fmt.Fprintf(&b, "%d ", r.Intn(10)+1)
		}
		b.WriteByte('\n')
		for i := 0; i < n-1; i++ {
			a := r.Intn(i+1) + 1
			bNode := i + 2
			fmt.Fprintf(&b, "%d %d\n", a, bNode)
		}
		k := r.Intn(n) + 1
		fmt.Fprintf(&b, "%d\n", k)
		for i := 0; i < k; i++ {
			fmt.Fprintf(&b, "%d ", r.Intn(10)+1)
		}
		b.WriteByte('\n')
		tests = append(tests, b.String())
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	cand := os.Args[1]
	official := "./officialA"
	if err := exec.Command("go", "build", "-o", official, "532A.go").Run(); err != nil {
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
