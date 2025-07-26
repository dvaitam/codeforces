package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func expected(a, b string) string {
	n := len(a)
	m := len(b)
	pref := make([]int, m)
	pos := 0
	for i := 0; i < m; i++ {
		for pos < n && a[pos] != b[i] {
			pos++
		}
		if pos == n {
			pref[i] = n
		} else {
			pref[i] = pos
			pos++
		}
	}
	suff := make([]int, m+1)
	suff[m] = n
	pos = n - 1
	for i := m - 1; i >= 0; i-- {
		for pos >= 0 && a[pos] != b[i] {
			pos--
		}
		if pos < 0 {
			suff[i] = -1
		} else {
			suff[i] = pos
			pos--
		}
	}
	bestL, bestR, bestLen := 0, m, m
	r := 0
	for l := 0; l <= m; l++ {
		prefixPos := -1
		if l > 0 {
			if pref[l-1] == n {
				break
			}
			prefixPos = pref[l-1]
		}
		if r < l {
			r = l
		}
		for r <= m {
			if r == m {
				break
			}
			if suff[r] == -1 || suff[r] <= prefixPos {
				r++
				continue
			}
			break
		}
		if r > m {
			break
		}
		delLen := r - l
		if delLen < bestLen {
			bestLen = delLen
			bestL = l
			bestR = r
		}
	}
	res := b[:bestL] + b[bestR:]
	if len(res) == 0 {
		return "-"
	}
	return res
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func randString(rng *rand.Rand, n int) string {
	letters := []rune("abcde")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))

	for t := 0; t < 100; t++ {
		a := randString(rng, rng.Intn(10)+1)
		b := randString(rng, rng.Intn(10)+1)
		input := fmt.Sprintf("%s\n%s\n", a, b)
		exp := expected(a, b)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\noutput:\n%s\n", t+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
