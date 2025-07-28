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

type pair struct{ u, v int }

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int, a []int) string {
	x := 0
	for i := 1; i <= n; i++ {
		if a[i-1] == 0 {
			x = i
		}
	}
	var sb strings.Builder
	if x == 0 {
		fmt.Fprintf(&sb, "%d ", n+1)
	}
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&sb, "%d ", i)
		if i == x {
			fmt.Fprintf(&sb, "%d ", n+1)
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(2)
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), expected(n, arr)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, want := generateCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
