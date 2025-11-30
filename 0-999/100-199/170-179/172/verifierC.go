package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesRaw = `5 4 12 8 2 9 2 15 10 2 7
5 2 3 3 2 2 9 1 5 9 5 6 3 10
5 4 6 8 1 2 10 4 10 4 6 5 7
5 1 7 5 7 5 8 3 3 1 2 2 2
1 2 5 10
4 4 7 8 7 5 1 5 5 5 3
4 4 10 4 10 4 10 4 10 4
2 4 6 6 8 1
3 3 1 10 2 7 1
4 1 4 3 6 1 10 2 9 1
1 4 5 4
2 5 8 4 1 6
3 5 9 10 6 7 8
2 2 8 6 6 8
5 2 8 1 5 4 7 3 6 10 9 10
5 1 4 6 10 7 10 1 1 5 9 6
5 4 4 4 5 5 2 1 3 5 9 8
2 1 6 8 2 9
2 3 4 3 9 2
5 1 8 2 1 2 8 2 8 10 3 6
1 4 8 6
1 5 4 2
5 4 10 7 10 1 10 9 10 3 3 3
2 3 6 6 8 8
1 3 6 10
5 1 9 2 4 4 9 3 3 3 10 8
5 4 8 1 8 5 9 1 5 2 10 5
5 3 3 10 3 5 7 7 7 4 2 1
5 1 4 1 6 9 2 9 9 8 2 9
5 4 10 6 3 7 9 2 1 4 1 7
2 1 10 1 1 9
2 2 6 5 3 5
5 1 4 5 6 3 4 4 7 9 6 10
5 3 1 3 9 10 2 4 4 1 6 9
5 1 10 2 9 3 2 1 3 3 3 1
2 2 10 4 1 10
2 1 10 8 5 10
5 1 8 3 7 4 7 3 4 2 6 6
3 3 8 8 10 8 9
2 4 3 6 7 2
4 4 1 1 9 5 8 3 2 1
4 2 8 10 10 6 10 10 4 3
3 2 3 6 6 8 3
1 3 3 3
2 2 9 4 5 8
5 1 7 2 10 3 9 6 6 3 6 2
2 4 9 9 5 7
4 3 7 7 7 9 3 5 3 4
5 2 1 1 7 1 9 1 5 10 3 3
4 2 10 1 4 7 4 7 2 1
1 2 9 10
4 4 5 7 3 1 5 1 9 10
4 1 3 5 5 7 10 9 7 8
4 3 8 3 8 10 7 2 9 7
5 4 10 5 4 3 6 5 1 2 2 8
1 5 5 9
4 3 8 1 9 3 1 6 1 4
3 3 2 3 5 7 3
4 2 7 8 3 4 4 10 5 7
2 1 10 5 5 10
4 1 10 1 1 10 4 10 10 10
3 1 5 10 8 10 5
2 3 9 6 8 2
1 1 9 10
3 3 9 10 4 2 6
4 1 4 7 4 5 7 9 9 10
4 1 9 9 3 3 9 7 8 8
2 5 8 10 5 8
4 4 4 3 6 10 3 6 2 10
2 5 4 2 5 7
4 1 6 9 3 3 10 1 4 1
3 2 10 10 4 8 2
1 3 5 10
5 3 7 9 3 9 6 9 2 10 2 9
5 3 10 1 1 1 7 8 10 10 6 5
4 2 9 2 7 1 7 5 2 7
5 3 3 10 1 7 10 2 9 8 8 2
4 2 2 3 2 2 9 3 3 2
5 2 9 2 10 2 4 3 6 2 6 4
5 4 4 1 10 1 10 2 5 2 6 4
2 5 5 4 2 6
1 1 2 9
4 3 7 7 1 10 8 8 1 5
2 5 6 2 9 4
5 1 3 2 7 4 5 9 5 2 3 8
1 1 2 9
1 2 4 3
2 5 5 7 7 5
5 3 1 1 9 4 2 9 4 9 6 5
2 5 2 7 10 10
5 3 2 8 1 3 9 2 6 9 10 6
5 4 1 1 10 2 3 5 6 9 1 5
2 2 2 2 8 3
4 2 3 9 9 1 6 3 1 4
2 5 2 4 6 9
1 2 1 9
3 4 2 6 6 10 7
3 3 6 2 1 7 4
4 4 10 3 9 6 6 7 1 7
4 1 4 3 1 6 1 4 10 1
1 1 4 3`

