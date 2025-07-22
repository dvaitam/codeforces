package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func expectedOutput(n int, k int64, arr []int) string {
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = arr[i-1]
	}
	sort.Slice(a[1:], func(i, j int) bool { return a[i+1] < a[j+1] })
	var sb strings.Builder
	for i := 0; i <= n; i++ {
		for j := 1; j <= n-i; j++ {
			if k <= 0 {
				return strings.TrimSpace(sb.String())
			}
			k--
			sb.WriteString(strconv.Itoa(i + 1))
			for p := n - i + 1; p <= n; p++ {
				sb.WriteByte(' ')
				sb.WriteString(strconv.Itoa(a[p]))
			}
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(a[j]))
			sb.WriteByte('\n')
		}
	}
	return strings.TrimSpace(sb.String())
}

func runCase(bin string, n int, k int64, arr []int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprint(v))
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expected := expectedOutput(n, k, arr)
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 1
		maxK := n * (n + 1) / 2
		k := int64(rng.Intn(maxK) + 1)
		arr := rand.Perm(n)
		for j := range arr {
			arr[j] = arr[j] + 1 + rng.Intn(9)*10 + j // ensure distinct
		}
		if err := runCase(bin, n, k, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
