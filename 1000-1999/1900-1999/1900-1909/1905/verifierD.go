package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcases = `4
0 2 3 1

2
1 0

8
0 6 2 3 5 1 7 4

7
0 3 2 4 6 1 5

1
0

2
1 0

1
0

5
4 3 1 2 0

7
2 1 0 6 4 5 3

2
1 0

3
2 0 1

5
0 4 1 2 3

6
1 0 2 3 5 4

6
1 4 3 2 0 5

6
1 3 2 0 5 4

5
4 1 3 0 2

8
6 4 5 7 3 0 2 1

7
1 5 2 4 0 3 6

7
1 0 3 6 2 4 5

1
0

5
1 2 3 4 0

4
0 2 1 3

3
1 0 2

6
0 4 3 1 5 2

8
3 1 2 0 7 4 5 6

7
0 3 4 2 1 6 5

5
3 1 0 2 4

6
2 4 5 1 3 0

8
0 1 3 6 4 2 7 5

1
0

1
0

1
0

6
0 1 4 5 3 2

3
2 0 1

6
0 3 1 5 4 2

2
1 0

3
2 0 1

5
3 4 0 2 1

2
1 0

6
4 5 0 3 1 2

6
2 0 4 3 1 5

2
1 0

4
1 0 3 2

5
3 1 4 0 2

3
0 2 1

5
0 1 4 2 3

5
2 1 0 4 3

3
1 2 0

8
3 7 2 0 1 5 4 6

5
1 3 0 2 4

5
3 2 4 1 0

4
1 0 2 3

1
0

8
3 0 5 6 2 4 7 1

7
3 1 5 2 6 4 0

8
4 3 6 7 5 1 2 0

3
2 1 0

3
2 1 0

7
6 4 3 5 1 2 0

1
0

8
3 2 4 1 7 0 6 5

1
0

2
0 1

1
0

2
0 1

6
3 4 5 0 2 1

7
3 1 0 6 5 2 4

7
3 5 2 6 4 1 0

2
1 0

1
0

6
1 0 2 5 4 3

7
6 4 1 5 3 2 0

8
5 7 3 2 4 6 1 0

7
5 2 4 0 6 3 1

5
1 4 2 3 0

7
1 3 4 6 2 0 5

8
6 5 7 4 0 3 1 2

7
6 4 3 5 2 0 1

8
1 3 5 6 0 2 7 4

2
1 0

1
0

6
5 2 3 0 1 4

4
0 1 2 3

5
0 4 1 3 2

2
0 1

4
1 0 2 3

7
2 1 0 3 6 4 5

3
1 2 0

8
2 6 0 1 3 4 5 7

4
0 3 2 1

4
0 1 3 2

5
3 4 0 2 1

8
3 5 0 2 4 7 1 6

8
7 0 4 1 6 5 3 2

8
7 2 4 5 0 3 6 1

2
1 0

5
2 0 3 4 1

3
1 0 2

5
2 3 4 0 1

7
2 1 3 4 6 5 0`

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

type segtree struct {
	n    int
	sum  []int64
	lazy []int64
	set  []bool
}

func newSegTree(arr []int64) *segtree {
	n := 1
	for n < len(arr) {
		n <<= 1
	}
	st := &segtree{
		n:    n,
		sum:  make([]int64, 2*n),
		lazy: make([]int64, 2*n),
		set:  make([]bool, 2*n),
	}
	for i, v := range arr {
		st.sum[n+i] = v
	}
	for i := n - 1; i > 0; i-- {
		st.sum[i] = st.sum[i<<1] + st.sum[i<<1|1]
	}
	return st
}

func (st *segtree) apply(idx int, val int64, length int) {
	st.sum[idx] = val * int64(length)
	st.lazy[idx] = val
	st.set[idx] = true
}

func (st *segtree) push(idx, length int) {
	if st.set[idx] {
		st.apply(idx<<1, st.lazy[idx], length>>1)
		st.apply(idx<<1|1, st.lazy[idx], length>>1)
		st.set[idx] = false
	}
}

func (st *segtree) update(idx, l, r, ql, qr int, val int64) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(idx, val, r-l+1)
		return
	}
	st.push(idx, r-l+1)
	mid := (l + r) >> 1
	st.update(idx<<1, l, mid, ql, qr, val)
	st.update(idx<<1|1, mid+1, r, ql, qr, val)
	st.sum[idx] = st.sum[idx<<1] + st.sum[idx<<1|1]
}

func (st *segtree) Update(l, r int, val int64) {
	if l > r {
		return
	}
	st.update(1, 0, st.n-1, l, r, val)
}

func (st *segtree) QueryAll() int64 {
	return st.sum[1]
}

func solve(perm []int) int64 {
	n := len(perm)
	pos := make([]int, n)
	for i, v := range perm {
		pos[v] = i
	}
	maxPos := make([]int, n+1)
	mx := -1
	for k := 1; k <= n; k++ {
		if pos[k-1] > mx {
			mx = pos[k-1]
		}
		maxPos[k] = mx
	}
	arr := make([]int64, n)
	for k := 1; k <= n; k++ {
		arr[k-1] = int64(maxPos[k] - n)
	}
	st := newSegTree(arr)
	sumLast := st.QueryAll()
	best := int64(0)
	for s := 0; s < n; s++ {
		cost := int64(s)*int64(n) - sumLast
		if cost > best {
			best = cost
		}
		v := perm[s]
		if v+1 <= n {
			st.Update(v, n-1, int64(s))
			sumLast = st.QueryAll()
		}
	}
	return best
}

type testCase struct {
	n    int
	perm []int
}

func parseCases(raw string) ([]testCase, error) {
	blocks := strings.Split(raw, "\n\n")
	var res []testCase
	for _, b := range blocks {
		b = strings.TrimSpace(b)
		if b == "" {
			continue
		}
		lines := strings.Split(b, "\n")
		if len(lines) < 2 {
			return nil, fmt.Errorf("invalid block: %q", b)
		}
		n, err := strconv.Atoi(strings.TrimSpace(lines[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid n: %w", err)
		}
		fields := strings.Fields(lines[1])
		if len(fields) != n {
			return nil, fmt.Errorf("perm length %d expected %d", len(fields), n)
		}
		perm := make([]int, n)
		for i, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("invalid perm value: %w", err)
			}
			perm[i] = v
		}
		res = append(res, testCase{n: n, perm: perm})
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases(testcases)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		want := solve(tc.perm)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.perm {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		gotStr, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(gotStr), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d output parse error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
