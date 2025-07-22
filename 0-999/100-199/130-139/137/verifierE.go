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

func isVowel(c byte) bool {
	switch c {
	case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
		return true
	}
	return false
}

func compute(s string) string {
	n := len(s)
	A := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if isVowel(s[i-1]) {
			A[i] = A[i-1] + 1
		} else {
			A[i] = A[i-1] - 2
		}
	}
	comp := make([]int, n+1)
	vals := append([]int(nil), A...)
	sort.Ints(vals)
	uniq := []int{vals[0]}
	for _, v := range vals[1:] {
		if v != uniq[len(uniq)-1] {
			uniq = append(uniq, v)
		}
	}
	m := len(uniq)
	for i := 0; i <= n; i++ {
		idx := sort.SearchInts(uniq, A[i])
		comp[i] = idx + 1
	}
	INF := n + 5
	tree := make([]int, m+1)
	for i := 1; i <= m; i++ {
		tree[i] = INF
	}
	update := func(pos, v int) {
		for pos <= m {
			if v < tree[pos] {
				tree[pos] = v
			}
			pos += pos & -pos
		}
	}
	query := func(pos int) int {
		res := INF
		for pos > 0 {
			if tree[pos] < res {
				res = tree[pos]
			}
			pos -= pos & -pos
		}
		return res
	}
	update(m-comp[0]+1, 0)
	maxlen := 0
	count := 0
	for r := 1; r <= n; r++ {
		cr := comp[r]
		lmin := query(m - cr + 1)
		if lmin < INF {
			length := r - lmin
			if length > maxlen {
				maxlen = length
				count = 1
			} else if length == maxlen {
				count++
			}
		}
		update(m-comp[r]+1, r)
	}
	if count > 0 {
		return fmt.Sprintf("%d %d", maxlen, count)
	}
	return "No solution"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	b := make([]byte, n)
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(len(letters))]
	}
	s := string(b)
	exp := compute(s)
	return s + "\n", exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
