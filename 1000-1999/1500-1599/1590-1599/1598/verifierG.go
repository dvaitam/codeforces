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

const md = 2006091501

type testCase struct {
	input    string
	expected string
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func exkmp(S, T string) []int {
	n, m := len(S), len(T)
	nxt := make([]int, m)
	nxt[0] = m
	l, r := 0, 0
	for i := 1; i < m; i++ {
		if i <= r {
			nxt[i] = min(r-i+1, nxt[i-l])
		}
		for i+nxt[i] < m && T[nxt[i]] == T[i+nxt[i]] {
			nxt[i]++
		}
		if i+nxt[i]-1 > r {
			l = i
			r = i + nxt[i] - 1
		}
	}
	le := make([]int, n)
	l, r = 0, 0
	for i := 0; i < n; i++ {
		if i <= r {
			le[i] = min(r-i+1, nxt[i-l])
		}
		for i+le[i] < n && le[i] < m && S[i+le[i]] == T[le[i]] {
			le[i]++
		}
		if i+le[i]-1 > r {
			l = i
			r = i + le[i] - 1
		}
	}
	return le
}

func solveG(S, T string) string {
	n, m := len(S), len(T)
	le := exkmp(S, T)
	ha := make([]int64, n+1)
	mi := make([]int64, n+1)
	mi[0] = 1
	var ht int64
	for i := 0; i < m; i++ {
		ht = (ht*10 + int64(T[i]-'0')) % md
	}
	for i := 0; i < n; i++ {
		mi[i+1] = (mi[i] * 10) % md
		ha[i+1] = (ha[i]*10 + int64(S[i]-'0')) % md
	}
	var ans string
	check := func(a, b, c int) {
		if a < 0 || c >= n || ans != "" {
			return
		}
		h1 := (ha[b+1] - ha[a]*mi[b-a+1]%md + md) % md
		h2 := (ha[c+1] - ha[b+1]*mi[c-b]%md + md) % md
		if (h1+h2)%md != ht {
			return
		}
		ans = fmt.Sprintf("%d %d\n%d %d", a+1, b+1, b+2, c+1)
	}
	for l := 0; l+m <= n; l++ {
		r := l + m - 1
		c := le[l]
		if c == m || S[l+c] > T[c] {
			continue
		}
		z := m - c
		check(l, r, r+z)
		check(l-z, l-1, r)
		z--
		if z > 0 {
			check(l, r, r+z)
			check(l-z, l-1, r)
		}
		if ans != "" {
			return ans
		}
	}
	if ans == "" && m > 1 {
		for l := 0; l+2*(m-1) <= n; l++ {
			check(l, l+m-2, l+2*m-3)
			if ans != "" {
				return ans
			}
		}
	}
	return ans
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	digits := []rune("0123456789")
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 3
		m := rng.Intn(2) + 2
		sbS := make([]rune, n)
		sbT := make([]rune, m)
		for j := 0; j < n; j++ {
			sbS[j] = digits[rng.Intn(10)]
		}
		for j := 0; j < m; j++ {
			sbT[j] = digits[rng.Intn(10)]
		}
		S := string(sbS)
		T := string(sbT)
		input := fmt.Sprintf("%s\n%s\n", S, T)
		exp := solveG(S, T)
		cases[i] = testCase{input: input, expected: exp}
	}
	return cases
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if out != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
