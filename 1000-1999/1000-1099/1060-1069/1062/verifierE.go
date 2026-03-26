package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ─── embedded correct solver ───

func solveE(input string) string {
	data := []byte(input)
	idx := 0
	nextInt := func() int {
		for idx < len(data) && (data[idx] < '0' || data[idx] > '9') {
			idx++
		}
		val := 0
		for idx < len(data) && data[idx] >= '0' && data[idx] <= '9' {
			val = val*10 + int(data[idx]-'0')
			idx++
		}
		return val
	}

	n := nextInt()
	q := nextInt()

	parent := make([]int, n+1)
	childCount := make([]int, n+1)
	for i := 2; i <= n; i++ {
		p := nextInt()
		parent[i] = p
		childCount[p]++
	}

	children := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		if childCount[i] > 0 {
			children[i] = make([]int, 0, childCount[i])
		}
	}
	for i := 2; i <= n; i++ {
		p := parent[i]
		children[p] = append(children[p], i)
	}

	K := 0
	for (1 << K) <= n {
		K++
	}

	tin := make([]int, n+1)
	depth := make([]int, n+1)
	up := make([][]int, K)
	for k := 0; k < K; k++ {
		up[k] = make([]int, n+1)
		up[k][1] = 1
	}

	timer := 1
	tin[1] = 1
	stack := make([]int, 1, n)
	stack[0] = 1
	ptr := make([]int, n+1)

	for len(stack) > 0 {
		u := stack[len(stack)-1]
		if ptr[u] < len(children[u]) {
			v := children[u][ptr[u]]
			ptr[u]++
			depth[v] = depth[u] + 1
			up[0][v] = u
			for k := 1; k < K; k++ {
				up[k][v] = up[k-1][up[k-1][v]]
			}
			timer++
			tin[v] = timer
			stack = append(stack, v)
		} else {
			stack = stack[:len(stack)-1]
		}
	}

	lg := make([]int, n+2)
	for i := 2; i <= n; i++ {
		lg[i] = lg[i>>1] + 1
	}

	mn := make([][]int, K)
	mx := make([][]int, K)
	lc := make([][]int, K)
	for k := 0; k < K; k++ {
		mn[k] = make([]int, n+1)
		mx[k] = make([]int, n+1)
		lc[k] = make([]int, n+1)
	}

	for i := 1; i <= n; i++ {
		mn[0][i] = i
		mx[0][i] = i
		lc[0][i] = i
	}

	lcaFunc := func(a, b int) int {
		if depth[a] < depth[b] {
			a, b = b, a
		}
		diff := depth[a] - depth[b]
		bit := 0
		for diff > 0 {
			if diff&1 == 1 {
				a = up[bit][a]
			}
			diff >>= 1
			bit++
		}
		if a == b {
			return a
		}
		for k := K - 1; k >= 0; k-- {
			if up[k][a] != up[k][b] {
				a = up[k][a]
				b = up[k][b]
			}
		}
		return up[0][a]
	}

	for k := 1; k < K; k++ {
		half := 1 << (k - 1)
		span := 1 << k
		limit := n - span + 1
		prevMn := mn[k-1]
		prevMx := mx[k-1]
		prevLc := lc[k-1]
		curMn := mn[k]
		curMx := mx[k]
		curLc := lc[k]
		for i := 1; i <= limit; i++ {
			a := prevMn[i]
			b := prevMn[i+half]
			if tin[a] < tin[b] {
				curMn[i] = a
			} else {
				curMn[i] = b
			}
			a = prevMx[i]
			b = prevMx[i+half]
			if tin[a] > tin[b] {
				curMx[i] = a
			} else {
				curMx[i] = b
			}
			curLc[i] = lcaFunc(prevLc[i], prevLc[i+half])
		}
	}

	rangeMin := func(l, r int) int {
		k := lg[r-l+1]
		a := mn[k][l]
		b := mn[k][r-(1<<k)+1]
		if tin[a] < tin[b] {
			return a
		}
		return b
	}
	rangeMax := func(l, r int) int {
		k := lg[r-l+1]
		a := mx[k][l]
		b := mx[k][r-(1<<k)+1]
		if tin[a] > tin[b] {
			return a
		}
		return b
	}
	rangeLCA := func(l, r int) int {
		if l > r {
			return 0
		}
		k := lg[r-l+1]
		return lcaFunc(lc[k][l], lc[k][r-(1<<k)+1])
	}
	mergeLCA := func(a, b int) int {
		if a == 0 {
			return b
		}
		if b == 0 {
			return a
		}
		return lcaFunc(a, b)
	}

	out := make([]byte, 0, q*16)

	for ; q > 0; q-- {
		l := nextInt()
		r := nextInt()

		x := rangeMin(l, r)
		y := rangeMax(l, r)

		l1 := mergeLCA(rangeLCA(l, x-1), rangeLCA(x+1, r))
		d1 := depth[l1]

		l2 := mergeLCA(rangeLCA(l, y-1), rangeLCA(y+1, r))
		d2 := depth[l2]

		if d1 >= d2 {
			out = strconv.AppendInt(out, int64(x), 10)
			out = append(out, ' ')
			out = strconv.AppendInt(out, int64(d1), 10)
			out = append(out, '\n')
		} else {
			out = strconv.AppendInt(out, int64(y), 10)
			out = append(out, ' ')
			out = strconv.AppendInt(out, int64(d2), 10)
			out = append(out, '\n')
		}
	}

	return strings.TrimSpace(string(out))
}

// ─── verifier ───

func runProg(bin, input string) (string, error) {
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
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	q := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		if i > 2 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", p))
	}
	if n > 1 {
		sb.WriteByte('\n')
	} else {
		sb.WriteString("\n")
	}
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n) + 1
		if l > r {
			l, r = r, l
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []string{
		"1 1\n\n1 1\n",
		"2 1\n1\n1 2\n",
	}
	for i := 0; i < 100; i++ {
		tests = append(tests, genCase(rng))
	}

	for idx, input := range tests {
		exp := solveE(input)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if exp != got {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\n got: %s\ninput:\n%s", idx+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
