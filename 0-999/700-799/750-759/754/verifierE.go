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

const mod int64 = 1000000007
const base1 int64 = 911382323
const base2 int64 = 972663749

type testCase struct {
	n, m    int
	board   []string
	r, c    int
	pattern []string
}

func solve(tc testCase) []string {
	n, m, r, c := tc.n, tc.m, tc.r, tc.c
	board := tc.board
	pattern := tc.pattern
	N := n + r - 1
	M := m + c - 1
	pow1 := make([]int64, N+1)
	pow2 := make([]int64, M+1)
	pow1[0] = 1
	for i := 1; i <= N; i++ {
		pow1[i] = pow1[i-1] * base1 % mod
	}
	pow2[0] = 1
	for i := 1; i <= M; i++ {
		pow2[i] = pow2[i-1] * base2 % mod
	}
	ans := make([][]bool, n)
	for i := 0; i < n; i++ {
		ans[i] = make([]bool, m)
		for j := range ans[i] {
			ans[i][j] = true
		}
	}
	prefix := make([][]int64, N+1)
	for i := range prefix {
		prefix[i] = make([]int64, M+1)
	}
	for ch := byte('a'); ch <= 'z'; ch++ {
		var patHash int64
		for i := 0; i < r; i++ {
			row := pattern[i]
			for j := 0; j < c; j++ {
				if row[j] == ch {
					patHash = (patHash + pow1[i]*pow2[j]) % mod
				}
			}
		}
		if patHash == 0 {
			continue
		}
		for i := 0; i <= N; i++ {
			for j := range prefix[i] {
				prefix[i][j] = 0
			}
		}
		for i := 0; i < N; i++ {
			rowIdx := i % n
			prev := prefix[i]
			cur := prefix[i+1]
			bRow := board[rowIdx]
			for j := 0; j < M; j++ {
				colIdx := j % m
				var val int64
				if bRow[colIdx] == ch {
					val = pow1[i] * pow2[j] % mod
				}
				cur[j+1] = (cur[j] + prev[j+1] - prev[j] + val) % mod
				if cur[j+1] < 0 {
					cur[j+1] += mod
				}
			}
		}
		for i := 0; i < n; i++ {
			base1Shift := pow1[i]
			for j := 0; j < m; j++ {
				if !ans[i][j] {
					continue
				}
				sub := prefix[i+r][j+c]
				sub -= prefix[i][j+c]
				if sub < 0 {
					sub += mod
				}
				sub -= prefix[i+r][j]
				if sub < 0 {
					sub += mod
				}
				sub += prefix[i][j]
				sub %= mod
				if sub < 0 {
					sub += mod
				}
				expected := patHash * base1Shift % mod * pow2[j] % mod
				if sub != expected {
					ans[i][j] = false
				}
			}
		}
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		buf := make([]byte, m)
		for j := 0; j < m; j++ {
			if ans[i][j] {
				buf[j] = '1'
			} else {
				buf[j] = '0'
			}
		}
		res[i] = string(buf)
	}
	return res
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i := 0; i < tc.n; i++ {
		sb.WriteString(tc.board[i])
		sb.WriteByte('\n')
	}
	fmt.Fprintf(&sb, "%d %d\n", tc.r, tc.c)
	for i := 0; i < tc.r; i++ {
		sb.WriteString(tc.pattern[i])
		sb.WriteByte('\n')
	}
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	expect := solve(tc)
	if len(lines) != len(expect) {
		return fmt.Errorf("wrong number of lines")
	}
	for i := range lines {
		if strings.TrimSpace(lines[i]) != expect[i] {
			return fmt.Errorf("line %d expected %s got %s", i+1, expect[i], lines[i])
		}
	}
	return nil
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	r := rng.Intn(n) + 1
	c := rng.Intn(m) + 1
	board := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			b[j] = byte('a' + rng.Intn(3))
		}
		board[i] = string(b)
	}
	pattern := make([]string, r)
	for i := 0; i < r; i++ {
		b := make([]byte, c)
		for j := 0; j < c; j++ {
			if rng.Intn(4) == 0 {
				b[j] = '?'
			} else {
				b[j] = byte('a' + rng.Intn(3))
			}
		}
		pattern[i] = string(b)
	}
	return testCase{n, m, board, r, c, pattern}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{1, 1, []string{"a"}, 1, 1, []string{"a"}},
		{2, 2, []string{"ab", "ba"}, 1, 1, []string{"?"}},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
