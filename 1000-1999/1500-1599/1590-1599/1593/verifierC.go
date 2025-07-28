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

func expected(n int64, xs []int64) string {
	sort.Slice(xs, func(i, j int) bool { return xs[i] > xs[j] })
	var dist int64
	saved := 0
	for _, x := range xs {
		need := n - x
		if dist+need < n {
			saved++
			dist += need
		} else {
			break
		}
	}
	return fmt.Sprintf("%d", saved)
}

func generateCase(rng *rand.Rand) (int64, []int64) {
	n := rng.Int63n(1000) + 2
	k := rng.Intn(7) + 1
	xs := make([]int64, k)
	for i := 0; i < k; i++ {
		xs[i] = rng.Int63n(n-1) + 1
	}
	return n, xs
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
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
		n, xs := generateCase(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		fmt.Fprintf(&sb, "%d %d\n", n, len(xs))
		for j, x := range xs {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", x)
		}
		sb.WriteByte('\n')
		input := sb.String()
		expectedOutput := expected(n, append([]int64(nil), xs...))
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expectedOutput {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:%s", i+1, expectedOutput, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
