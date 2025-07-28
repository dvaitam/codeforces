package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type BIT struct {
	n    int
	tree []int64
}

func NewBIT(n int) *BIT {
	b := &BIT{n: n + 2, tree: make([]int64, n+3)}
	return b
}

func (b *BIT) Add(i int, val int64) {
	for i <= b.n {
		b.tree[i] += val
		i += i & -i
	}
}

func (b *BIT) Sum(i int) int64 {
	if i > b.n {
		i = b.n
	}
	var s int64
	for i > 0 {
		s += b.tree[i]
		i -= i & -i
	}
	return s
}

func sumSubarrayMin(a []int64) int64 {
	n := len(a)
	prev := make([]int, n)
	stack := make([]int, 0, n)
	for i := 0; i < n; i++ {
		for len(stack) > 0 && a[stack[len(stack)-1]] > a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			prev[i] = -1
		} else {
			prev[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	next := make([]int, n)
	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			next[i] = n
		} else {
			next[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	var res int64
	for i := 0; i < n; i++ {
		left := i - prev[i]
		right := next[i] - i
		res += a[i] * int64(left) * int64(right)
	}
	return res
}

func solveE(r *bufio.Reader) string {
	var T int
	fmt.Fscan(r, &T)
	var out strings.Builder
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(r, &n)
		var s string
		fmt.Fscan(r, &s)
		pref := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			if s[i-1] == '(' {
				pref[i] = pref[i-1] + 1
			} else {
				pref[i] = pref[i-1] - 1
			}
		}
		vals := make([]int64, len(pref))
		copy(vals, pref)
		sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
		uniq := vals[:0]
		for _, v := range vals {
			if len(uniq) == 0 || uniq[len(uniq)-1] != v {
				uniq = append(uniq, v)
			}
		}
		idx := make(map[int64]int, len(uniq))
		for i, v := range uniq {
			idx[v] = i + 1
		}
		bitCount := NewBIT(len(uniq))
		bitSum := NewBIT(len(uniq))
		totalSum := int64(0)
		T1 := int64(0)
		id0 := idx[pref[0]]
		bitCount.Add(id0, 1)
		bitSum.Add(id0, pref[0])
		totalSum += pref[0]
		for r := 1; r <= n; r++ {
			val := pref[r]
			id := idx[val]
			cntLe := bitCount.Sum(id)
			sumLe := bitSum.Sum(id)
			sumGt := totalSum - sumLe
			T1 += val*cntLe + sumGt
			bitCount.Add(id, 1)
			bitSum.Add(id, val)
			totalSum += val
		}
		totalMin := sumSubarrayMin(pref) - func(arr []int64) int64 {
			var s int64
			for _, v := range arr {
				s += v
			}
			return s
		}(pref)
		result := T1 - totalMin
		out.WriteString(fmt.Sprintf("%d\n", result))
	}
	return strings.TrimSpace(out.String())
}

func generateCaseE(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	var s strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			s.WriteByte('(')
		} else {
			s.WriteByte(')')
		}
	}
	return fmt.Sprintf("1\n%d\n%s\n", n, s.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseE(rng)
		expect := solveE(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
