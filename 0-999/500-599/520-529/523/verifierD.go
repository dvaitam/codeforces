package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
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

type minHeap []int64

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func expectedOutput(n, k int, sVals, mVals []int64) string {
	res := make([]int64, n)
	h := &minHeap{}
	heap.Init(h)
	free := k
	type job struct {
		idx int
		m   int64
	}
	queue := make([]job, 0)
	process := func(t int64) {
		for h.Len() > 0 {
			ft := (*h)[0]
			if ft > t {
				break
			}
			heap.Pop(h)
			if len(queue) > 0 {
				j := queue[0]
				queue = queue[1:]
				start := ft
				finish := start + j.m
				res[j.idx] = finish
				heap.Push(h, finish)
			} else {
				free++
			}
		}
	}
	for i := 0; i < n; i++ {
		ti := sVals[i]
		process(ti)
		if free > 0 {
			free--
			start := ti
			finish := start + mVals[i]
			res[i] = finish
			heap.Push(h, finish)
		} else {
			queue = append(queue, job{i, mVals[i]})
		}
	}
	for len(queue) > 0 {
		ft := heap.Pop(h).(int64)
		j := queue[0]
		queue = queue[1:]
		start := ft
		finish := start + j.m
		res[j.idx] = finish
		heap.Push(h, finish)
	}
	var b strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "%d", res[i])
		if i+1 < n {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

const testcasesDRaw = `100
4 2
6 10
14 10
16 10
17 8
5 5
4 4
12 9
21 8
28 3
32 3
7 6
1 2
4 10
5 5
6 5
14 10
21 7
28 10
8 3
6 2
7 3
15 4
20 7
25 7
34 7
44 6
53 10
7 5
4 6
5 5
15 3
21 9
31 10
33 4
43 5
5 1
2 8
10 2
16 2
23 3
24 5
7 7
7 2
8 10
18 1
25 10
31 9
36 9
40 1
5 1
2 2
12 9
13 4
20 5
30 5
3 3
1 6
7 6
10 7
7 4
9 7
19 9
21 10
30 5
37 4
42 7
47 9
5 5
6 1
13 10
19 1
26 10
36 3
1 1
8 6
6 5
5 8
6 10
7 1
13 5
21 5
31 10
6 2
6 3
12 6
22 5
27 7
29 1
39 3
5 5
4 5
8 6
11 7
13 2
23 6
6 6
4 8
7 2
13 4
23 8
28 4
30 1
4 3
10 3
15 6
17 10
23 10
3 2
5 9
10 8
16 7
5 4
10 7
11 7
14 4
15 8
25 9
7 5
4 1
12 9
17 9
23 4
25 10
30 2
34 1
1 1
7 10
1 1
8 2
3 3
5 4
6 9
15 7
1 1
6 3
5 5
8 1
14 4
18 2
27 2
30 4
5 2
1 8
11 7
12 5
16 5
26 9
7 1
8 6
9 1
12 1
14 1
16 8
17 2
26 9
8 6
3 6
5 6
12 7
22 5
28 5
32 6
39 2
42 9
1 1
2 10
3 1
6 8
16 9
23 1
7 1
6 8
12 7
19 8
20 4
24 9
29 10
31 7
4 4
3 1
9 6
18 5
20 8
2 2
2 6
12 9
2 1
8 3
12 7
1 1
10 2
7 2
1 6
3 1
5 8
10 10
15 2
16 10
25 9
4 1
9 2
18 1
27 6
37 3
2 1
3 4
11 10
7 3
6 10
13 6
22 7
24 7
33 4
40 3
47 10
8 3
7 3
10 2
18 8
27 8
37 3
40 5
44 3
54 9
6 2
9 5
16 10
26 10
31 4
36 1
41 8
7 2
3 10
9 4
15 8
18 7
26 10
30 8
40 9
1 1
2 7
1 1
4 4
2 1
5 4
9 5
3 1
10 1
15 3
16 6
3 2
2 2
4 2
9 5
1 1
8 10
6 1
1 6
7 7
14 8
16 4
26 8
33 3
6 1
5 2
12 2
20 9
25 2
34 6
40 8
5 3
2 6
12 9
21 2
29 9
35 1
5 5
3 3
6 6
14 2
16 9
19 6
7 5
5 3
13 8
18 3
20 2
23 9
32 10
39 6
2 2
5 7
6 3
1 1
9 5
4 3
6 7
14 9
16 6
24 2
3 2
10 2
12 10
14 3
4 4
7 3
17 10
20 7
24 9
3 3
3 4
8 6
13 1
8 7
7 6
16 10
21 8
30 5
38 1
48 4
49 2
53 8
3 3
8 4
12 9
16 1
8 2
10 5
13 3
21 2
31 1
32 6
42 4
51 2
59 9
1 1
6 6
6 6
3 2
13 1
15 6
19 2
23 7
27 8
6 1
1 7
3 4
6 7
14 8
16 9
23 4
8 5
1 8
9 7
17 3
25 1
30 6
36 8
45 6
55 7
4 1
4 5
10 3
18 9
22 3
4 1
3 10
10 9
13 1
16 2
3 2
8 3
9 1
16 8
6 4
1 1
5 7
6 7
14 1
18 4
20 7
8 4
3 6
13 2
19 2
29 1
34 5
42 5
50 4
59 5
1 1
6 6
2 1
7 2
17 10
1 1
1 2
1 1
9 1
8 1
4 9
10 4
18 6
26 6
27 7
32 10
39 2
44 3
7 1
9 7
18 6
27 7
30 7
39 6
42 6
49 8
4 4
8 6
13 3
22 10
29 8
1 1
3 1
8 2
2 6
6 10
7 10
8 8
16 6
22 1
24 4
31 2
6 5
5 2
13 2
17 4
18 3
21 10
22 2
4 3
4 4
13 9
20 9
30 6
4 4
3 10
5 1
7 10
8 2
4 3
2 2
10 7
14 10
16 8
6 4
10 8
12 5
22 8
29 4
31 9
32 8
5 1
6 6
10 8
12 9
18 7
20 10
4 2
6 1
12 4
19 8
21 5
4 3
3 4
7 10
15 9
22 6
4 4
8 7
16 10
17 5
18 3
2 1
3 5
12 9
1 1
1 4
4 3
8 7
9 6
17 4
22 3
2 2
5 7
13 2
4 2
8 5
15 6
18 7
23 8
8 4
6 5
11 1
19 6
25 5
29 9
30 1
33 9
36 9
1 1
1 1
4 4
6 6
15 1
23 3
27 1
5 4
6 1
16 9
18 8
23 5
27 8
7 6
5 6
6 1
13 1
16 10
20 3
27 9
33 9
3 2
1 3
2 1
10 1
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	reader := bufio.NewReader(strings.NewReader(testcasesDRaw))
	var tcases int
	if _, err := fmt.Fscan(reader, &tcases); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read test count: %v\n", err)
		os.Exit(1)
	}
	for caseNum := 1; caseNum <= tcases; caseNum++ {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		sVals := make([]int64, n)
		mVals := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &sVals[i], &mVals[i])
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", n, k)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&input, "%d %d\n", sVals[i], mVals[i])
		}
		want := expectedOutput(n, k, sVals, mVals)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", caseNum, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tcases)
}
