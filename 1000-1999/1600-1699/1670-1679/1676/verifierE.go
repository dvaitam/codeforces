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

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	q := rng.Intn(10) + 1
	arr := make([]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, q)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(100) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	queries := make([]int, q)
	for i := 0; i < q; i++ {
		queries[i] = rng.Intn(n*100 + 1)
		fmt.Fprintf(&sb, "%d\n", queries[i])
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
	prefix := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		sum += arr[i]
		prefix[i] = sum
	}
	var out strings.Builder
	for _, x := range queries {
		idx := sort.Search(len(prefix), func(i int) bool { return prefix[i] >= x })
		if idx == len(prefix) {
			out.WriteString("-1\n")
		} else {
			fmt.Fprintf(&out, "%d\n", idx+1)
		}
	}
	return sb.String(), strings.TrimSpace(out.String())
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected \n%s\ngot \n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
