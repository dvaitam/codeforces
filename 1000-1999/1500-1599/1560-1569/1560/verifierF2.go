package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n int
	k int
}

func countDistinct(mask int) int {
	return bits.OnesCount(uint(mask))
}

func nextBeautiful(n int, k int) int {
	s := strconv.Itoa(n)
	L := len(s)
	prefixMask := make([]int, L+1)
	for i := 0; i < L; i++ {
		d := int(s[i] - '0')
		prefixMask[i+1] = prefixMask[i] | (1 << d)
	}
	if countDistinct(prefixMask[L]) <= k {
		return n
	}
	best := int64(1<<63 - 1)
	for j := L - 1; j >= 0; j-- {
		mask := prefixMask[j]
		if countDistinct(mask) > k {
			continue
		}
		start := int(s[j]-'0') + 1
		for d := start; d < 10; d++ {
			m2 := mask | (1 << d)
			if countDistinct(m2) > k {
				continue
			}
			fill := byte('0')
			for t := 0; t < 10; t++ {
				m3 := m2
				if (m3 & (1 << t)) == 0 {
					m3 |= 1 << t
				}
				if countDistinct(m3) <= k {
					fill = byte('0' + t)
					break
				}
			}
			cand := s[:j] + string('0'+byte(d)) + strings.Repeat(string(fill), L-j-1)
			if val, err := strconv.ParseInt(cand, 10, 64); err == nil {
				if val < best {
					best = val
				}
			}
		}
	}
	if best != int64(1<<63-1) {
		return int(best)
	}
	if k == 1 {
		ans := 0
		for i := 0; i < L+1; i++ {
			ans = ans*10 + 1
		}
		return ans
	}
	pow := 1
	for i := 0; i < L; i++ {
		pow *= 10
	}
	return pow
}

func genTests() []testCase {
	r := rand.New(rand.NewSource(7))
	tests := make([]testCase, 100)
	for i := range tests {
		tests[i].n = r.Intn(1000000) + 1
		tests[i].k = r.Intn(10) + 1
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
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d %d\n", tc.n, tc.k)
		want := fmt.Sprintf("%d", nextBeautiful(tc.n, tc.k))
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
