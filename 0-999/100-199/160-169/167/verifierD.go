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

const mod int64 = 1e9 + 9

func solve(n, k int, xs, ys []int64, a, b, c, d int64, queries [][2]int64) []int {
	for i := k; i < n; i++ {
		xs[i] = (a*xs[i-1] + b) % mod
		ys[i] = (c*ys[i-1] + d) % mod
	}
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool { return xs[idx[i]] < xs[idx[j]] })
	sortedX := make([]int64, n)
	for i, id := range idx {
		sortedX[i] = xs[id]
	}
	res := make([]int, len(queries))
	for qi, q := range queries {
		L, R := q[0], q[1]
		l := sort.Search(n, func(i int) bool { return sortedX[i] >= L })
		r := sort.Search(n, func(i int) bool { return sortedX[i] > R }) - 1
		cnt := 0
		if l < n && r >= l {
			cnt = r - l + 1
		}
		res[qi] = cnt / 2
	}
	return res
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(8) + 2
	k := rng.Intn(n) + 1
	xs := make([]int64, n)
	ys := make([]int64, n)
	for i := 0; i < k; i++ {
		xs[i] = rng.Int63n(mod)
		ys[i] = rng.Int63n(mod)
	}
	a := rng.Int63n(10) + 1
	b := rng.Int63n(10) + 1
	c := rng.Int63n(10) + 1
	d := rng.Int63n(10) + 1
	m := rng.Intn(5) + 1
	queries := make([][2]int64, m)
	for i := 0; i < m; i++ {
		L := rng.Int63n(mod)
		R := L + rng.Int63n(mod-L)
		queries[i] = [2]int64{L, R}
	}
	ans := solve(n, k, append([]int64(nil), xs...), append([]int64(nil), ys...), a, b, c, d, queries)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < k; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", xs[i], ys[i]))
	}
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", a, b, c, d))
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 0; i < m; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", queries[i][0], queries[i][1]))
	}
	return sb.String(), ans
}

func run(bin, input string) ([]int, error) {
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
		return nil, fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	fields := strings.Fields(out.String())
	res := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid int on line %d", i+1)
		}
		res[i] = v
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if len(out) != len(exp) {
			fmt.Fprintf(os.Stderr, "case %d wrong number of lines\ninput:\n%s", i+1, in)
			os.Exit(1)
		}
		for j := range exp {
			if out[j] != exp[j] {
				fmt.Fprintf(os.Stderr, "case %d mismatch on line %d: expected %d got %d\ninput:\n%s", i+1, j+1, exp[j], out[j], in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
