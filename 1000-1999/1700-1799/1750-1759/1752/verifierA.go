package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildReference(dir string) (string, error) {
	refSource := "1752A.go"
	binPath := filepath.Join(dir, "ref_bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GO111MODULE=off")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v, %s", err, string(out))
	}
	return binPath, nil
}

func runBinary(bin string, input []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("run error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("failed to get working directory:", err)
		os.Exit(1)
	}
	refBin, err := buildReference(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rand.Seed(42)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(5) + 1  // 1..5
		k := n + rand.Intn(50) // n..n+49
		a := rand.Float64()
		b := rand.Float64()
		c := rand.Float64()
		d := rand.Float64()
		e := rand.Float64()
		P := rand.Intn(16) + 1 // 1..16

		input := fmt.Sprintf("%d\n%d\n%.2f %.2f\n%.2f %.2f\n%.2f %.2f %d\n", n, k, a, b, c, d, e, b, P)

		expect, err := runBinary(refBin, []byte(input))
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", t, err)
			os.Exit(1)
		}

		got, err := runBinary(candidate, []byte(input))
		if err != nil {
			fmt.Printf("candidate failed on test %d: %v\n", t, err)
			os.Exit(1)
		}

		if got != expect {
			fmt.Printf("wrong answer on test %d\nexpected:\n%s\ngot:\n%s\n", t, expect, got)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}
