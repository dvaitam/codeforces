package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type query struct{ l, r, k int }

func solveCase(s string, w int, q query) (int, int) {
	n := len(s)
	digits := make([]int, n+1)
	for i := 1; i <= n; i++ {
		digits[i] = int(s[i-1] - '0')
	}
	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = (pref[i-1] + digits[i]) % 9
	}
	vals := make([][]int, 9)
	for i := 1; i+w-1 <= n; i++ {
		v := (pref[i+w-1] - pref[i-1]) % 9
		if v < 0 {
			v += 9
		}
		if len(vals[v]) < 2 {
			vals[v] = append(vals[v], i)
		}
	}
	l, r, k := q.l, q.r, q.k
	cur := (pref[r] - pref[l-1]) % 9
	if cur < 0 {
		cur += 9
	}
	ans1, ans2 := -1, -1
	for a := 0; a < 9; a++ {
		if len(vals[a]) == 0 {
			continue
		}
		for b := 0; b < 9; b++ {
			if len(vals[b]) == 0 {
				continue
			}
			if (a*cur+b)%9 != k {
				continue
			}
			if a == b {
				if len(vals[a]) < 2 {
					continue
				}
				p, q := vals[a][0], vals[a][1]
				if ans1 == -1 || p < ans1 || (p == ans1 && q < ans2) {
					ans1, ans2 = p, q
				}
			} else {
				p, q := vals[a][0], vals[b][0]
				if ans1 == -1 || p < ans1 || (p == ans1 && q < ans2) {
					ans1, ans2 = p, q
				}
			}
		}
	}
	return ans1, ans2
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct {
		s string
		w int
		q query
	}
	tests := make([]test, 0, 120)
	// deterministic small case
	tests = append(tests, test{s: "12345", w: 2, q: query{l: 1, r: 5, k: 0}})
	for len(tests) < 120 {
		n := rng.Intn(10) + 2
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteByte(byte('0' + rng.Intn(10)))
		}
		s := sb.String()
		w := rng.Intn(n-1) + 1
		l := rng.Intn(n-w+1) + 1
		r := rng.Intn(n-l+1) + l
		k := rng.Intn(9)
		tests = append(tests, test{s: s, w: w, q: query{l: l, r: r, k: k}})
	}

	for i, tc := range tests {
		a, b := solveCase(tc.s, tc.w, tc.q)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%s\n", tc.s))
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.w, 1))
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.q.l, tc.q.r, tc.q.k))
		input := sb.String()
		expected := fmt.Sprintf("%d %d", a, b)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
