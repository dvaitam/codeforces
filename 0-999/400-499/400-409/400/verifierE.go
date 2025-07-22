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

func bruteSum(a []int) int64 {
	n := len(a)
	total := int64(0)
	cur := make([]int, n)
	copy(cur, a)
	for len(cur) > 0 {
		for _, v := range cur {
			total += int64(v)
		}
		if len(cur) == 1 {
			break
		}
		next := make([]int, len(cur)-1)
		for i := 0; i < len(cur)-1; i++ {
			next[i] = cur[i] & cur[i+1]
		}
		cur = next
	}
	return total
}

func expected(n, m int, arr []int, queries [][2]int) string {
	a := make([]int, n)
	copy(a, arr)
	var sb strings.Builder
	for i := 0; i < m; i++ {
		p := queries[i][0]
		v := queries[i][1]
		a[p-1] = v
		s := bruteSum(a)
		sb.WriteString(fmt.Sprintf("%d\n", s))
	}
	return strings.TrimSpace(sb.String())
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		n := rng.Intn(6) + 1
		m := rng.Intn(5) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(32)
		}
		queries := make([][2]int, m)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		sb.WriteByte('\n')
		for i := 0; i < m; i++ {
			p := rng.Intn(n) + 1
			v := rng.Intn(32)
			queries[i] = [2]int{p, v}
			sb.WriteString(fmt.Sprintf("%d %d\n", p, v))
		}
		input := sb.String()
		expectedOut := expected(n, m, arr, queries)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tcase+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", tcase+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
