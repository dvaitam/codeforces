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
H4sIAAsHK2kC/z2N0Q0AMQhC/zuOiPuPdhS5hqCWPCtOHUptX0ETlJbrTUZ9Xg6LNpVPuHWF4XuPqQ7LkHtnr056eUYo+m6ntrfb/8LZr934AF4dr5jIAAAA
`

type pair struct{ first, second int64 }

func f(n, o int64) []pair {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		var res []pair
		for n != o {
			res = append(res, pair{0, n})
			res = append(res, pair{n, n})
			n *= 2
		}
		return res
	}
	if n == 2 && o == 2 {
		return []pair{{0, 1}, {1, 1}}
	}
	if (n == 3 && o == 4) || (n == 4 && o == 4) {
		return []pair{{1, 3}, {2, 2}, {0, 4}}
	}
	if n == 5 && o == 8 {
		return []pair{{3, 5}, {2, 2}, {4, 4}, {0, 1}, {1, 1}, {0, 2}, {2, 2}, {0, 4}, {4, 4}, {0, 8}, {0, 8}}
	}
	if n == 6 && o == 8 {
		return []pair{{3, 5}, {2, 6}, {4, 4}, {0, 1}, {1, 1}, {2, 2}, {0, 4}, {4, 4}, {0, 8}, {0, 8}}
	}
	if n == o {
		return f(n-1, o)
	}
	if o >= 2*n {
		res := f(n, o/2)
		nn := n
		for nn >= 4 {
			res = append(res, pair{o / 2, o / 2})
			res = append(res, pair{0, o})
			nn -= 2
		}
		if nn == 3 {
			res = append(res, pair{o / 2, o / 2})
			res = append(res, pair{0, o / 2})
			res = append(res, pair{o / 2, o / 2})
			res = append(res, pair{0, o})
		} else if nn == 2 {
			res = append(res, pair{o / 2, o / 2})
			res = append(res, pair{0, o})
		}
		return res
	}
	var res []pair
	for i := o - n; i < o/2; i++ {
		res = append(res, pair{i, o - i})
	}
	res1 := f(o-n-1, o)
	res2 := f(n-o/2, o/2)
	for i := range res2 {
		res2[i].first *= 2
		res2[i].second *= 2
	}
	if n-o/2 > 2 {
		res2 = res2[:len(res2)-1]
		res = append(res, res2...)
		res = append(res, res1...)
		res = append(res, pair{0, o / 2}, pair{o / 2, o / 2}, pair{0, o})
		return res
	} else if o-n-1 > 2 {
		res1 = res1[:len(res1)-1]
		res = append(res, res1...)
		res = append(res, res2...)
		res = append(res, pair{0, o / 2}, pair{o / 2, o / 2}, pair{0, o})
		return res
	}
	panic("unreachable")
}

func nextPow2(n int64) int64 {
	if n&(n-1) == 0 {
		return n
	}
	o := int64(1)
	for o < n {
		o <<= 1
	}
	return o
}

func solve(n int64) string {
	if n == 2 {
		return "-1"
	}
	o := nextPow2(n)
	res := f(n, o)
	var sb strings.Builder
	sb.WriteString(fmt.Sprint(len(res)))
	sb.WriteByte('\n')
	for idx, p := range res {
		if idx > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d %d", p.first, p.second))
	}
	return strings.TrimSpace(sb.String())
}

func decodeTestcases() ([]int64, error) {
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
	parts := strings.Fields(out.String())
	var res []int64
	for _, p := range parts {
		val, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, err
		}
		res = append(res, val)
	}
	return res, nil
}

func runCandidate(bin string, n int64) (string, error) {
	input := fmt.Sprintf("1\n%d\n", n)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := decodeTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, n := range tests {
		expect := solve(n)
		got, err := runCandidate(bin, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
