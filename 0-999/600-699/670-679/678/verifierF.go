package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded testcases (same format as original file).
const embeddedTestcases = `100
10
3 0
3 5
3 -5
3 -2
3 -5
1 -4 0
2 1
3 3
1 4 -2
1 -2 1

5
1 1 -3
1 -3 4
3 2
1 -3 -5
1 -2 -2

3
1 -1 0
1 3 5
3 -2

3
3 -2
3 -1
1 0 1

3
1 -1 -4
2 1
3 4

1
3 5

6
1 -1 0
2 1
3 0
1 2 2
3 -3
1 -1 -5

6
3 -5
3 1
3 1
3 -5
3 -5
3 -3

10
1 -4 -2
2 1
3 0
3 -1
3 -4
3 0
3 -5
3 -4
1 0 3
3 0

3
3 -1
3 3
1 -1 5

6
3 -3
1 5 -3
3 -1
2 1
3 -5
1 4 3

7
1 -2 4
2 1
3 5
3 -3
1 5 -5
2 1
1 -3 4

3
3 1
1 -3 1
2 1

1
3 -1

3
3 4
1 3 2
2 1

8
3 -1
3 1
1 -4 1
3 -3
3 2
2 1
1 2 -1
3 3

9
3 -4
3 4
3 -5
3 0
3 5
3 2
3 -1
3 5
1 4 -5

8
3 -1
3 5
3 2
3 3
3 5
3 0
3 5
3 1

6
1 2 0
2 1
3 -3
1 0 2
2 1
3 5

7
1 4 4
3 5
2 1
3 3
3 -1
3 -5
1 -3 4

8
3 5
1 -2 5
1 5 -5
2 1
1 -5 -3
1 0 -3
2 2
3 -5

7
3 0
3 5
3 -4
3 -2
1 0 -5
2 1
3 1

2
3 3
3 -5

9
3 -1
1 -1 3
3 0
3 -1
2 1
3 1
3 5
3 0
3 -3

3
3 1
3 2
1 -3 4

2
3 5
1 1 -4

6
3 4
3 -3
3 5
3 1
3 1
1 2 -1

8
3 1
3 -3
3 4
3 -1
3 -1
3 -5
3 -1
3 -1

3
3 -5
1 5 4
2 1

5
1 -3 1
1 2 3
3 -1
1 2 -5
1 2 -1

3
3 -1
3 2
3 2

9
3 4
3 -4
1 -3 -1
2 1
3 -1
3 3
1 2 0
2 1
1 -5 -1

9
3 5
1 5 3
1 -5 1
2 2
2 1
3 5
3 -2
3 2
1 -1 -5

7
3 -1
3 2
3 -3
1 2 3
2 1
3 -3
3 4

9
1 -4 -2
2 1
1 5 -5
2 1
1 1 2
1 -4 -4
1 4 5
1 -3 -1
3 2

3
1 4 -3
2 1
1 4 -2

1
3 -4

7
3 -4
3 -5
3 4
1 1 4
1 -5 1
1 0 5
3 2

8
3 5
3 -5
1 -1 -3
3 5
3 3
2 1
3 -4
3 -2

5
3 3
1 -2 0
2 1
1 -3 -5
3 0

2
3 5
1 -4 -4

9
3 2
1 1 2
2 1
1 3 -5
3 1
1 -3 1
3 -2
3 -4
1 5 2

4
1 4 1
2 1
3 0
3 0

7
3 -3
3 0
3 -1
3 1
3 3
1 4 4
1 -3 1

2
1 -5 -5
1 -2 -2

1
3 2

6
1 1 2
2 1
3 2
3 -4
1 -3 5
1 -4 0

7
3 -3
3 -1
1 -1 -3
2 1
1 2 -5
2 1
3 -5

1
3 2

8
1 -4 5
2 1
3 5
1 -3 -5
1 -5 4
1 5 -4
3 1
3 0

5
1 2 0
2 1
3 2
1 -2 -3
3 -2

1
3 4

7
3 -1
1 2 -5
2 1
1 -5 -5
3 3
3 0
1 -4 -2

8
1 2 -5
3 -2
3 -5
2 1
1 -5 -1
2 1
3 5
1 -5 -2

3
3 3
3 1
1 3 -2

2
3 -4
1 3 -3

10
1 -2 4
1 2 3
2 2
2 1
3 5
1 -5 -1
3 2
1 3 -4
1 -1 2
1 -1 1

6
1 -3 3
1 2 -2
3 -4
2 1
2 2
3 -2

1
3 -2

5
3 -4
1 1 0
2 1
3 5
1 4 -5

8
1 5 -5
3 -5
3 1
2 1
3 3
1 0 -4
2 1
1 -1 2

6
3 -5
1 3 3
1 -5 -4
1 -3 0
2 2
3 -5

4
3 1
1 0 -2
2 1
1 -3 4

10
1 -4 -1
3 4
1 5 -3
2 2
2 1
3 0
3 -2
3 4
3 -3
3 -5

3
3 2
3 4
1 -4 -3

7
3 -5
1 -5 -5
3 -2
2 1
3 2
1 -2 -1
1 -1 3

9
3 -2
3 0
1 5 -5
3 2
2 1
1 -3 -2
2 1
1 -4 -3
2 1

5
1 0 2
3 2
3 4
1 -2 1
3 -5

2
3 4
3 4

3
3 -5
3 5
1 4 -1

9
3 3
3 3
3 1
3 4
3 3
3 -5
3 -3
3 0
3 -5

9
1 -1 -1
1 -5 -4
2 2
2 1
3 -2
1 -5 -4
3 -1
2 1
3 3

2
1 -2 0
3 1

7
1 5 1
2 1
3 3
1 5 4
2 1
3 -4
3 3

6
1 -4 -2
2 1
3 4
1 -3 -4
3 1
1 -2 5

2
3 4
3 5

9
3 1
3 -2
3 2
3 4
1 -1 1
2 1
3 0
1 1 -5
1 -4 -1

10
3 4
3 -3
3 -2
3 4
1 5 0
3 -4
1 2 -3
3 -4
2 1
2 2

4
3 4
3 -4
1 -2 5
1 -4 -4

6
1 0 4
1 4 2
3 0
1 3 4
3 -4
2 1

2
1 -2 0
2 1

8
3 -4
1 1 3
2 1
1 -4 4
2 1
3 -5
3 1
3 -4

8
3 -1
1 -4 2
2 1
3 1
1 -1 2
3 -4
3 -4
2 1

6
1 3 1
1 -2 1
1 -1 5
3 -4
1 -1 -4
3 2

7
1 3 1
1 -4 2
2 1
1 -5 -3
2 2
1 -1 3
2 2

9
1 -5 -2
1 2 -5
2 1
2 2
3 -5
1 -3 1
3 1
3 5
3 4

10
1 -2 -5
2 1
3 1
1 -3 0
1 -1 -2
2 1
2 2
3 1
1 -1 -3
2 1

5
1 5 5
1 -3 2
1 -4 5
2 1
2 2

4
3 -5
3 2
3 4
3 -3

8
1 -4 0
3 -5
1 -3 1
1 -1 -3
3 0
2 1
1 -1 3
3 3

2
3 -4
1 5 1

6
3 0
3 3
3 2
1 4 -3
2 1
3 -5

7
1 -2 2
3 -5
1 -5 -2
2 1
3 4
1 -1 -5
3 -2

1
3 1

6
1 3 -1
2 1
1 5 -4
2 1
3 4
3 1

7
3 -3
3 5
1 5 2
1 -2 3
3 1
1 3 -1
3 0

7
1 2 0
1 2 -5
3 4
1 4 2
1 -5 -4
1 0 -3
2 3

5
3 2
1 -4 -1
2 1
3 -1
3 -3

4
3 -1
3 2
1 4 -3
3 -5

6
3 -4
1 2 -5
3 4
2 1
1 4 1
1 1 3`

