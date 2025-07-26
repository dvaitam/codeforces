package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	h int
	p []int
}

func minCrystals(h int, p []int) int {
	p = append(append([]int(nil), p...), 0)
	idx := 1
	cur := h
	ans := 0
	for cur > 2 {
		for idx < len(p) && p[idx] >= cur {
			idx++
		}
		if idx >= len(p) {
			break
		}
		if p[idx] == cur-1 {
			next := 0
			if idx+1 < len(p) {
				next = p[idx+1]
			}
			if cur-next > 2 {
				ans++
				cur -= 2
				for idx < len(p) && p[idx] >= cur {
					idx++
				}
			} else {
				cur = next
				idx += 2
			}
		} else {
			cur--
		}
	}
	return ans
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

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	for i := 0; i < 100; i++ {
		h := rng.Intn(50) + 3
		n := rng.Intn(10) + 1
		p := make([]int, n)
		cur := h - 1
		for j := 0; j < n; j++ {
			if cur <= 0 {
				cur = 0
			} else {
				cur -= rng.Intn(3) + 1
				if cur < 0 {
					cur = 0
				}
			}
			p[j] = cur
		}
		sort.Sort(sort.Reverse(sort.IntSlice(p)))
		tests = append(tests, testCase{h: h, p: p})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("1\n%d %d\n", t.h, len(t.p))
		for _, v := range t.p {
			input += fmt.Sprintf("%d ", v)
		}
		input = strings.TrimSpace(input) + "\n"
		want := fmt.Sprintf("%d", minCrystals(t.h, append([]int(nil), t.p...)))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\ninput:%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
