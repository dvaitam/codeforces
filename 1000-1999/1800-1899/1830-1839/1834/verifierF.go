package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func precomputeF(p []int) []int {
	n := len(p)
	diff := make([]int, n+1)
	for j, v := range p {
		if v == n {
			continue
		}
		start := (j + 1) % n
		end := (j - v + n) % n
		if start <= end {
			diff[start]++
			diff[end+1]--
		} else {
			diff[start]++
			diff[n]--
			diff[0]++
			diff[end+1]--
		}
	}
	res := make([]int, n)
	cur := 0
	for i := 0; i < n; i++ {
		cur += diff[i]
		res[i] = cur
	}
	return res
}

type queryF struct {
	t int
	k int
}

type testCaseF struct {
	p       []int
	queries []queryF
}

func expectedF(tc testCaseF) []string {
	n := len(tc.p)
	counts := precomputeF(tc.p)
	rev := make([]int, n)
	for i, v := range tc.p {
		rev[n-1-i] = v
	}
	countsRev := precomputeF(rev)

	offset := 0
	reversed := false
	res := []string{}
	if reversed {
		res = append(res, fmt.Sprint(countsRev[offset]))
	} else {
		res = append(res, fmt.Sprint(counts[offset]))
	}
	for _, q := range tc.queries {
		if q.t == 1 {
			offset = (offset + q.k) % n
		} else if q.t == 2 {
			offset = (offset - q.k) % n
			if offset < 0 {
				offset += n
			}
		} else {
			reversed = !reversed
			offset = (n - offset) % n
		}
		if reversed {
			res = append(res, fmt.Sprint(countsRev[offset]))
		} else {
			res = append(res, fmt.Sprint(counts[offset]))
		}
	}
	return res
}

func genTestsF() []testCaseF {
	rand.Seed(6)
	tests := make([]testCaseF, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(6) + 1
		p := rand.Perm(n)
		for i := range p {
			p[i]++
		}
		qnum := rand.Intn(5) + 1
		queries := make([]queryF, qnum)
		for j := 0; j < qnum; j++ {
			t := rand.Intn(3) + 1
			k := 0
			if t == 1 || t == 2 {
				k = rand.Intn(n) + 1
			}
			queries[j] = queryF{t: t, k: k}
		}
		tests = append(tests, testCaseF{p: p, queries: queries})
	}
	return tests
}

func runCase(bin string, tc testCaseF) error {
	var input strings.Builder
	n := len(tc.p)
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range tc.p {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprint(v))
	}
	input.WriteByte('\n')
	input.WriteString(fmt.Sprintf("%d\n", len(tc.queries)))
	for _, q := range tc.queries {
		if q.t == 1 || q.t == 2 {
			input.WriteString(fmt.Sprintf("%d %d\n", q.t, q.k))
		} else {
			input.WriteString(fmt.Sprintf("%d\n", q.t))
		}
	}

	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	expect := expectedF(tc)
	if len(gotLines) != len(expect) {
		return fmt.Errorf("expected %d lines got %d", len(expect), len(gotLines))
	}
	for i := range expect {
		if strings.TrimSpace(gotLines[i]) != expect[i] {
			return fmt.Errorf("mismatch on line %d: expected %s got %s", i+1, expect[i], gotLines[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsF()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
