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

type query struct {
	i int
	j int64
}

func permutations(p []int, c int) [][]int {
	n := len(p)
	set := make(map[string]struct{})
	var res [][]int
	var rec func(start int, remain int)
	rec = func(start int, remain int) {
		key := fmt.Sprint(p)
		if _, ok := set[key]; !ok {
			set[key] = struct{}{}
			arr := make([]int, n)
			copy(arr, p)
			res = append(res, arr)
		}
		for l := start; l < n; l++ {
			for r := l + 1; r < n; r++ {
				cost := r - l
				if cost <= remain {
					for x, y := l, r; x < y; x, y = x+1, y-1 {
						p[x], p[y] = p[y], p[x]
					}
					rec(r+1, remain-cost)
					for x, y := l, r; x < y; x, y = x+1, y-1 {
						p[x], p[y] = p[y], p[x]
					}
				}
			}
		}
	}
	rec(0, c)
	sort.Slice(res, func(i, j int) bool {
		a, b := res[i], res[j]
		for k := 0; k < n; k++ {
			if a[k] != b[k] {
				return a[k] < b[k]
			}
		}
		return false
	})
	return res
}

func solveCase(p []int, c int, qs []query) []int {
	perms := permutations(p, c)
	ans := make([]int, len(qs))
	for i, q := range qs {
		idx := q.i - 1
		ans[i] = perms[idx][q.j-1]
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	p := rng.Perm(n)
	for i := 0; i < n; i++ {
		p[i]++
	}
	c := rng.Intn(3) + 1
	q := rng.Intn(3) + 1
	qs := make([]query, q)
	for i := 0; i < q; i++ {
		qs[i].i = rng.Intn(3) + 1
		qs[i].j = int64(rng.Intn(n) + 1)
	}
	input := fmt.Sprintf("1\n%d %d %d\n", n, c, q)
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", p[i])
	}
	input += "\n"
	for _, qv := range qs {
		input += fmt.Sprintf("%d %d\n", qv.i, qv.j)
	}
	ans := solveCase(p, c, qs)
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return input, sb.String()
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
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
