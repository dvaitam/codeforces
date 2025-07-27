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

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 2
	p := rng.Intn(3) + 1
	if p > 7 {
		p = 7
	}
	if p >= n {
		p = n - 1
	}
	k := rng.Intn(n-p) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, p, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(rng.Intn(20) + 1))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		for j := 0; j < p; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(rng.Intn(20) + 1))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
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
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "1316E.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := generateCase(rng)
		candOut, cErr := runBinary(candidate, input)
		refOut, rErr := runBinary(ref, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "case %d candidate error: %v\ninput:%s", t+1, cErr, input)
			os.Exit(1)
		}
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "case %d reference error: %v\ninput:%s", t+1, rErr, input)
			os.Exit(1)
		}
		if candOut != refOut {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:%sexpected:%sactual:%s\n", t+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
