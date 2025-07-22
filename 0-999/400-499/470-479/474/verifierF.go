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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expectedF(arr []int, queries [][2]int) []int {
	res := make([]int, len(queries))
	for i, q := range queries {
		l, r := q[0]-1, q[1]-1
		g := arr[l]
		for j := l + 1; j <= r; j++ {
			g = gcd(g, arr[j])
		}
		cnt := 0
		for j := l; j <= r; j++ {
			if arr[j] == g {
				cnt++
			}
		}
		res[i] = (r - l + 1) - cnt
	}
	return res
}

func genCase(rng *rand.Rand) ([]int, [][2]int) {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(30) + 1
	}
	t := rng.Intn(10) + 1
	queries := make([][2]int, t)
	for i := range queries {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		queries[i] = [2]int{l, r}
	}
	return arr, queries
}

func runCase(bin string, arr []int, queries [][2]int, exp []int) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(arr))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", len(queries))
	for _, q := range queries {
		fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Fields(strings.TrimSpace(out.String()))
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(lines))
	}
	for i, l := range lines {
		var val int
		fmt.Sscan(l, &val)
		if val != exp[i] {
			return fmt.Errorf("query %d expected %d got %d", i+1, exp[i], val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		arr, queries := genCase(rng)
		exp := expectedF(arr, queries)
		if err := runCase(bin, arr, queries, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
