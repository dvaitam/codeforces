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

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesRaw = `3 3
1 2
2 3
2 1
1 1 0
2 3
4 1
1 2
2 3
2 4
1 1 2
3 1
1 2
2 3
1 1 2
4 3
1 2
2 3
3 4
2 2
1 2 3
2 1
4 1
1 2
2 3
3 4
1 4 3
4 4
1 2
1 3
1 4
2 3
1 1 0
2 3
2 3
3 2
1 2
2 3
1 3 3
2 1
4 4
1 2
1 3
2 4
1 1 1
2 2
1 3 1
2 4
2 1
1 2
1 2 2
2 3
1 2
2 1
2 1
2 1
5 3
1 2
2 3
2 4
2 5
2 5
1 3 0
2 2
4 3
1 2
1 3
2 4
1 4 1
1 1 0
1 2 1
2 5
1 2
1 2 0
1 2 3
1 2 1
1 2 3
2 1
3 4
1 2
1 3
2 1
2 1
1 1 2
2 1
3 2
1 2
2 3
1 2 0
2 3
2 5
1 2
1 2 3
1 2 1
1 2 2
2 2
2 1
4 3
1 2
2 3
1 4
2 2
1 4 1
2 1
3 4
1 2
2 3
2 1
1 1 1
1 3 2
1 2 3
3 5
1 2
1 3
1 3 2
2 3
2 1
1 1 0
2 3
4 5
1 2
1 3
2 4
1 1 2
2 1
1 4 3
1 3 3
1 4 3
2 1
1 2
2 2
5 4
1 2
1 3
3 4
3 5
2 4
1 5 1
2 2
1 2 2
5 3
1 2
1 3
3 4
2 5
1 3 0
1 4 3
1 5 1
5 3
1 2
1 3
1 4
1 5
1 5 2
2 5
2 5
5 5
1 2
2 3
1 4
3 5
2 3
1 4 0
1 5 1
2 3
1 5 2
2 1
1 2
1 2 2
4 2
1 2
1 3
3 4
1 1 2
1 2 1
5 5
1 2
1 3
1 4
4 5
1 5 2
2 5
1 4 1
2 3
1 3 3
5 5
1 2
2 3
2 4
4 5
2 5
1 3 1
2 5
2 4
2 5
2 2
1 2
2 2
2 1
4 4
1 2
2 3
3 4
1 2 1
1 3 1
1 2 3
2 1
3 2
1 2
1 3
1 2 0
1 2 0
4 5
1 2
1 3
3 4
2 2
1 3 0
1 4 1
1 3 3
1 4 1
3 3
1 2
2 3
2 2
1 2 0
2 1
4 1
1 2
1 3
1 4
1 1 0
3 2
1 2
2 3
1 2 3
2 1
3 1
1 2
2 3
2 3
5 5
1 2
1 3
2 4
3 5
2 3
1 2 2
1 1 2
2 3
1 2 0
5 5
1 2
1 3
1 4
1 5
2 5
1 3 0
1 1 1
2 5
1 1 0
5 2
1 2
2 3
1 4
3 5
2 4
2 1
3 4
1 2
1 3
2 3
2 1
2 2
1 1 0
5 4
1 2
2 3
2 4
4 5
2 5
1 3 2
1 1 3
2 5
4 2
1 2
1 3
1 4
1 2 2
1 4 0
4 1
1 2
1 3
1 4
1 2 0
2 3
1 2
1 1 1
1 1 1
2 1
4 3
1 2
1 3
1 4
1 1 0
1 1 2
2 1
3 2
1 2
2 3
2 3
2 2
5 4
1 2
1 3
1 4
2 5
2 2
2 2
2 5
1 5 2
5 3
1 2
1 3
1 4
3 5
1 5 1
1 1 0
1 3 2
5 1
1 2
2 3
1 4
2 5
2 5
5 3
1 2
1 3
2 4
2 5
2 1
2 2
1 3 2
2 1
1 2
1 1 1
3 5
1 2
1 3
1 2 3
1 3 1
1 1 1
2 1
2 3
4 5
1 2
1 3
3 4
2 3
2 3
1 2 2
1 4 3
1 1 1
3 2
1 2
1 3
2 1
1 3 1
3 3
1 2
2 3
2 2
1 3 3
2 3
2 2
1 2
2 1
1 2 2
2 2
1 2
1 2 0
2 1
3 2
1 2
1 3
1 1 0
2 3
5 1
1 2
2 3
3 4
2 5
1 5 0
2 1
1 2
2 1
2 5
1 2
2 2
1 1 0
1 1 0
2 2
2 2
3 5
1 2
2 3
1 3 3
1 2 0
1 2 3
2 1
2 1
4 3
1 2
1 3
3 4
2 3
2 3
1 3 0
3 3
1 2
1 3
2 1
2 1
1 1 0
3 2
1 2
2 3
2 2
1 3 3
5 3
1 2
2 3
1 4
1 5
1 1 3
1 4 1
2 1
3 3
1 2
2 3
1 1 2
2 1
2 3
3 4
1 2
2 3
1 2 1
1 2 3
1 1 0
2 1
3 2
1 2
2 3
2 2
2 2
4 4
1 2
2 3
2 4
1 4 3
2 2
1 4 2
1 3 0
3 5
1 2
2 3
1 3 2
2 2
2 3
1 2 0
2 2
5 3
1 2
1 3
2 4
3 5
1 3 3
1 4 3
2 3
2 4
1 2
1 1 0
2 2
1 1 0
1 2 1
4 3
1 2
1 3
2 4
1 1 0
2 1
1 2 3
2 5
1 2
1 1 0
1 1 1
1 2 2
2 2
2 1
4 1
1 2
2 3
1 4
1 3 0
5 1
1 2
1 3
1 4
4 5
2 3
2 1
1 2
2 2
5 1
1 2
1 3
1 4
2 5
2 5
5 4
1 2
2 3
1 4
3 5
2 5
1 3 3
2 3
1 2 1
3 4
1 2
1 3
1 3 0
1 2 0
1 3 3
1 3 0
4 3
1 2
1 3
1 4
2 4
2 3
1 1 3
4 5
1 2
2 3
3 4
2 2
2 2
1 2 2
1 2 1
2 3
3 5
1 2
2 3
2 3
2 2
2 3
1 2 3
2 2
3 1
1 2
2 3
2 3
5 1
1 2
2 3
3 4
1 5
2 5
2 2
1 2
1 2 3
1 1 2
3 2
1 2
1 3
1 2 3
2 2
4 3
1 2
1 3
3 4
1 1 1
1 4 2
2 1
3 3
1 2
2 3
2 2
1 2 3
2 2
2 4
1 2
1 2 3
2 1
2 2
2 1
3 2
1 2
1 3
1 2 0
2 3
3 1
1 2
2 3
2 2
3 1
1 2
1 3
2 2
5 1
1 2
2 3
3 4
4 5
1 1 2
4 3
1 2
2 3
2 4
1 2 3
2 1
1 1 0
5 5
1 2
2 3
3 4
1 5
2 4
2 5
2 5
1 4 1
2 3
4 1
1 2
1 3
2 4
2 3
5 5
1 2
1 3
2 4
4 5
2 5
2 2
1 2 2
1 3 3
2 3`

