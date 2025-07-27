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

func genCase() string {
	n := rand.Intn(5) + 1
	m := rand.Intn(n) + 1
	r := rand.Intn(3)
	used := make(map[[2]int]bool)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, r))
	for i := 0; i < n; i++ {
		for {
			x := rand.Intn(10) + 1
			y := rand.Intn(10) + 1
			if !used[[2]int{x, y}] {
				used[[2]int{x, y}] = true
				b := rand.Intn(10) + 1
				sb.WriteString(fmt.Sprintf("%d %d %d\n", x, y, b))
				break
			}
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	refBin := "./1425D_ref"
	if err := exec.Command("go", "build", "-o", refBin, "1425D.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	rand.Seed(42)
	for i := 0; i < 100; i++ {
		input := genCase()
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
