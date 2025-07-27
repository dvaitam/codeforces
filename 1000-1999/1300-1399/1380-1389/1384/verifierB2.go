package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type caseB struct {
	n int
	k int
	l int
	d []int
}

func generateTests() []caseB {
	r := rand.New(rand.NewSource(3))
	tests := []caseB{
		{2, 1, 3, []int{1, 2}},
		{3, 2, 5, []int{1, 4, 2}},
	}
	for len(tests) < 120 {
		n := r.Intn(12) + 1
		k := r.Intn(4) + 1
		l := r.Intn(20) + 5
		d := make([]int, n)
		for i := 0; i < n; i++ {
			d[i] = r.Intn(l + k + 5)
		}
		tests = append(tests, caseB{n, k, l, d})
	}
	return tests
}

func runCandidate(bin, input string) (string, error) {
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

func solveCase(n, k, l int, d []int) bool {
	p := make([]int, 0, 2*k)
	for i := 0; i < k; i++ {
		p = append(p, i)
	}
	p = append(p, k)
	for i := k - 1; i >= 1; i-- {
		p = append(p, i)
	}
	cycle := len(p)
	type state struct{ pos, tm int }
	queue := []state{{0, 0}}
	vis := make([][]bool, n+2)
	for i := range vis {
		vis[i] = make([]bool, cycle)
	}
	vis[0][0] = true
	for idx := 0; idx < len(queue); idx++ {
		cur := queue[idx]
		pos, tm := cur.pos, cur.tm
		if pos == n+1 {
			return true
		}
		nt := (tm + 1) % cycle
		if pos == 0 || pos == n+1 {
			if !vis[pos][nt] {
				vis[pos][nt] = true
				queue = append(queue, state{pos, nt})
			}
		} else if int64(d[pos-1])+int64(p[nt]) <= int64(l) {
			if !vis[pos][nt] {
				vis[pos][nt] = true
				queue = append(queue, state{pos, nt})
			}
		}
		np := pos + 1
		if np == n+1 {
			if !vis[np][nt] {
				vis[np][nt] = true
				queue = append(queue, state{np, nt})
			}
		} else if np <= n && int64(d[np-1])+int64(p[nt]) <= int64(l) {
			if !vis[np][nt] {
				vis[np][nt] = true
				queue = append(queue, state{np, nt})
			}
		}
	}
	return false
}

func verify(tc caseB, out string) error {
	ans := "No"
	if solveCase(tc.n, tc.k, tc.l, tc.d) {
		ans = "Yes"
	}
	out = strings.TrimSpace(out)
	if strings.EqualFold(out, ans) {
		return nil
	}
	return fmt.Errorf("expected %s, got %s", ans, out)
}

func main() {
	var bin string
	if len(os.Args) == 2 {
		bin = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		bin = os.Args[2]
	} else {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.k, tc.l)
		for j, v := range tc.d {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verify(tc, out); err != nil {
			fmt.Printf("wrong answer on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