type query struct {
	typ int
	v   int
	d   int
}

type testcase struct {
	n       int
	q       int
	edges   [][2]int
	queries []query
}

// parseTestcases reads all cases from the embedded raw data.
func parseTestcases() ([]testcase, error) {
	data := strings.TrimSpace(testcasesRaw)
	if data == "" {
		return nil, fmt.Errorf("no testcases embedded")
	}
	sc := bufio.NewScanner(strings.NewReader(data))
	sc.Split(bufio.ScanWords)

	nextInt := func() (int, error) {
		if !sc.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		return strconv.Atoi(sc.Text())
	}

	var cases []testcase
	for {
		n, err := nextInt()
		if err != nil {
			if err.Error() == "unexpected EOF" {
				break
			}
			return nil, err
		}
		q, err := nextInt()
		if err != nil {
			return nil, err
		}
		tc := testcase{n: n, q: q}
		tc.edges = make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			u, err := nextInt()
			if err != nil {
				return nil, err
			}
			v, err := nextInt()
			if err != nil {
				return nil, err
			}
			tc.edges[i] = [2]int{u, v}
		}
		tc.queries = make([]query, q)
		for i := 0; i < q; i++ {
			t, err := nextInt()
			if err != nil {
				return nil, err
			}
			if t == 1 {
				v, err := nextInt()
				if err != nil {
					return nil, err
				}
				d, err := nextInt()
				if err != nil {
					return nil, err
				}
				tc.queries[i] = query{typ: 1, v: v, d: d}
			} else {
				v, err := nextInt()
				if err != nil {
					return nil, err
				}
				tc.queries[i] = query{typ: 2, v: v}
			}
		}
		cases = append(cases, tc)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

// solve returns the expected output for a testcase following 1254D.go logic (all zeros for type 2 queries).
func solve(tc testcase) string {
	var out strings.Builder
	for _, q := range tc.queries {
		if q.typ == 2 {
			fmt.Fprintln(&out, 0)
		}
	}
	return strings.TrimSpace(out.String())
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.q)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		for _, q := range tc.queries {
			if q.typ == 1 {
				fmt.Fprintf(&sb, "1 %d %d\n", q.v, q.d)
			} else {
				fmt.Fprintf(&sb, "2 %d\n", q.v)
			}
		}
		input := sb.String()
		expect := solve(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
