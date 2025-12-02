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

// Go implementation of the solution logic (Segment Tree + Fenwick Tree)
func solveOracle(n int, a []int, queries [][2]int) []int {
	// Adjust to 1-based indexing for internal logic consistency with the C++ reference
	// a is passed as 0-based from generator, convert to 1-based
	arr := make([]int, n+1)
	for i := 0; i < n; i++ {
		arr[i+1] = a[i]
	}
	
	// Constants
	// N is n.
	
	// Fenwick Tree
	tree := make([]int, n+1)
	add := func(x int) {
		for x > 0 {
			tree[x]++
			x -= x & -x
		}
	}
	ask := func(x int) int {
		s := 0
		for x <= n {
			s += tree[x]
			x += x & -x
		}
		return s
	}

	// Segment Tree
	mx := make([]int, 4*n+1)
	
	var upd func(k, p, l, r int)
	upd = func(k, p, l, r int) {
		if l == r {
			mx[p] = arr[k]
			return
		}
		mid := (l + r) / 2
		if k <= mid {
			upd(k, p*2, l, mid)
		} else {
			upd(k, p*2+1, mid+1, r)
		}
		v1 := mx[p*2]
		v2 := mx[p*2+1]
		if v1 > v2 {
			mx[p] = v1
		} else {
			mx[p] = v2
		}
	}

	var querySeg func(L, R, p, l, r int) int
	querySeg = func(L, R, p, l, r int) int {
		if L <= l && r <= R {
			return mx[p]
		}
		mid := (l + r) / 2
		s := 0
		if L <= mid {
			v := querySeg(L, R, p*2, l, mid)
			if v > s {
				s = v
			}
		}
		if R > mid {
			v := querySeg(L, R, p*2+1, mid+1, r)
			if v > s {
				s = v
			}
		}
		return s
	}

	// Preprocess positions
	// Array values are up to n.
	posList := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		v := arr[i]
		posList[v] = append(posList[v], i)
	}

	qRes := make([][]int, n+1)
	qL := make([]int, len(queries))
	
	for i, q := range queries {
		l, r := q[0], q[1]
		qL[i] = l
		qRes[r] = append(qRes[r], i)
	}

	ans := make([]int, len(queries))

	for t := 1; t <= n; t++ {
		y := 0
		if t < len(posList) {
			for _, x := range posList[t] {
				var val int
				if y == 0 {
					val = t
				} else {
					val = querySeg(y, x, 1, 1, n)
				}
				add(val)
				y = x
			}
			for _, x := range posList[t] {
				upd(x, 1, 1, n)
			}
		}
		if t < len(qRes) {
			for _, qIdx := range qRes[t] {
				ans[qIdx] = ask(qL[qIdx])
			}
		}
	}

	return ans
}

func expectedF(n int, arr []int, queries [][2]int) string {
	res := solveOracle(n, arr, queries)
	var sb strings.Builder
	for _, v := range res {
		sb.WriteString(fmt.Sprintf("%d\n", v))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	q := rng.Intn(50) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(n) + 1 // 1 <= a_i <= n
	}
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l // l <= r <= n
		queries[i] = [2]int{l, r}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	for _, qu := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", qu[0], qu[1]))
	}
	expect := expectedF(n, arr, queries)
	return sb.String(), expect
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected:\n%q\ngot:\n%q", expStr, outStr)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
