package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesE.txt.
const embeddedTestcasesE = `100
2
3 0
7 3
2
1 1
1 2
5
3 2
5 0
8 1
9 2
12 1
2
3 3
3 1
3
4 1
6 1
10 2
1
2 3
3
5 2
10 1
14 3
5
2 2
1 2
1 1
3 2
3 2
2
2 0
6 1
4
1 2
2 1
1 1
1 2
2
3 3
4 0
3
1 2
1 2
2 1
4
5 0
8 1
12 2
14 2
4
2 3
1 3
3 3
1 4
4
2 0
3 0
4 1
9 1
5
1 4
1 3
3 4
2 4
4 4
2
4 3
6 3
2
1 2
2 2
2
2 3
5 1
3
1 2
1 2
1 2
5
2 2
5 3
8 3
13 1
16 2
4
4 1
2 3
5 2
3 1
6
4 3
8 3
10 0
12 1
13 2
14 3
4
2 5
2 2
3 5
4 5
1
1 0
4
1 1
1 1
1 1
1 1
3
4 1
8 3
9 0
1
3 2
4
4 3
5 1
8 3
12 0
5
2 3
2 2
3 3
1 2
3 1
4
5 3
10 0
15 1
19 2
3
2 2
1 1
3 4
5
3 3
8 3
12 1
15 2
18 1
4
1 5
2 5
3 2
3 5
1
4 1
3
1 1
1 1
1 1
1
2 1
5
1 1
1 1
1 1
1 1
1 1
6
3 3
7 3
10 2
14 3
19 0
22 1
5
4 6
3 4
1 5
1 4
6 3
3
1 2
5 1
9 0
2
1 1
1 1
5
5 2
6 1
8 0
10 0
13 0
1
3 1
6
3 1
5 3
7 0
12 2
13 1
18 3
2
2 3
1 4
3
5 3
10 2
14 2
2
1 2
1 2
1
3 1
1
1 1
1
2 1
5
1 1
1 1
1 1
1 1
1 1
4
5 0
7 2
12 2
16 2
1
2 3
6
1 2
6 3
9 0
11 0
15 1
16 1
1
3 5
6
3 0
8 0
9 1
14 2
19 0
20 0
4
2 6
2 6
3 4
1 4
4
1 2
6 2
8 3
11 1
2
1 4
3 3
4
4 2
9 1
12 2
17 0
1
4 3
6
2 2
3 1
8 0
10 2
12 3
13 2
1
4 5
2
1 1
6 0
1
2 2
5
2 0
4 0
6 1
9 0
14 2
3
5 3
2 2
1 1
6
1 2
5 0
7 1
11 3
16 3
21 2
3
4 6
1 1
2 3
3
5 1
10 2
14 2
5
3 1
1 2
1 3
1 1
3 3
5
4 0
7 3
9 1
13 3
16 2
1
2 4
6
4 1
5 2
8 0
9 3
11 0
13 1
5
3 5
1 6
2 5
1 2
2 3
3
5 2
7 2
12 2
3
3 2
1 3
1 3
1
1 1
2
1 1
1 1
1
1 1
5
1 1
1 1
1 1
1 1
1 1
3
5 2
6 1
10 0
4
1 2
3 1
1 3
1 1
1
3 3
1
1 1
6
2 0
3 0
5 3
7 0
12 2
17 2
3
2 5
5 6
1 5
3
5 3
8 2
10 2
5
2 1
2 1
3 3
3 3
1 2
6
5 3
8 1
11 1
14 2
17 1
22 3
5
3 3
2 6
2 3
2 2
3 1
1
1 1
2
1 1
1 1
4
4 2
9 0
12 0
13 1
1
1 4
2
4 0
9 1
3
2 2
1 1
1 2
4
3 3
8 0
12 1
13 3
4
1 2
1 2
4 3
3 4
3
4 0
7 3
12 0
4
2 3
1 3
1 2
2 3
6
2 1
6 1
11 2
15 3
19 0
23 0
4
3 4
3 6
3 3
5 6
6
4 1
8 3
12 2
17 0
21 2
26 1
5
1 2
5 6
1 1
2 3
3 4
3
1 0
2 3
3 1
2
2 3
1 3
1
2 1
2
1 1
1 1
6
3 0
8 2
11 0
12 2
13 3
14 1
2
1 4
3 1
6
1 1
5 3
10 3
11 1
13 0
15 3
5
5 4
3 5
3 6
1 3
5 3
3
1 1
6 1
8 3
5
1 3
1 3
3 3
1 2
2 3
5
1 2
4 1
6 0
9 3
12 2
1
1 4
2
4 2
6 3
2
1 1
1 1
5
4 2
6 1
9 3
13 2
16 0
4
3 4
1 3
5 5
1 4
6
5 0
10 2
15 0
17 3
22 1
26 3
2
1 3
2 4
2
4 3
7 3
3
2 2
1 1
1 1
3
3 0
5 2
6 3
4
1 2
2 2
2 3
2 2
2
4 2
7 0
2
1 2
1 1
6
5 0
8 1
9 3
13 2
17 1
18 0
2
4 3
1 6
5
4 1
5 2
8 2
11 3
14 0
5
2 5
3 5
1 5
3 5
1 4
4
5 3
10 1
14 1
17 2
3
1 4
2 3
4 3
5
3 0
6 3
9 2
11 1
13 2
5
1 3
4 4
4 5
4 2
3 3
3
4 1
6 3
7 3
5
1 1
1 3
2 2
1 1
1 2
6
3 2
4 2
9 0
14 3
18 3
20 1
3
4 2
3 4
6 6
2
2 1
6 1
4
1 1
1 2
1 2
1 1
3
2 3
3 2
5 1
5
1 3
1 2
1 3
1 1
1 1
3
3 2
7 3
9 0
5
2 2
3 3
3 1
2 3
1 3
6
3 0
8 3
9 3
13 3
18 3
23 3
3
5 1
6 6
4 5
5
3 3
4 2
6 3
8 3
11 0
2
4 3
5 5
2
1 0
5 1
4
1 1
1 1
2 1
1 2
5
1 3
5 3
8 3
10 0
14 3
4
5 4
3 3
2 3
4 2
2
4 2
8 1
3
2 1
1 2
1 1
5
1 2
5 3
10 0
15 2
17 1
1
2 4
2
1 3
4 1
1
1 1
6
5 3
8 0
11 3
13 3
16 1
20 1
3
1 1
3 4
3 6
3
1 2
6 0
9 3
4
2 3
1 3
3 3
2 3
1
4 3
5
1 1
1 1
1 1
1 1
1 1
2
3 1
8 1
3
1 1
2 2
1 1
1
4 0
2
1 1
1 1
2
5 3
6 3
5
2 2
1 2
2 1
2 2
1 2
1
3 1
5
1 1
1 1
1 1
1 1
1 1
2
1 3
3 3
3
1 2
1 2
2 2
6
5 3
7 1
11 1
14 2
17 3
22 3
4
6 6
2 6
2 4
5 6
3
1 3
2 2
7 3
2
3 3
2 3
2
1 3
6 1
2
1 2
2 1
4
5 2
9 3
11 0
12 1
2
2 3
1 3
2
5 0
10 1
5
1 1
1 2
1 1
2 1
2 2
1
4 3
3
1 1
1 1
1 1
5
2 0
3 0
5 3
8 2
10 1
4
1 4
2 3
1 2
4 1
1
5 1
3
1 1
1 1
1 1
3
4 3
5 3
9 0
3
1 1
2 3
3 2
4
4 1
8 1
10 2
11 2
1
4 2
5
2 0
5 2
8 1
11 2
14 2
5
3 1
1 5
2 4
2 5
2 4`

