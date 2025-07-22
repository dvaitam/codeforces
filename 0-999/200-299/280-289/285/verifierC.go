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

func runCandidate(bin, input string) (string, error) {
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

func expected(a []int64) string {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	var sum int64
	for i, v := range a {
		idx := int64(i + 1)
		if v > idx {
			sum += v - idx
		} else {
			sum += idx - v
		}
	}
	return fmt.Sprintf("%d", sum)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rng.Int63n(2000) - 1000
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	expect := expected(append([]int64(nil), arr...))
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	edgeCases := [][]int64{
		{1},
		{1, 2, 3},
		{3, 2, 1},
		{-5, 5, 0},
		{10, -10},
	}
	for _, arr := range edgeCases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		exp := expected(append([]int64(nil), arr...))
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "edge case failed: %v\ninput:\n%s", err, sb.String())
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "edge case failed: expected %s got %s\ninput:\n%s", exp, out, sb.String())
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
