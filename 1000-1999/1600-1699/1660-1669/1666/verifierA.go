package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func dirIdx(c byte) int {
	switch c {
	case 'U':
		return 0
	case 'D':
		return 1
	case 'L':
		return 2
	default:
		return 3
	}
}

func solve(s string) string {
	n := len(s)
	pref := make([][4]int, n+1)
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i]
		pref[i+1][dirIdx(s[i])]++
	}

	divisors := make([][]int, n+1)
	for l := 1; l <= n; l++ {
		for d := 1; d*d <= l; d++ {
			if l%d == 0 {
				divisors[l] = append(divisors[l], d)
				if d*d != l {
					divisors[l] = append(divisors[l], l/d)
				}
			}
		}
	}

	balanced := func(l, r int) bool {
		u := pref[r][0] - pref[l][0]
		d := pref[r][1] - pref[l][1]
		if u != d {
			return false
		}
		L := pref[r][2] - pref[l][2]
		R := pref[r][3] - pref[l][3]
		return L == R
	}

	check := func(start, length, rows, cols int) bool {
		dirs := []struct{ dr, dc int }{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		indeg := make([]byte, length)
		for idxAbs := 0; idxAbs < length; idxAbs++ {
			ch := s[start+idxAbs]
			var dir int
			switch ch {
			case 'U':
				dir = 0
			case 'D':
				dir = 1
			case 'L':
				dir = 2
			default:
				dir = 3
			}
			r := idxAbs / cols
			c := idxAbs % cols
			nr := r + dirs[dir].dr
			nc := c + dirs[dir].dc
			if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
				return false
			}
			indeg[nr*cols+nc]++
		}
		for _, v := range indeg {
			if v != 1 {
				return false
			}
		}
		return true
	}

	ans := 0
	for l := 0; l < n; l++ {
		for r := l + 1; r <= n; r++ {
			length := r - l
			if length%2 == 1 {
				continue
			}
			if !balanced(l, r) {
				continue
			}
			ok := false
			for _, rows := range divisors[length] {
				cols := length / rows
				if rows%2 == 1 && cols%2 == 1 {
					continue
				}
				if check(l, length, rows, cols) {
					ok = true
					break
				}
			}
			if ok {
				ans++
			}
		}
	}

	return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(1666))
	var tests []test

	samples := []string{
		"RDUL",
		"RDRU",
		"RLRLRL",
	}
	for _, s := range samples {
		inp := s + "\n"
		tests = append(tests, test{inp, solve(s)})
	}

	edge := []string{
		"U",
		"UD",
		"LR",
		"UDLR",
		"RRRR",
		"LLLLUUUUDDDD",
	}
	for _, s := range edge {
		inp := s + "\n"
		tests = append(tests, test{inp, solve(s)})
	}

	for len(tests) < 120 {
		l := rng.Intn(60) + 1
		sb := strings.Builder{}
		for i := 0; i < l; i++ {
			switch rng.Intn(4) {
			case 0:
				sb.WriteByte('U')
			case 1:
				sb.WriteByte('D')
			case 2:
				sb.WriteByte('L')
			default:
				sb.WriteByte('R')
			}
		}
		s := sb.String()
		inp := s + "\n"
		tests = append(tests, test{inp, solve(s)})
	}

	return tests
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
