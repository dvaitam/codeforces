package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	if !strings.Contains(bin, "/") {
		bin = "./" + bin
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func genTest() string {
	n := rand.Intn(99) + 2
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('A' + rand.Intn(26)))
	}
	return fmt.Sprintf("%d\n%s\n", n, sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s /path/to/binary\n", os.Args[0])
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "refB"
	if err := exec.Command("go", "build", "-o", ref, "1212B.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(1)
	for i := 0; i < 100; i++ {
		input := genTest()
		exp, err := runBinary(ref, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference failed:", err)
			os.Exit(1)
		}
		got, err := runBinary(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed to run: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Fprintf(os.Stderr, "test %d failed:\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