// ---------- correct reference solver ----------

func nextIntBuf(data []byte, idx *int) int64 {
	n := len(data)
	for *idx < n {
		c := data[*idx]
		if c > ' ' {
			break
		}
		*idx++
	}
	var v int64
	for *idx < n {
		c := data[*idx]
		if c < '0' || c > '9' {
			break
		}
		v = v*10 + int64(c-'0')
		*idx++
	}
	return v
}

func upperBound(p []int64, n int, x int64) int {
	l, r := 1, n+1
	for l < r {
		m := (l + r) >> 1
		if p[m] > x {
			r = m
		} else {
			l = m + 1
		}
	}
	return l - 1
}

func updateI32(seg []int32, size, pos int, val int32) {
	i := size + pos - 1
	seg[i] = val
	for i >>= 1; i > 0; i >>= 1 {
		if seg[i<<1] > seg[i<<1|1] {
			seg[i] = seg[i<<1]
		} else {
			seg[i] = seg[i<<1|1]
		}
	}
}

func queryMaxI32(seg []int32, size, l, r int) int32 {
	var res int32
	l += size - 1
	r += size - 1
	for l <= r {
		if l&1 == 1 {
			if seg[l] > res {
				res = seg[l]
			}
			l++
		}
		if r&1 == 0 {
			if seg[r] > res {
				res = seg[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func queryMaxI64(seg []int64, size, l, r int) int64 {
	var res int64
	l += size - 1
	r += size - 1
	for l <= r {
		if l&1 == 1 {
			if seg[l] > res {
				res = seg[l]
			}
			l++
		}
		if r&1 == 0 {
			if seg[r] > res {
				res = seg[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func solve500E(input string) string {
	data := []byte(input)
	idx := 0

	n := int(nextIntBuf(data, &idx))
	p := make([]int64, n+1)
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		pi := nextIntBuf(data, &idx)
		li := nextIntBuf(data, &idx)
		p[i] = pi
		a[i] = pi + li
	}

	size := 1
	for size < n {
		size <<= 1
	}

	segR := make([]int32, size<<1)
	R := make([]int32, n+2)

	for i := n; i >= 1; i-- {
		ri := upperBound(p, n, a[i])
		val := int32(ri)
		if i+1 <= ri {
			mx := queryMaxI32(segR, size, i+1, ri)
			if mx > val {
				val = mx
			}
		}
		R[i] = val
		updateI32(segR, size, i, val)
	}

	segA := make([]int64, size<<1)
	for i := 1; i <= n; i++ {
		segA[size+i-1] = a[i]
	}
	for i := size - 1; i >= 1; i-- {
		if segA[i<<1] > segA[i<<1|1] {
			segA[i] = segA[i<<1]
		} else {
			segA[i] = segA[i<<1|1]
		}
	}

	K := 1
	for (1 << K) <= n+1 {
		K++
	}

	up := make([][]int32, K)
	mx := make([][]int32, K)
	sum := make([][]int64, K)
	for k := 0; k < K; k++ {
		up[k] = make([]int32, n+2)
		mx[k] = make([]int32, n+2)
		sum[k] = make([]int64, n+2)
	}

	sentinel := n + 1
	up[0][sentinel] = int32(sentinel)
	mx[0][sentinel] = int32(sentinel)

	for i := 1; i <= n; i++ {
		up[0][i] = R[i] + 1
		mx[0][i] = R[i]
		if R[i] < int32(n) {
			best := queryMaxI64(segA, size, i, int(R[i]))
			sum[0][i] = p[int(R[i])+1] - best
		}
	}

	for k := 1; k < K; k++ {
		upPrev := up[k-1]
		mxPrev := mx[k-1]
		sumPrev := sum[k-1]
		upCur := up[k]
		mxCur := mx[k]
		sumCur := sum[k]
		for i := 1; i <= n+1; i++ {
			mid := int(upPrev[i])
			upCur[i] = upPrev[mid]
			mxCur[i] = mxPrev[mid]
			sumCur[i] = sumPrev[i] + sumPrev[mid]
		}
	}

	q := int(nextIntBuf(data, &idx))
	out := make([]byte, 0, q*16)

	for ; q > 0; q-- {
		x := int(nextIntBuf(data, &idx))
		y := int32(nextIntBuf(data, &idx))
		cur := x
		var ans int64
		for k := K - 1; k >= 0; k-- {
			if mx[k][cur] < y {
				ans += sum[k][cur]
				cur = int(up[k][cur])
			}
		}
		out = strconv.AppendInt(out, ans, 10)
		out = append(out, '\n')
	}

	return strings.TrimSpace(string(out))
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(embeddedTestcasesE))
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, 1024), 1<<20)
	nextInt := func() int {
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "unexpected EOF in test data")
			os.Exit(1)
		}
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
			os.Exit(1)
		}
		return v
	}

	T := nextInt()
	for caseIdx := 1; caseIdx <= T; caseIdx++ {
		n := nextInt()
		coords := make([][2]int64, n)
		for i := 0; i < n; i++ {
			x := int64(nextInt())
			l := int64(nextInt())
			coords[i] = [2]int64{x, l}
		}
		q := nextInt()
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			x := nextInt()
			y := nextInt()
			queries[i] = [2]int{x, y}
		}

		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		for _, c := range coords {
			fmt.Fprintf(&input, "%d %d\n", c[0], c[1])
		}
		fmt.Fprintf(&input, "%d\n", q)
		for _, qu := range queries {
			fmt.Fprintf(&input, "%d %d\n", qu[0], qu[1])
		}

		wantRaw := solve500E(input.String())
		wantFields := strings.Fields(wantRaw)

		gotRaw, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseIdx, err)
			os.Exit(1)
		}
		gotFields := strings.Fields(gotRaw)
		if len(gotFields) != len(wantFields) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d outputs, got %d\n", caseIdx, len(wantFields), len(gotFields))
			os.Exit(1)
		}
		for i := range gotFields {
			if gotFields[i] != wantFields[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at query %d: expected %s got %s\n", caseIdx, i+1, wantFields[i], gotFields[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d testcases passed\n", T)
}
