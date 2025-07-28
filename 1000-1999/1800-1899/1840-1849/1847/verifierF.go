package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func solveF(n int, q int, a []int, queries []int) []int {
	orAll := 0
	for _, v := range a {
		orAll |= v
	}
	B := 30
	a2 := make([]int, 2*n)
	copy(a2, a)
	copy(a2[n:], a)
	nxt := make([][]int, B)
	for b := 0; b < B; b++ {
		nxt[b] = make([]int, 2*n+1)
		nxt[b][2*n] = 2 * n
	}
	for i := 2*n - 1; i >= 0; i-- {
		val := a2[i]
		for b := 0; b < B; b++ {
			nxt[b][i] = nxt[b][i+1]
		}
		for b := 0; b < B; b++ {
			if (val>>b)&1 == 1 {
				nxt[b][i] = i
			}
		}
	}
	mp := make(map[int]int)
	for s := 0; s < n; s++ {
		pos := s
		length := 1
		orv := a[pos]
		idx := (length-1)*n + (s + 1)
		if cur, ok := mp[orv]; !ok || idx < cur {
			mp[orv] = idx
		}
		mask := orAll &^ orv
		for mask != 0 && length < n {
			nextpos := 2 * n
			mm := mask
			for mm != 0 {
				lb := mm & -mm
				bit := bits.TrailingZeros(uint(lb))
				np := nxt[bit][pos+1]
				if np < nextpos {
					nextpos = np
				}
				mm -= lb
			}
			if nextpos >= s+n {
				break
			}
			pos = nextpos
			length = pos - s + 1
			orv |= a2[pos]
			idx = (length-1)*n + (s + 1)
			if cur, ok := mp[orv]; !ok || idx < cur {
				mp[orv] = idx
			}
			mask = orAll &^ orv
		}
		if orv != orAll {
			idx = (n-1)*n + (s + 1)
			if cur, ok := mp[orAll]; !ok || idx < cur {
				mp[orAll] = idx
			}
		}
	}
	vals := make([]int, 0, len(mp))
	for v := range mp {
		vals = append(vals, v)
	}
	sort.Ints(vals)
	suf := make([]int, len(vals))
	for i := len(vals) - 1; i >= 0; i-- {
		val := mp[vals[i]]
		if i == len(vals)-1 || suf[i+1] > val {
			suf[i] = val
		} else {
			suf[i] = suf[i+1]
		}
	}
	res := make([]int, q)
	for i, v := range queries {
		j := sort.Search(len(vals), func(k int) bool { return vals[k] > v })
		if j == len(vals) {
			res[i] = -1
		} else {
			res[i] = suf[j]
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(5)
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(10) + 1
		q := rand.Intn(10) + 1
		a := make([]int, n)
		for j := range a {
			a[j] = rand.Intn(1 << 10) // smaller values
		}
		queries := make([]int, q)
		for j := range queries {
			queries[j] = rand.Intn(1 << 10)
		}
		fmt.Fprintf(&input, "%d %d\n", n, q)
		for j := 0; j < n; j++ {
			if j+1 == n {
				fmt.Fprintf(&input, "%d\n", a[j])
			} else {
				fmt.Fprintf(&input, "%d ", a[j])
			}
		}
		for _, v := range queries {
			fmt.Fprintln(&input, v)
		}
		ans := solveF(n, q, a, queries)
		expected[i] = strings.TrimSpace(strings.Join(intSliceToStr(ans), " "))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing binary: %v\n", err)
		os.Exit(1)
	}
	outputs := splitNonEmptyLines(string(out))
	if len(outputs) != t {
		fmt.Fprintf(os.Stderr, "expected %d lines of output, got %d\n", t, len(outputs))
		fmt.Fprint(os.Stderr, string(out))
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		if outputs[i] != expected[i] {
			fmt.Fprintf(os.Stderr, "mismatch on case %d: expected %q got %q\n", i+1, expected[i], outputs[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}

func intSliceToStr(a []int) []string {
	res := make([]string, len(a))
	for i, v := range a {
		res[i] = strconv.Itoa(v)
	}
	return res
}

func splitNonEmptyLines(s string) []string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	res := lines[:0]
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			res = append(res, line)
		}
	}
	return res
}