type Student struct {
	t   int64
	x   int64
	idx int
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// expected mirrors 172C.go.
func expected(n, m int, students []Student) []int64 {
	ans := make([]int64, n)
	queue := []Student{}
	nextI := 0
	var curTime int64
	processed := 0
	for processed < n {
		if len(queue) == 0 && nextI < n && curTime < students[nextI].t {
			curTime = students[nextI].t
		}
		for nextI < n && students[nextI].t <= curTime && len(queue) < m {
			queue = append(queue, students[nextI])
			nextI++
		}
		for len(queue) < m && nextI < n {
			curTime = students[nextI].t
			queue = append(queue, students[nextI])
			nextI++
		}
		departTime := curTime
		k := len(queue)
		batch := make([]Student, k)
		copy(batch, queue)
		sort.Slice(batch, func(i, j int) bool { return batch[i].x < batch[j].x })
		var unloadTime int64
		lastX := int64(-1)
		var cntAtX int64
		for _, st := range batch {
			if st.x != lastX {
				if lastX != -1 {
					unloadTime += 1 + cntAtX/2
				}
				lastX = st.x
				cntAtX = 1
			} else {
				cntAtX++
			}
		}
		if lastX != -1 {
			unloadTime += 1 + cntAtX/2
		}
		maxX := batch[k-1].x
		for _, st := range batch {
			ans[st.idx] = departTime + st.x
		}
		curTime = departTime + maxX + unloadTime + maxX
		processed += k
		queue = queue[:0]
	}
	return ans
}

func parseTestcases() ([]struct {
	n int
	m int
	s []Student
}, error) {
	tokens := strings.Fields(testcasesRaw)
	var cases []struct {
		n int
		m int
		s []Student
	}
	pos := 0
	for pos < len(tokens) {
		if pos+2 > len(tokens) {
			return nil, fmt.Errorf("truncated tokens at %d", pos)
		}
		n, err := strconv.Atoi(tokens[pos])
		if err != nil {
			return nil, fmt.Errorf("parse n at token %d: %v", pos, err)
		}
		m, err := strconv.Atoi(tokens[pos+1])
		if err != nil {
			return nil, fmt.Errorf("parse m at token %d: %v", pos+1, err)
		}
		pos += 2
		if pos+2*n > len(tokens) {
			return nil, fmt.Errorf("not enough tokens for case with n=%d", n)
		}
		students := make([]Student, n)
		for i := 0; i < n; i++ {
			ti, err := strconv.Atoi(tokens[pos+2*i])
			if err != nil {
				return nil, fmt.Errorf("parse t[%d]: %v", i, err)
			}
			xi, err := strconv.Atoi(tokens[pos+2*i+1])
			if err != nil {
				return nil, fmt.Errorf("parse x[%d]: %v", i, err)
			}
			students[i] = Student{t: int64(ti), x: int64(xi), idx: i}
		}
		pos += 2 * n
		cases = append(cases, struct {
			n int
			m int
			s []Student
		}{n: n, m: m, s: students})
	}
	return cases, nil
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
		want := expected(tc.n, tc.m, tc.s)
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		for _, st := range tc.s {
			fmt.Fprintf(&input, "%d %d\n", st.t, st.x)
		}
		gotStr, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		outFields := strings.Fields(gotStr)
		if len(outFields) != tc.n {
			fmt.Fprintf(os.Stderr, "case %d wrong output length\n", idx+1)
			os.Exit(1)
		}
		for i := 0; i < tc.n; i++ {
			val, err := strconv.ParseInt(outFields[i], 10, 64)
			if err != nil || val != want[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at student %d: expected %d got %s\n", idx+1, i+1, want[i], outFields[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
