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

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "838F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	X := rng.Intn(1000001)
	weights := make([]int, n+1)
	total := 0
	for i := 0; i <= n; i++ {
		weights[i] = rng.Intn(100) + 1
		total += weights[i]
	}
	p := make([]int, n+1)
	sum := 0
	for i := 0; i < n; i++ {
		p[i] = weights[i] * 1000000 / total
		sum += p[i]
	}
	p[n] = 1000000 - sum
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, X))
	for i := 0; i <= n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", p[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(exe, ref, input string) error {
	cmdRef := exec.Command(ref)
	cmdRef.Stdin = strings.NewReader(input)
	var refOut bytes.Buffer
	cmdRef.Stdout = &refOut
	cmdRef.Stderr = &refOut
	if err := cmdRef.Run(); err != nil {
		return fmt.Errorf("reference runtime error: %v\n%s", err, refOut.String())
	}
	expected := strings.TrimSpace(refOut.String())

	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected\n%s\ngot\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		if err := runCase(exe, ref, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
