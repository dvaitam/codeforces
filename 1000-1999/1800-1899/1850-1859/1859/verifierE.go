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

func runReference(input string) (string, error) {
	cmd := exec.Command("go", "run", "1859E.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func generateCase(rng *rand.Rand) (int, int, []int64, []int64) {
	n := rng.Intn(5) + 1
	k := rng.Intn(n) + 1
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(11) - 5)
		b[i] = int64(rng.Intn(11) - 5)
	}
	return n, k, a, b
}

func buildInput(n, k int, a, b []int64) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, n, k int, a, b []int64) error {
	input := buildInput(n, k, a, b)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expOut, err := runReference(input)
	if err != nil {
		return fmt.Errorf("reference failed: %v\n%s", err, expOut)
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expOut) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expOut), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, a, b := generateCase(rng)
		if err := runCase(bin, n, k, a, b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
