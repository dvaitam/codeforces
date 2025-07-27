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

func run(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	refBin := "./1425A_ref"
	if err := exec.Command("go", "build", "-o", refBin, "1425A.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	rand.Seed(42)
	for i := 0; i < 100; i++ {
		n := rand.Int63n(1e12) + 1
		input := fmt.Sprintf("1\n%d\n", n)
		expect, err := run(refBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference failed on test", i+1, ":", err)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "program failed on test", i+1, ":", err)
			os.Exit(1)
		}
		if expect != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %s got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
