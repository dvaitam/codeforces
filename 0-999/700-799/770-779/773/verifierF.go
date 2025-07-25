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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	refBin := "./773F_ref"
	if err := exec.Command("go", "build", "-o", refBin, "773F.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	rand.Seed(42)
	for t := 0; t < 100; t++ {
		maxn := rand.Intn(10) + 1
		maxa := maxn + rand.Intn(100)
		q := rand.Intn(1000) + 104
		input := fmt.Sprintf("%d %d %d\n", maxn, maxa, q)
		expect, err := run(refBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference failed on test", t+1, ":", err)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "program failed on test", t+1, ":", err)
			os.Exit(1)
		}
		if expect != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %s got %s\n", t+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
