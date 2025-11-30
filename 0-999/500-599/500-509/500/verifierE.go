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

type block struct {
	start int
	reach int64
}

func solve500E(n int, coords [][2]int64, queries [][2]int) []int64 {
	p := make([]int64, n)
	r := make([]int64, n)
	for i := 0; i < n; i++ {
		p[i] = coords[i][0]
		r[i] = coords[i][0] + coords[i][1]
	}

	blocks := make([]block, 0, n)
	for i := 0; i < n; i++ {
		start := i
		reach := r[i]
		for len(blocks) > 0 && blocks[len(blocks)-1].reach >= p[i] {
			prev := blocks[len(blocks)-1]
			blocks = blocks[:len(blocks)-1]
			if prev.start < start {
				start = prev.start
			}
			if prev.reach > reach {
				reach = prev.reach
			}
		}
		blocks = append(blocks, block{start: start, reach: reach})
	}

	bcnt := len(blocks)
	blk := make([]int, n)
	for bi := 0; bi < bcnt; bi++ {
		s := blocks[bi].start
		var e int
		if bi+1 < bcnt {
			e = blocks[bi+1].start - 1
		} else {
			e = n - 1
		}
		for i := s; i <= e; i++ {
			blk[i] = bi
		}
	}

	gaps := make([]int64, bcnt-1)
	for bi := 0; bi+1 < bcnt; bi++ {
		nextStart := blocks[bi+1].start
		gaps[bi] = p[nextStart] - blocks[bi].reach
	}
	prefix := make([]int64, bcnt)
	for i := 1; i < bcnt; i++ {
		prefix[i] = prefix[i-1] + gaps[i-1]
	}

	ans := make([]int64, len(queries))
	for i, q := range queries {
		x := q[0] - 1
		y := q[1] - 1
		bx := blk[x]
		by := blk[y]
		if bx >= by {
			ans[i] = 0
		} else {
			ans[i] = prefix[by] - prefix[bx]
		}
	}
	return ans
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

		wantVals := solve500E(n, coords, queries)
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		for _, c := range coords {
			fmt.Fprintf(&input, "%d %d\n", c[0], c[1])
		}
		fmt.Fprintf(&input, "%d\n", q)
		for _, qu := range queries {
			fmt.Fprintf(&input, "%d %d\n", qu[0], qu[1])
		}

		gotRaw, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseIdx, err)
			os.Exit(1)
		}
		gotFields := strings.Fields(gotRaw)
		if len(gotFields) != len(wantVals) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d outputs, got %d\n", caseIdx, len(wantVals), len(gotFields))
			os.Exit(1)
		}
		for i, g := range gotFields {
			val, err := strconv.ParseInt(g, 10, 64)
			if err != nil || val != wantVals[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at query %d: expected %d got %s\n", caseIdx, i+1, wantVals[i], g)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d testcases passed\n", T)
}
