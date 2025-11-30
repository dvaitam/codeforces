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

// Embedded gzipped+base64 testcases from testcasesG.txt.
const encodedTestcases = `
H4sIAI4EK2kC/1VU2bHEIAz7pwpKWF8Q+m/sgSycvNlJhgUfkiziPbp262O/vY9m/fy0yw9L2Y+32Eex1w+CAv9Pku89P6F5bH1h25iw6yFEqtTCie/VRMTiiXZtzqKBwqdprhQR0YAKkdk+gDT6bH7jkIOOPyAFMkeWowpYaTWM3P4B28CG3g2FGIpo45Fv2IZGWdbAf+/zHI0v2uS1OwbECZIblCDfglwDjlEyMSnhoTGgnGcVykD9yYkdkDfGqYMmHKdEB+sqWZ8sUQPRRPRBkwM5Q8SUS7n0Q2YctEkWkNpllZEPIUn1la3gqanFVUuqtM1RiuMbnLTnhlxdDKBAiOJaT1Umnbhbl7hxK/jHgJJ9NqvLd4HvOUkTK04sXZSELzVlktRVkCoSnE6SkuZ0XkDEh3jJyOCkZL9K39FSXfmUyzLxOqSUnVRfwCDXQfsop/yUk1/BJGdYpnlXo0TTjKr5jna9cu/Q/97neYrR+52w4kXHIGfWrXmIJWp4VrPR8i90KLtHzSv1sVJt0v5B9qtQC7VzmEiJy8s0NIl9pLh3EPMKxhowyv1MZqy0P8ffpetMBQAA
`

type edge struct {
	u, v int
	w    int
}

type testCase struct {
	n     int
	edges []edge
}

type dsu struct {
	parent []int
	size   []int
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n+1), size: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
}

func (d *dsu) connected(n int) bool {
	root := d.find(1)
	for i := 2; i <= n; i++ {
		if d.find(i) != root {
			return false
		}
	}
	return true
}

func solve(tc testCase) string {
	cur := tc.edges
	ans := 0
	for bit := 30; bit >= 0; bit-- {
		d := newDSU(tc.n)
		nextEdges := make([]edge, 0, len(cur))
		mask := 1 << uint(bit)
		for _, e := range cur {
			if e.w&mask == 0 {
				d.union(e.u, e.v)
				nextEdges = append(nextEdges, e)
			}
		}
		if d.connected(tc.n) {
			cur = nextEdges
		} else {
			ans |= mask
		}
	}
	return fmt.Sprint(ans)
}

func decodeTestcases() (string, error) {
	data, err := base64.StdEncoding.DecodeString(encodedTestcases)
	if err != nil {
		return "", err
	}
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer r.Close()
	var out bytes.Buffer
	if _, err := out.ReadFrom(r); err != nil {
		return "", err
	}
	return out.String(), nil
}

func parseTestcases() ([]testCase, error) {
	raw, err := decodeTestcases()
	if err != nil {
		return nil, err
	}
	var cases []testCase
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tokens := strings.Fields(line)
		if len(tokens) < 2 {
			continue
		}
		ints := make([]int, 0, len(tokens))
		for _, tok := range tokens {
			val, err := strconv.Atoi(tok)
			if err != nil {
				return nil, err
			}
			ints = append(ints, val)
		}
		n := ints[0]
		m := ints[1]
		needed := 2 + 3*m
		for len(ints) < needed {
			ints = append(ints, 0) // pad missing edge data with zeros
		}
		edges := make([]edge, m)
		maxNode := n
		pos := 2
		for i := 0; i < m; i++ {
			u := ints[pos]
			v := ints[pos+1]
			w := ints[pos+2]
			edges[i] = edge{u: u, v: v, w: w}
			if u > maxNode {
				maxNode = u
			}
			if v > maxNode {
				maxNode = v
			}
			pos += 3
		}
		cases = append(cases, testCase{n: maxNode, edges: edges})
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1 %d %d\n", tc.n, len(tc.edges))
	for _, e := range tc.edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
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
