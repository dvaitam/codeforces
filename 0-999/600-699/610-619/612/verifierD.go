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

func runCmd(exe, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	cmd := exec.CommandContext(ctx, exe)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func buildRef() (string, error) {
	refBin := "./refD_bin"
	cmd := exec.Command("go", "build", "-o", refBin, "612D.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return refBin, nil
}

func generateTests() []string {
	r := rand.New(rand.NewSource(42))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(8) + 1
		k := r.Intn(n) + 1
		sb := strings.Builder{}
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for j := 0; j < n; j++ {
			l := r.Intn(21) - 10
			r2 := l + r.Intn(11)
			fmt.Fprintf(&sb, "%d %d\n", l, r2)
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	refBin, err := buildRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, input := range tests {
		expect, err := runCmd(refBin, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCmd(candidate, input)
		if err != nil {
			fmt.Printf("candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(expect) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", idx+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
