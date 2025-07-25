package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	ref := filepath.Join(dir, "refE")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "534E.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for t := 0; t < 100; t++ {
		n := rng.Intn(10) + 2
		a := make([]int, n)
		cur := 0
		for i := 0; i < n; i++ {
			cur += rng.Intn(5) + 1
			a[i] = cur
		}
		m := rng.Intn(10) + 1
		b := make([]int, m)
		for i := 0; i < m; i++ {
			b[i] = rng.Intn(n) + 1
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(a[i]))
		}
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for i := 0; i < m; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(b[i]))
		}
		sb.WriteString("\n")
		input := sb.String()
		candOut, cErr := runBinary(candidate, input)
		refOut, rErr := runBinary(ref, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate error: %v\n", t+1, cErr)
			os.Exit(1)
		}
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: reference error: %v\n", t+1, rErr)
			os.Exit(1)
		}
		if candOut != refOut {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:%sexpected:%sactual:%s\n", t+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
