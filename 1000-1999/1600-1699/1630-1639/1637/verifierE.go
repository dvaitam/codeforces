package main

import (
	"bytes"
	"compress/gzip"
	"container/heap"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded gzipped+base64 testcases from testcasesE.txt.
const encodedTestcases = `
H4sIAOsFK2kC/0VU17HkMAz7VxUqwaRy/40dAvfd7IzXlhgBkKfdfnt+PXu8vnu03aLH6pF99Pj6a/E1fODq9tNj4/FoC7/2Gl6Dh9l5Aju4xGnXLrCfMlgd/7Nlg+Nqo/EYaXGA1OE4jzZbMYLpTj/MzTomIk2eDGcOxtBBMJWC4ZlsBP/4TiZ5DoI3Hg1kQ9jXFosLNQS3wP1WkaNiMeOqW/YEP2R2rcc2pyAbjENvtkawZmOrC9gAx8H8LIlAJO6CwRcQmI3VIinaaMbtsauL89GJxKMPoUW5i31fg9MASghWRya+MLkGjZQhiBoXVaG3j83BUQzSb7C3QWSSufiiSn9lk9AwHZcdK8xRv/gpi7pRDSxbnaWqrBRwBGbodktB26b2xm85PlyzFMIqXe5jbY8+IR5kK9nBt+gTaGi4Cv4qvyDafxRdOL3SG3HHzWIfBEIyWiL1P2VLZaaOHoqTrYl0Prb+DDcJtw4ky01IpppLIiLcTzHhDi9Nnr5hSM0gJv5AMDyGTv2UytFydbtMpeJeptNAhbogBFlOcNmSKj1gvIS0gFgeEspA5YgNsZLWwXZrvB2lbOfyrH6i0UMNKLsQ9a1w1Rr41QVTTVUa/VU6QUmWdXLpVD46Zbemj1eM50mRZ2PGNHdHhE4WoakmpuEYAklcBu/zh0IaHRb+iAEFR/3HqYRDsmFzXiNSGcs1SZbJ5khpIXCUkJ1DGadkO4V0Sr/eQCptE+up1SgRY4+R/yg9DG9JJMx+f4CfGsutTGTfG9akow5uB8N0pLihbYc2FxdQrdvrPRXe6V4unp4tOLwfriFzGWJAe29pFf/tdrlFr0VmrlOUDH+zhnilLm9QtqR51aAODQypolqpM3IB4v4B9R0LYH0GAAA=
`

type node struct {
	sum int64
	i   int
	j   int
}

type maxHeap []node

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].sum > h[j].sum }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) {
	*h = append(*h, x.(node))
}
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type testCase struct {
	n     int
	m     int
	arr   []int64
	pairs [][2]int64
}

func pairKey(a, b int64) int64 {
	if a > b {
		a, b = b, a
	}
	return (a << 32) | b
}

// solve mirrors 1637E.go.
func solve(tc testCase) string {
	freq := make(map[int64]int)
	for _, x := range tc.arr {
		freq[x]++
	}
	groups := make(map[int][]int64)
	for v, c := range freq {
		groups[c] = append(groups[c], v)
	}
	counts := make([]int, 0, len(groups))
	for c := range groups {
		counts = append(counts, c)
		sort.Slice(groups[c], func(i, j int) bool { return groups[c][i] > groups[c][j] })
	}
	sort.Ints(counts)

	bad := make(map[int64]struct{})
	for _, p := range tc.pairs {
		bad[pairKey(p[0], p[1])] = struct{}{}
	}

	var ans int64
	for i := 0; i < len(counts); i++ {
		c1 := counts[i]
		arr1 := groups[c1]
		for j := i; j < len(counts); j++ {
			c2 := counts[j]
			arr2 := groups[c2]
			if c1 == c2 && len(arr1) < 2 {
				continue
			}
			h := &maxHeap{}
			visited := make(map[int64]struct{})
			if c1 == c2 {
				heap.Push(h, node{int64(arr1[0] + arr2[1]), 0, 1})
				visited[int64(0)<<32|1] = struct{}{}
			} else {
				heap.Push(h, node{int64(arr1[0] + arr2[0]), 0, 0})
				visited[int64(0)<<32|0] = struct{}{}
			}
			found := false
			for h.Len() > 0 && !found {
				cur := heap.Pop(h).(node)
				x := arr1[cur.i]
				y := arr2[cur.j]
				if x != y {
					if _, ok := bad[pairKey(x, y)]; !ok {
						val := int64(c1+c2) * (x + y)
						if val > ans {
							ans = val
						}
						found = true
						break
					}
				}
				ni, nj := cur.i+1, cur.j
				if ni < len(arr1) {
					key := int64(ni)<<32 | int64(nj)
					if _, ok := visited[key]; !ok {
						visited[key] = struct{}{}
						heap.Push(h, node{int64(arr1[ni] + arr2[nj]), ni, nj})
					}
				}
				ni, nj = cur.i, cur.j+1
				if nj < len(arr2) {
					key := int64(ni)<<32 | int64(nj)
					if _, ok := visited[key]; !ok {
						visited[key] = struct{}{}
						heap.Push(h, node{int64(arr1[ni] + arr2[nj]), ni, nj})
					}
				}
			}
		}
	}
	return fmt.Sprint(ans)
}

func decodeTestcases() ([]testCase, error) {
	data, err := base64.StdEncoding.DecodeString(encodedTestcases)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var out bytes.Buffer
	if _, err := out.ReadFrom(r); err != nil {
		return nil, err
	}
	fields := strings.Fields(out.String())
	pos := 0
	var cases []testCase
	for pos+1 < len(fields) {
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, err
		}
		m, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return nil, err
		}
		pos += 2
		if pos+n > len(fields) {
			break
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			val, err := strconv.ParseInt(fields[pos+i], 10, 64)
			if err != nil {
				return nil, err
			}
			arr[i] = val
		}
		pos += n
		if pos+2*m > len(fields) {
			break
		}
		pairs := make([][2]int64, m)
		for i := 0; i < m; i++ {
			x, _ := strconv.ParseInt(fields[pos], 10, 64)
			y, _ := strconv.ParseInt(fields[pos+1], 10, 64)
			pairs[i] = [2]int64{x, y}
			pos += 2
		}
		cases = append(cases, testCase{n: n, m: m, arr: arr, pairs: pairs})
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", tc.n, tc.m)
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for _, p := range tc.pairs {
		fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := decodeTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
