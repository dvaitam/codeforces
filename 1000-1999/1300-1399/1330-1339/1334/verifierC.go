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

func expected(a, b []int64) string {
	n := len(a)
	extra := make([]int64, n)
	var sum int64
	for i := 0; i < n; i++ {
		prev := (i + n - 1) % n
		if a[i] > b[prev] {
			extra[i] = a[i] - b[prev]
		} else {
			extra[i] = 0
		}
		sum += extra[i]
	}
	ans := sum + a[0] - extra[0]
	for i := 1; i < n; i++ {
		cand := sum + a[i] - extra[i]
		if cand < ans {
			ans = cand
		}
	}
	return fmt.Sprintf("%d", ans)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		a := make([]int64, n)
		bArr := make([]int64, n)
		for j := 0; j < n; j++ {
			a[j] = int64(rng.Intn(20) + 1)
			bArr[j] = int64(rng.Intn(20) + 1)
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", a[j], bArr[j]))
		}
		input := sb.String()
		exp := expected(append([]int64(nil), a...), append([]int64(nil), bArr...))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
