package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded gzipped+base64 testcases from testcasesA.txt.
const encodedTestcases = `
H4sIAEwJK2kC/21Uy7HEMAi7u4qUwM/9t/ZskDA78w7JbgBjJAHx+SeffnaeOP/vY58vO9ZjW/vbn2VEfmfExq9f34lwnBdk0WOPzGHw5JORMiy7Mnx77Tz17rhWZ8yKY5esQlbGZFU3h+G9s6pbU2Q1fiKkI6TR3NO6DHfdmrSzStYr6bPhc2ATVOTDTiR+b217AImBGa9n+cdctL9bDPfeKhlXmbW+GkMkO5UnMioQc9ETo4NrQYbMfs4mi+C5UBaOsjmiqPf1Sn39oHuoJ7pArLcuDt2RI7WjlzVUr90brqrFv6KP0EmZJxK7t0IG9MWJUO/lYNe6L4vTGF38eOXtN0LBgzeLMbiJfHzVhBgUVmq59uucqq2Vzdq696Tr0ZHBWdFiBvbrryVQvzYL3hppR1JFxI67a3revJZ2xaugcnITyzH9I2+zy93AjuIMVkfuwcB7K3pLGx0nelb4n8ezN2zsD24Ibc21bdwB1n0n3CqYLukpzOif+ZsbTbDn2H+7dwcV1Mb8eGQtkedqip4q6MPsARsT02+glTe5mc17P3CCvPFH9zY7v3bS3IOcParIXabYzAZvQMGTFf4/0HQCfSYGAAA=
`

type testCase struct {
	n int
	m int
	f []int
	edges [][2]int
}

// solve mirrors 164A.go.
func solve(tc testCase) string {
	n := tc.n
	maxNode := n
	for _, e := range tc.edges {
		if e[0] > maxNode {
			maxNode = e[0]
		}
		if e[1] > maxNode {
			maxNode = e[1]
		}
	}
	f := make([]int, maxNode+1)
	copy(f[1:], tc.f)
	adj := make([][]int, maxNode+1)
	radj := make([][]int, maxNode+1)
	for _, e := range tc.edges {
		a, b := e[0], e[1]
		adj[a] = append(adj[a], b)
		radj[b] = append(radj[b], a)
	}
	fwd := make([]bool, maxNode+1)
	q := make([]int, 0, maxNode)
	for i := 1; i <= maxNode && i <= n; i++ {
		if f[i] == 1 {
			fwd[i] = true
			q = append(q, i)
		}
	}
	for qi := 0; qi < len(q); qi++ {
		u := q[qi]
		for _, v := range adj[u] {
			if f[v] == 1 {
				continue
			}
			if !fwd[v] {
				fwd[v] = true
				q = append(q, v)
			}
		}
	}
	bakNonAssign := make([]bool, maxNode+1)
	q = q[:0]
	for i := 1; i <= maxNode && i <= n; i++ {
		if f[i] == 2 {
			bakNonAssign[i] = true
			q = append(q, i)
		}
	}
	for qi := 0; qi < len(q); qi++ {
		u := q[qi]
		for _, v := range radj[u] {
			if f[v] == 1 {
				continue
			}
			if !bakNonAssign[v] {
				bakNonAssign[v] = true
				q = append(q, v)
			}
		}
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		interesting := false
		if fwd[i] {
			if f[i] == 1 {
				for _, v := range adj[i] {
					if bakNonAssign[v] {
						interesting = true
						break
					}
				}
			} else if bakNonAssign[i] {
				interesting = true
			}
		}
		if interesting {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
		if i < n {
			sb.WriteByte(' ')
		}
	}
	return sb.String()
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
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	pos := 0
	// first token is number of cases
	_, err = strconv.Atoi(fields[pos])
	if err != nil {
		return nil, err
	}
	pos++
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
		fRaw := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[pos+j])
			if err != nil {
				return nil, err
			}
			fRaw[j] = v
		}
		pos += n
		if pos+2*m > len(fields) {
			break
		}
		edges := make([][2]int, m)
		maxNode := n
		for j := 0; j < m; j++ {
			a, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, err
			}
			b, err := strconv.Atoi(fields[pos+1])
			if err != nil {
				return nil, err
			}
			edges[j] = [2]int{a, b}
			pos += 2
			if a > maxNode {
				maxNode = a
			}
			if b > maxNode {
				maxNode = b
			}
		}
		f := make([]int, maxNode)
		copy(f, fRaw)
		cases = append(cases, testCase{n: maxNode, m: m, f: f, edges: edges})
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i, v := range tc.f {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
