package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func solveCase(arr []int) string {
	n := len(arr)
	h := append([]int(nil), arr...)
	sort.Ints(h)
	if n == 2 {
		return fmt.Sprintf("%d %d\n", h[0], h[1])
	}
	idx := 0
	minDiff := h[1] - h[0]
	for i := 1; i < n-1; i++ {
		if d := h[i+1] - h[i]; d < minDiff {
			minDiff = d
			idx = i
		}
	}
	res := make([]int, 0, n)
	res = append(res, h[idx])
	for i := idx + 2; i < n; i++ {
		res = append(res, h[i])
	}
	for i := 0; i < idx; i++ {
		res = append(res, h[i])
	}
	res = append(res, h[idx+1])
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 2
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(100) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", n)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	expect := solveCase(arr)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), strings.TrimSpace(out), in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