// Embedded correct solver for 678F
const NEG int64 = -1 << 62

func solve678F(input string) string {
	data := []byte(input)
	p := 0
	nextInt64 := func() int64 {
		for p < len(data) && data[p] <= ' ' {
			p++
		}
		sign := int64(1)
		if p < len(data) && data[p] == '-' {
			sign = -1
			p++
		}
		var v int64
		for p < len(data) {
			c := data[p]
			if c < '0' || c > '9' {
				break
			}
			v = v*10 + int64(c-'0')
			p++
		}
		return sign * v
	}

	n := int(nextInt64())
	slope := make([]int64, n+1)
	intercept := make([]int64, n+1)
	removeAt := make([]int32, n+1)

	addIdx := make([]int, 0)
	queryTimes := make([]int, 0)
	queryVals := make([]int64, 0)
	xsAll := make([]int64, 0)

	for i := 1; i <= n; i++ {
		t := int(nextInt64())
		if t == 1 {
			a := nextInt64()
			b := nextInt64()
			slope[i] = a
			intercept[i] = b
			removeAt[i] = int32(n + 1)
			addIdx = append(addIdx, i)
		} else if t == 2 {
			idx := int(nextInt64())
			removeAt[idx] = int32(i)
		} else {
			q := nextInt64()
			queryTimes = append(queryTimes, i)
			queryVals = append(queryVals, q)
			xsAll = append(xsAll, q)
		}
	}

	qCount := len(queryTimes)
	if qCount == 0 {
		return ""
	}

	sort.Slice(xsAll, func(i, j int) bool { return xsAll[i] < xsAll[j] })
	k := 0
	for _, v := range xsAll {
		if k == 0 || xsAll[k-1] != v {
			xsAll[k] = v
			k++
		}
	}
	xs := xsAll[:k]
	m := len(xs)

	lowerBound64 := func(a []int64, x int64) int32 {
		l, r := 0, len(a)
		for l < r {
			mid := (l + r) >> 1
			if a[mid] < x {
				l = mid + 1
			} else {
				r = mid
			}
		}
		return int32(l)
	}
	lowerBoundInt := func(a []int, x int) int {
		l, r := 0, len(a)
		for l < r {
			mid := (l + r) >> 1
			if a[mid] < x {
				l = mid + 1
			} else {
				r = mid
			}
		}
		return l
	}
	upperBoundInt := func(a []int, x int) int {
		l, r := 0, len(a)
		for l < r {
			mid := (l + r) >> 1
			if a[mid] <= x {
				l = mid + 1
			} else {
				r = mid
			}
		}
		return l
	}

	queryIdxs := make([]int32, qCount)
	for i, v := range queryVals {
		queryIdxs[i] = lowerBound64(xs, v)
	}

	sizeT := 1
	for sizeT < qCount {
		sizeT <<= 1
	}

	type Edge struct {
		line int32
		next int32
	}
	type Change struct {
		node int32
		prev int32
	}

	countEdges := 0
	for _, idx := range addIdx {
		l := idx
		r := int(removeAt[idx]) - 1
		li := lowerBoundInt(queryTimes, l)
		ri := upperBoundInt(queryTimes, r) - 1
		if li > ri {
			continue
		}
		L := li + 1 + sizeT - 1
		R := ri + 1 + sizeT - 1
		for L <= R {
			if L&1 == 1 {
				countEdges++
				L++
			}
			if R&1 == 0 {
				countEdges++
				R--
			}
			L >>= 1
			R >>= 1
		}
	}

	headT := make([]int32, 2*sizeT)
	edgesT := make([]Edge, 1, countEdges+1)

	for _, idx := range addIdx {
		l := idx
		r := int(removeAt[idx]) - 1
		li := lowerBoundInt(queryTimes, l)
		ri := upperBoundInt(queryTimes, r) - 1
		if li > ri {
			continue
		}
		L := li + 1 + sizeT - 1
		R := ri + 1 + sizeT - 1
		for L <= R {
			if L&1 == 1 {
				edgesT = append(edgesT, Edge{int32(idx), headT[L]})
				headT[L] = int32(len(edgesT) - 1)
				L++
			}
			if R&1 == 0 {
				edgesT = append(edgesT, Edge{int32(idx), headT[R]})
				headT[R] = int32(len(edgesT) - 1)
				R--
			}
			L >>= 1
			R >>= 1
		}
	}

	lcTree := make([]int32, m*4+5)
	changes := make([]Change, 0, len(addIdx)*(bits.Len(uint(m))+1)+1)
	out := make([]byte, 0, qCount*24)

	var insertLine func(id int32)
	insertLine = func(id int32) {
		node, l, r := 1, 0, m-1
		mi, bi := slope[int(id)], intercept[int(id)]
		for {
			cur := lcTree[node]
			if cur == 0 {
				changes = append(changes, Change{int32(node), 0})
				lcTree[node] = id
				return
			}

			ci := int(cur)
			mc, bc := slope[ci], intercept[ci]
			mid := (l + r) >> 1
			xm := xs[mid]

			if mi*xm+bi > mc*xm+bc {
				changes = append(changes, Change{int32(node), cur})
				lcTree[node] = id
				id = cur
				mi, bi, mc, bc = mc, bc, mi, bi
			}

			if l == r {
				return
			}

			xl := xs[l]
			if mi*xl+bi > mc*xl+bc {
				node <<= 1
				r = mid
				continue
			}

			xr := xs[r]
			if mi*xr+bi > mc*xr+bc {
				node = node<<1 | 1
				l = mid + 1
				continue
			}

			return
		}
	}

	queryVal := func(idx int32) int64 {
		x := xs[int(idx)]
		res := NEG
		node, l, r := 1, 0, m-1
		for {
			cur := lcTree[node]
			if cur != 0 {
				ci := int(cur)
				v := slope[ci]*x + intercept[ci]
				if v > res {
					res = v
				}
			}
			if l == r {
				break
			}
			mid := (l + r) >> 1
			if int(idx) <= mid {
				node <<= 1
				r = mid
			} else {
				node = node<<1 | 1
				l = mid + 1
			}
		}
		return res
	}

	var dfs func(node, l, r int)
	dfs = func(node, l, r int) {
		if l > qCount {
			return
		}

		cp := len(changes)

		for e := headT[node]; e != 0; e = edgesT[e].next {
			insertLine(edgesT[e].line)
		}

		if l == r {
			res := queryVal(queryIdxs[l-1])
			if res == NEG {
				out = append(out, "EMPTY SET\n"...)
			} else {
				out = strconv.AppendInt(out, res, 10)
				out = append(out, '\n')
			}
		} else {
			mid := (l + r) >> 1
			dfs(node<<1, l, mid)
			dfs(node<<1|1, mid+1, r)
		}

		for len(changes) > cp {
			c := changes[len(changes)-1]
			lcTree[int(c.node)] = c.prev
			changes = changes[:len(changes)-1]
		}
	}

	dfs(1, 1, sizeT)

	return strings.TrimSpace(string(out))
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scan := bufio.NewScanner(strings.NewReader(embeddedTestcases))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t := 0
	fmt.Sscan(scan.Text(), &t)

	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		if !scan.Scan() {
			fmt.Printf("missing q for case %d\n", caseIdx)
			os.Exit(1)
		}
		q, _ := strconv.Atoi(scan.Text())
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", q))
		for i := 0; i < q; i++ {
			if !scan.Scan() {
				fmt.Printf("bad test file at case %d\n", caseIdx)
				os.Exit(1)
			}
			tok := scan.Text()
			tt, _ := strconv.Atoi(tok)
			sb.WriteString(tok)
			switch tt {
			case 1:
				var a, b int
				if !scan.Scan() {
					fmt.Printf("bad test file at case %d\n", caseIdx)
					os.Exit(1)
				}
				fmt.Sscan(scan.Text(), &a)
				if !scan.Scan() {
					fmt.Printf("bad test file at case %d\n", caseIdx)
					os.Exit(1)
				}
				fmt.Sscan(scan.Text(), &b)
				sb.WriteString(fmt.Sprintf(" %d %d", a, b))
			case 2:
				var idx int
				if !scan.Scan() {
					fmt.Printf("bad test file at case %d\n", caseIdx)
					os.Exit(1)
				}
				fmt.Sscan(scan.Text(), &idx)
				sb.WriteString(fmt.Sprintf(" %d", idx))
			default:
				var v int
				if !scan.Scan() {
					fmt.Printf("bad test file at case %d\n", caseIdx)
					os.Exit(1)
				}
				fmt.Sscan(scan.Text(), &v)
				sb.WriteString(fmt.Sprintf(" %d", v))
			}
			sb.WriteByte('\n')
		}
		input := sb.String()

		expect := solve678F(input)

		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", caseIdx, err, got)
			os.Exit(1)
		}

		// Compare using fields (whitespace-insensitive)
		gotFields := strings.Fields(got)
		expFields := strings.Fields(expect)
		match := len(gotFields) == len(expFields)
		if match {
			for k := range gotFields {
				if gotFields[k] != expFields[k] {
					match = false
					break
				}
			}
		}
		if !match {
			fmt.Printf("case %d failed:\nexpected:\n%s\ngot:\n%s\n", caseIdx, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
