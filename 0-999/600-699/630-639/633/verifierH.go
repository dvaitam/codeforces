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

func solveH(n, m int, arr []int64, queries [][2]int) string {
	fib := make([]int64, n+2)
	if n >= 1 {
		fib[1] = int64(1 % m)
	}
	if n >= 2 {
		fib[2] = int64(1 % m)
	}
	for i := 3; i <= n+1; i++ {
		fib[i] = (fib[i-1] + fib[i-2]) % int64(m)
	}
	var sb strings.Builder
	for _, qr := range queries {
		l := qr[0] - 1
		r := qr[1] - 1
		uniq := make(map[int64]struct{})
		for i := l; i <= r; i++ {
			uniq[arr[i]] = struct{}{}
		}
		values := make([]int64, 0, len(uniq))
		for k := range uniq {
			values = append(values, k)
		}
		sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
		res := int64(0)
		for i, v := range values {
			idx := i + 1
			for idx >= len(fib) {
				fib = append(fib, (fib[len(fib)-1]+fib[len(fib)-2])%int64(m))
			}
			res = (res + (v%int64(m))*fib[idx]%int64(m)) % int64(m)
		}
		sb.WriteString(fmt.Sprintf("%d\n", res%int64(m)))
	}
	return sb.String()
}

func generateCaseH(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(10) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(20))
	}
	q := rng.Intn(4) + 1
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		queries[i] = [2]int{l, r}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for _, qr := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", qr[0], qr[1]))
	}
	expect := solveH(n, m, arr, queries)
	return sb.String(), expect
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
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCaseH(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
