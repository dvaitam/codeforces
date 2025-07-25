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

func solveCase(input string) string {
	reader := strings.NewReader(input)
	var n int
	fmt.Fscan(reader, &n)
	sumA, sumB, sumC := int64(0), int64(0), int64(0)
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		sumA += x
	}
	for i := 0; i < n-1; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		sumB += x
	}
	for i := 0; i < n-2; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		sumC += x
	}
	missing1 := sumA - sumB
	missing2 := sumB - sumC
	return fmt.Sprintf("%d %d", missing1, missing2)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	a := make([]int64, n)
	for i := range a {
		a[i] = rng.Int63n(20)
	}
	b := append([]int64(nil), a...)
	idx := rng.Intn(len(b))
	b = append(b[:idx], b[idx+1:]...)
	c := append([]int64(nil), b...)
	idx2 := rng.Intn(len(c))
	c = append(c[:idx2], c[idx2+1:]...)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range c {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		inp := generateCase(rng)
		exp := solveCase(inp)
		if err := runCase(bin, inp, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, inp)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
