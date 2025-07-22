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

func randWord(r *rand.Rand) string {
	length := r.Intn(5) + 1 // 1..5 letters
	b := make([]byte, length)
	for i := range b {
		b[i] = byte('a' + r.Intn(26))
	}
	return string(b)
}

func generateTests() []string {
	r := rand.New(rand.NewSource(43))
	tests := make([]string, 100)
	for t := 0; t < 100; t++ {
		n := r.Intn(4) + 1
		lesha := make([]string, 0, n)
		seen := map[string]bool{}
		for len(lesha) < n {
			w := randWord(r)
			if !seen[w] {
				seen[w] = true
				lesha = append(lesha, w)
			}
		}
		m := r.Intn(10) + 1
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		sb.WriteString(strings.Join(lesha, " "))
		sb.WriteByte('\n')
		sb.WriteString(strconv.Itoa(m))
		sb.WriteByte('\n')
		for i := 0; i < m; i++ {
			k := r.Intn(20) + 1
			words := make([]string, k)
			for j := 0; j < k; j++ {
				words[j] = randWord(r)
			}
			sb.WriteString(strconv.Itoa(k))
			sb.WriteByte(' ')
			sb.WriteString(strings.Join(words, " "))
			sb.WriteByte('\n')
		}
		tests[t] = sb.String()
	}
	return tests
}

func permutations(n int) [][]int {
	res := [][]int{}
	perm := make([]int, n)
	used := make([]bool, n)
	var dfs func(int)
	dfs = func(pos int) {
		if pos == n {
			cp := make([]int, n)
			copy(cp, perm)
			res = append(res, cp)
			return
		}
		for i := 0; i < n; i++ {
			if !used[i] {
				used[i] = true
				perm[pos] = i
				dfs(pos + 1)
				used[i] = false
			}
		}
	}
	dfs(0)
	return res
}

func solve(input string) string {
	f := strings.Fields(input)
	idx := 0
	readInt := func() int {
		v, _ := strconv.Atoi(f[idx])
		idx++
		return v
	}
	n := readInt()
	lesha := make([]string, n)
	for i := 0; i < n; i++ {
		lesha[i] = f[idx]
		idx++
	}
	m := readInt()
	perms := permutations(n)
	totalInv := n * (n - 1) / 2
	bestP := -1
	bestIdx := -1
	for problemIdx := 1; problemIdx <= m; problemIdx++ {
		k := readInt()
		archive := make([]string, k)
		for i := 0; i < k; i++ {
			archive[i] = f[idx]
			idx++
		}
		minInv := totalInv + 1
		for _, p := range perms {
			pos := 0
			for _, w := range archive {
				if w == lesha[p[pos]] {
					pos++
					if pos == n {
						break
					}
				}
			}
			if pos == n {
				inv := 0
				for i := 0; i < n; i++ {
					for j := i + 1; j < n; j++ {
						if p[i] > p[j] {
							inv++
						}
					}
				}
				if inv < minInv {
					minInv = inv
				}
			}
		}
		if minInv <= totalInv {
			p := totalInv - minInv + 1
			if p > bestP {
				bestP = p
				bestIdx = problemIdx
			}
		}
	}
	if bestIdx == -1 {
		return "Brand new problem!"
	}
	return fmt.Sprintf("%d\n:%s:", bestIdx, strings.Repeat("|", bestP))
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, test := range tests {
		expected := solve(test)
		actual, err := runBinary(bin, test)
		if err != nil {
			fmt.Printf("Test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if actual != expected {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, test, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
