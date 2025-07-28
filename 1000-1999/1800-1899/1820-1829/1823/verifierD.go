package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	n   int
	k   int
	pos []int
	cnt []int
}

func possible(n int, k int, pos []int, cnt []int) bool {
	lp, lc := 0, 0
	for i := 0; i < k; i++ {
		if pos[i]-lp < cnt[i]-lc {
			return false
		}
		if pos[i] <= 3 {
			if cnt[i] != pos[i] {
				return false
			}
		} else {
			if cnt[i] < 3 {
				return false
			}
		}
		lp = pos[i]
		lc = cnt[i]
	}
	return true
}

func isPal(s string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

func pCount(s string, m int) int {
	seen := map[string]struct{}{}
	for i := 0; i < m; i++ {
		for j := i + 1; j <= m; j++ {
			sub := s[i:j]
			if isPal(sub) {
				seen[sub] = struct{}{}
			}
		}
	}
	return len(seen)
}

func checkString(s string, n int, pos []int, cnt []int) bool {
	if len(s) != n {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < 'a' || s[i] > 'z' {
			return false
		}
	}
	for i := range pos {
		if pCount(s, pos[i]) != cnt[i] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	rand.Seed(4)
	const cases = 100
	tests := make([]Test, cases)
	for i := range tests {
		n := rand.Intn(7) + 3
		k := rand.Intn(3) + 1
		pos := make([]int, k)
		cnt := make([]int, k)
		lastP, lastC := 0, 0
		for j := 0; j < k; j++ {
			lastP += rand.Intn(n/k+1) + 1
			if lastP > n {
				lastP = n
			}
			pos[j] = lastP
			lastC += rand.Intn(3)
			cnt[j] = lastC
		}
		tests[i] = Test{n: n, k: k, pos: pos, cnt: cnt}
	}

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.pos {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
		for i, v := range tc.cnt {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("error running binary:", err)
		fmt.Print(out.String())
		return
	}

	reader := bufio.NewReader(bytes.NewReader(out.Bytes()))
	for idx, tc := range tests {
		token, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("test %d: failed to read output\n", idx+1)
			return
		}
		token = strings.TrimSpace(token)
		poss := possible(tc.n, tc.k, tc.pos, tc.cnt)
		if token == "NO" {
			if poss {
				fmt.Printf("test %d: expected YES but got NO\n", idx+1)
				return
			}
			continue
		}
		if token != "YES" {
			fmt.Printf("test %d: expected YES/NO got %s\n", idx+1, token)
			return
		}
		sline, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("test %d: failed to read string\n", idx+1)
			return
		}
		s := strings.TrimSpace(sline)
		if !poss {
			fmt.Printf("test %d: expected NO but got YES\n", idx+1)
			return
		}
		if !checkString(s, tc.n, tc.pos, tc.cnt) {
			fmt.Printf("test %d: invalid string\n", idx+1)
			return
		}
	}
	fmt.Printf("verified %d test cases\n", len(tests))
}
