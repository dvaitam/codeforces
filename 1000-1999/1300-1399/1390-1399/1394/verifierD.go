package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func generateTests() []string {
	r := rand.New(rand.NewSource(42))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(5) + 2
		var b strings.Builder
		fmt.Fprintf(&b, "%d\n", n)
		for j := 0; j < n; j++ {
			fmt.Fprintf(&b, "%d ", r.Int63n(100)+1)
		}
		b.WriteByte('\n')
		for j := 0; j < n; j++ {
			fmt.Fprintf(&b, "%d ", r.Intn(20)+1)
		}
		b.WriteByte('\n')
		for j := 2; j <= n; j++ {
			u := j
			v := r.Intn(j-1) + 1
			fmt.Fprintf(&b, "%d %d\n", u, v)
		}
		tests[i] = b.String()
	}
	return tests
}

func runCmd(cmdPath string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, cmdPath)
	c.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = io.Discard
	err := c.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	solPath := "./solD.bin"
	if err := exec.Command("go", "build", "-o", solPath, "1394D.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(solPath)
	bin := os.Args[1]

	tests := generateTests()
	for i, t := range tests {
		expect, err := runCmd(solPath, t)
		if err != nil {
			fmt.Printf("reference solution failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runCmd(bin, t)
		if err != nil {
			fmt.Printf("binary failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if expect != strings.TrimSpace(got) {
			fmt.Printf("mismatch on test %d\nexpected: %s\n got: %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
