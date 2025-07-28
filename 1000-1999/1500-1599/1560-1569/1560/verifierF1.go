package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCase struct {
	n int
	k int
}

func smallestForLen(length int, digits []int) string {
	sort.Ints(digits)
	first := -1
	for _, d := range digits {
		if d != 0 {
			first = d
			break
		}
	}
	if first == -1 {
		return ""
	}
	res := make([]byte, length)
	res[0] = byte('0' + first)
	fill := digits[0]
	for i := 1; i < length; i++ {
		res[i] = byte('0' + fill)
	}
	return string(res)
}

func attemptSameLen(n string, digits []int) (string, bool) {
	sort.Ints(digits)
	L := len(n)
	res := make([]byte, L)
	var dfs func(pos int, tight bool) bool
	dfs = func(pos int, tight bool) bool {
		if pos == L {
			return true
		}
		nd := int(n[pos] - '0')
		for _, d := range digits {
			if pos == 0 && d == 0 {
				continue
			}
			if tight {
				if d < nd {
					continue
				}
				res[pos] = byte('0' + d)
				if d == nd {
					if dfs(pos+1, true) {
						return true
					}
				} else {
					for i := pos + 1; i < L; i++ {
						res[i] = byte('0' + digits[0])
					}
					return true
				}
			} else {
				res[pos] = byte('0' + d)
				for i := pos + 1; i < L; i++ {
					res[i] = byte('0' + digits[0])
				}
				return true
			}
		}
		return false
	}
	if dfs(0, true) {
		return string(res), true
	}
	return "", false
}

func buildCandidate(n string, digits []int) string {
	if cand, ok := attemptSameLen(n, digits); ok {
		return cand
	}
	return smallestForLen(len(n)+1, digits)
}

func cmpLess(a, b string) bool {
	if len(a) != len(b) {
		return len(a) < len(b)
	}
	return a < b
}

func solve(n, k int) string {
	s := strconv.Itoa(n)
	best := ""
	if k == 1 {
		for d := 1; d <= 9; d++ {
			cand := buildCandidate(s, []int{d})
			if best == "" || cmpLess(cand, best) {
				best = cand
			}
		}
	} else {
		for i := 0; i <= 9; i++ {
			for j := i; j <= 9; j++ {
				digits := []int{i}
				if j != i {
					digits = append(digits, j)
				}
				if len(digits) == 1 && digits[0] == 0 {
					continue
				}
				cand := buildCandidate(s, digits)
				if cand == "" {
					continue
				}
				if best == "" || cmpLess(cand, best) {
					best = cand
				}
			}
		}
	}
	return best
}

func genTests() []testCase {
	r := rand.New(rand.NewSource(6))
	tests := make([]testCase, 100)
	for i := range tests {
		tests[i].n = r.Intn(1000000) + 1
		tests[i].k = r.Intn(2) + 1
	}
	return tests
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d %d\n", tc.n, tc.k)
		want := solve(tc.n, tc.k)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, want, got)
			fmt.Printf("input:\n%s", input)
			return
		}
	}
	fmt.Println("All tests passed")
}
