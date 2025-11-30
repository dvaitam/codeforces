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

// Embedded gzipped+base64 testcases from testcasesE.txt.
const encodedTestcases = `
H4sIAHUDK2kC/21UWa7EIAz75xQ9Alm4/9UGEseEp6eq6gzZbRP/5NN41373M5wnvs/kk2G0z/3LhsT/GbbjadvzxErEnn+2rfN8hyJbZI9IQaSxRuRgJmEmDfv57nPknvvVyFKdzsgYdVBtohtFjjuXhfeKDtNb0F/2k3ZpVk8fWNe2R7+DcazjgYIwruNj0bUiiwWa11qIOGqvsBozFzMG9Czto7DRyHHQLqY8T8lHnhUKAlzPiSMmsxjYNUwmH5/nd9pX4CGNZWt9Oti6iGQHSn0JFZPW8EAdY02FZ0fcMYMCSWeNznX2IA8PqRODTlZO/ehSgMdlXcqTirHGlwcOwT+wsz/cKisKuZOrtkdpCos2zLXNvCK/txtkYNiawqFx9tszz0fP0CZiFcotLaW3M1obL9K6SJaMd7yjLlFFqP7i5KpzgX9rO0HYValTyMDdN0ZE9eHLc7J/uVy5Q6DDUkShO2Pi96bk7XFWW20HvOx43XRYrWmq36nit6Zx4tN3ztthcbm4ZYyb7t7H3Gw/Ro2LC9YFAAA=
`

type testCase struct {
	n int
	q int
	k int
	a []int64
	qs [][2]int
}

// solve mirrors 1601E.go.
func solve(tc testCase) string {
	const maxInt64 = int64(^uint64(0) >> 1)

	windowMin := make([]int64, tc.n+1)
	deque := make([]int, 0)
	for i := 1; i <= tc.n; i++ {
		if len(deque) > 0 && deque[0] <= i-tc.k {
			deque = deque[1:]
		}
		for len(deque) > 0 && tc.a[deque[len(deque)-1]] >= tc.a[i] {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, i)
		windowMin[i] = tc.a[deque[0]]
	}

	var out strings.Builder
	for _, q := range tc.qs {
		l, r := q[0], q[1]
		if l > r {
			l, r = r, l
		}
		curMin := tc.a[l]
		ans := curMin
		for s := l + tc.k; s <= r; s += tc.k {
			segMin := windowMin[s]
			if segMin < curMin {
				curMin = segMin
			}
			ans += curMin
		}
		fmt.Fprintln(&out, ans)
	}
	return strings.TrimSpace(out.String())
}

func decodeTestcases() ([]string, error) {
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
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	var cases []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			cases = append(cases, l)
		}
	}
	return cases, nil
}

func parseCase(line string) (testCase, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return testCase{}, fmt.Errorf("invalid case")
	}
	nextInt := func(idx int) (int, error) {
		v, err := strconv.Atoi(fields[idx])
		return v, err
	}
	n, err := nextInt(0)
	if err != nil {
		return testCase{}, err
	}
	q, err := nextInt(1)
	if err != nil {
		return testCase{}, err
	}
	k, err := nextInt(2)
	if err != nil {
		return testCase{}, err
	}
	pos := 3
	if len(fields) < pos+n+2*q {
		return testCase{}, fmt.Errorf("truncated data")
	}
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		val, err := strconv.ParseInt(fields[pos+i-1], 10, 64)
		if err != nil {
			return testCase{}, err
		}
		a[i] = val
	}
	pos += n
	qs := make([][2]int, q)
	for i := 0; i < q; i++ {
		l, err := strconv.Atoi(fields[pos+2*i])
		if err != nil {
			return testCase{}, err
		}
		r, err := strconv.Atoi(fields[pos+2*i+1])
		if err != nil {
			return testCase{}, err
		}
		qs[i] = [2]int{l, r}
	}
	return testCase{n: n, q: q, k: k, a: a, qs: qs}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
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

	lines, err := decodeTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, line := range lines {
		tc, err := parseCase(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expect := solve(tc)
		got, err := runCandidate(bin, line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
