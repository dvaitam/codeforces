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
	input    string
	expected []int
}

// BIT implementation
func bitAdd(bit []int, idx, delta int) {
	for idx < len(bit) {
		bit[idx] += delta
		idx += idx & -idx
	}
}

func bitSum(bit []int, idx int) int {
	res := 0
	for idx > 0 {
		res += bit[idx]
		idx &= idx - 1
	}
	return res
}

func bitKth(bit []int, k int) int {
	idx := 0
	bitlen := 1
	for bitlen<<1 < len(bit) {
		bitlen <<= 1
	}
	for bitlen > 0 {
		nxt := idx + bitlen
		if nxt < len(bit) && bit[nxt] <= k {
			k -= bit[nxt]
			idx = nxt
		}
		bitlen >>= 1
	}
	return idx + 1
}

func uniqueInts(a []int) []int {
	if len(a) == 0 {
		return a
	}
	j := 0
	for i := 1; i < len(a); i++ {
		if a[i] != a[j] {
			j++
			a[j] = a[i]
		}
	}
	return a[:j+1]
}

func solveCase(l, r []int) []int {
	n := len(l)
	uniq := append([]int(nil), l...)
	sort.Ints(uniq)
	uniq = uniqueInts(uniq)

	pos := make(map[int]int, len(uniq))
	for i, v := range uniq {
		pos[v] = i + 1
	}

	bit := make([]int, len(uniq)+2)
	total := 0
	res := make([]int, n)
	for i := 0; i < n; i++ {
		u := sort.Search(len(uniq), func(j int) bool { return uniq[j] > r[i] }) + 1
		left := bitSum(bit, u-1)
		if total-left > 0 {
			idx := bitKth(bit, left)
			bitAdd(bit, idx, -1)
			total--
		}
		bitAdd(bit, pos[l[i]], 1)
		total++
		res[i] = total
	}
	return res
}

func buildCase(l, r []int) testCase {
	var sb strings.Builder
	n := len(l)
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", l[i], r[i]))
	}
	return testCase{input: sb.String(), expected: solveCase(l, r)}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	l := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		l[i] = rng.Intn(50) + 1
		r[i] = l[i] + rng.Intn(50)
	}
	return buildCase(l, r)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != len(tc.expected) {
		return fmt.Errorf("expected %d numbers got %d", len(tc.expected), len(fields))
	}
	for i, f := range fields {
		var val int
		if _, err := fmt.Sscan(f, &val); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if val != tc.expected[i] {
			return fmt.Errorf("expected %v got %v", tc.expected, fields)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	// simple case
	cases = append(cases, buildCase([]int{1}, []int{1}))
	// large deterministic case
	l := make([]int, 200)
	r := make([]int, 200)
	for i := 0; i < 200; i++ {
		l[i] = i + 1
		r[i] = i + 1
	}
	cases = append(cases, buildCase(l, r))

	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
