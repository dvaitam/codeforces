package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runProgram(binary, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func buildSA(s string) []int {
	n := len(s)
	sa := make([]int, n)
	rank := make([]int, n)
	tmp := make([]int, n)
	for i := 0; i < n; i++ {
		sa[i] = i
		rank[i] = int(s[i])
	}
	for k := 1; k < n; k <<= 1 {
		sort.Slice(sa, func(i, j int) bool {
			a, b := sa[i], sa[j]
			if rank[a] != rank[b] {
				return rank[a] < rank[b]
			}
			ra, rb := -1, -1
			if a+k < n {
				ra = rank[a+k]
			}
			if b+k < n {
				rb = rank[b+k]
			}
			return ra < rb
		})
		tmp[sa[0]] = 0
		for i := 1; i < n; i++ {
			prev, cur := sa[i-1], sa[i]
			prev2, cur2 := -1, -1
			if prev+k < n {
				prev2 = rank[prev+k]
			}
			if cur+k < n {
				cur2 = rank[cur+k]
			}
			if rank[prev] != rank[cur] || prev2 != cur2 {
				tmp[cur] = tmp[prev] + 1
			} else {
				tmp[cur] = tmp[prev]
			}
		}
		copy(rank, tmp)
		if rank[sa[n-1]] == n-1 {
			break
		}
	}
	return sa
}

func buildLCP(s string, sa []int) []int {
	n := len(s)
	rank := make([]int, n)
	for i := 0; i < n; i++ {
		rank[sa[i]] = i
	}
	lcp := make([]int, n)
	h := 0
	for i := 0; i < n; i++ {
		if rank[i] > 0 {
			j := sa[rank[i]-1]
			for i+h < n && j+h < n && s[i+h] == s[j+h] {
				h++
			}
			lcp[rank[i]] = h
			if h > 0 {
				h--
			}
		}
	}
	return lcp
}

func solve(s string, k int64) string {
	n := len(s)
	sa := buildSA(s)
	lcp := buildLCP(s, sa)
	var total int64
	for i := 0; i < n; i++ {
		suffixLen := int64(n - sa[i])
		common := int64(lcp[i])
		cnt := suffixLen - common
		if total+cnt >= k {
			t := k - total
			length := int(common + t)
			return s[sa[i] : sa[i]+length]
		}
		total += cnt
	}
	return "No such line."
}

func randomTest(rng *rand.Rand) (string, int64) {
	l := rng.Intn(10) + 1
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('a' + rng.Intn(3))
	}
	s := string(b)
	maxK := l*(l+1)/2 + 1
	k := int64(rng.Intn(maxK) + 1)
	return s, k
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		s, k := randomTest(rng)
		input := fmt.Sprintf("%s\n%d\n", s, k)
		expected := solve(s, k)
		out, err := runProgram(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(out) != expected {
			fmt.Printf("Test %d failed: expected %s got %s\n", i+1, expected, strings.TrimSpace(out))
			return
		}
	}
	fmt.Println("All tests passed")
}
