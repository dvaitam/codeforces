package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `
100
5
0 5 1 4 5
2
3 0
1
4
3
2 3 2
2
0 3
4
5 4 4 4
1
4
5
5 4 3 2 0
3
4 3 0
5
5 2 0 4 2
1
4
2
2 5
2
2 3
2
3 2
5
0 3 3 5 1
1
4
4
3 0 4 4
3
5 1 0
2
2 5
5
4 4 0 5 2
4
3 2 2 3
5
2 5 3 0 2
4
0 2 2 1
5
3 0 5 2 3
3
2 5 4
1
2
4
1 3 2 3
4
1 5 5 1
5
2 4 4 0 4
5
3 3 2 1 3
4
4 3 4 0
3
0 4 2
5
2 3 4 5 1
1
4
2
1 3
1
1
1
1
3
4 3 0
4
5 3 1 0
2
0 1
3
2 2 2
2
1 1
2
0 2
4
2 5 2 1
4
5 1 5 0
5
4 1 2 3 4
2
5 5
3
5 3 0
3
1 4 2
4
0 2 4 4
2
1 3
3
1 4 5
3
2 0 0
5
4 4 2 3 0
1
3
2
0 1
3
3 3 1
2
2 1
2
1 3
5
5 3 2 1 5
1
4
3
0 4 2
4
4 1 1 1
2
1 4
4
5 3 4 2
2
1 5
2
2 1
4
3 4 3 5
5
3 2 1 2 4
5
0 2 3 1 3
1
2
5
5 1 5 3 3
4
2 0 5 4
3
5 0 4
4
1 3 2 3
2
4 0
3
1 2 1
1
2
2
4 0
2
2 5
3
5 1 3
4
3 4 5 0
2
1 1
5
4 1 0 5 5
3
4 0 0
2
0 4
4
0 3 4 4
1
2
2
5 4
3
0 1 2
3
1 5 4
3
2 0 4
4
1 5 1 4
4
5 3 0 5
4
3 5 3 2
4
4 2 0 0
4
3 1 0 5
4
0 5 2 4
4
4 2 2 5
2
2 2
5
4 2 0 5 2
5
3 0 4 1 2
4
3 4 5 4
4
4 1 5 5
2
3 5
5
3 3 3 2 4
4
2 4 3 2
2
4 1
5
5 0 0 5 1
5
2 5 2 3 2
5
0 1 0 0 1
1
2
4
1 1 3 5
3
0 3 3
2
2 2
5
1 3 4 2 1
4
1 5 1 0
1
1
3
2 4 2
1
3
4
1 5 1 2
3
4 5 1
1
1
5
1 4 4 4 5
4
4 2 2 0
5
3 0 2 2 4
2
5 3
3
2 0 3
5
2 2 4 1 4
5
5 4 0 4 5
3
4 3 4
2
3 5
5
1 4 2 4 5
5
5 4 5 4 2
4
4 3 5 1
5
0 2 5 4 1
1
1
3
1 2 2
2
0 0
3
2 2 1
1
2
4
4 5 5 1
5
1 0 4 2 3
4
0 1 4 1
4
0 3 3 3
4
0 5 5 0
3
0 3 1
1
1
1
2
3
3 0 2
4
3 4 0 3
5
5 5 1 5 1
3
2 1 2
3
2 4 2
3
5 5 0
2
4 1
4
1 2 2 3
3
0 5 5
5
3 1 2 4 3
3
0 3 4
4
2 2 2 3
3
1 0 4
4
1 5 3 4
5
3 4 0 5 4
1
4
5
4 1 5 0 5
2
5 5
1
4
4
2 1 0 1
2
1 2
5
2 0 2 0 0
3
3 3 1
5
2 5 2 4 5
4
1 0 4 3
2
1 3
4
4 3 2 5
4
5 5 4 5
3
1 1 4
1
4
3
5 3 4
2
1 5
1
1
3
2 1 3
4
4 0 1 3
5
2 1 1 2 2
4
4 3 5 1
4
0 0 2 0
3
5 0 0
1
4
5
5 0 5 3 1
4
2 2 5 1
5
0 2 5 0 4
5
0 5 2 4 4
2
2 4
3
1 3 3
5
5 2 0 3 5
1
5
3
1 5 5
4
5 3 4 4
4
5 1 4 5
1
2
3
1 2 1
5
4 5 5 0 5
2
2 5
5
3 0 5 2 1
1
5
4
0 4 4 0
2
3 1
1
5
2
0 3
2
1 4
3
5 3 5
3
4 2 0
5
5 4 1 2 5
4
3 4 1 1
5
4 2 4 5 3
2
4 5
3
0 5 5
4
5 2 5 0
2
3 4
5
4 2 0 0 2
4
2 4 3 5
5
1 4 0 3 2
4
1 0 4 1
4
1 4 1 2
2
3 1
4
5 2 3 0 3
5
3 4 4 4 0
2
2 0
2
1 2
4
2 5 3 3
4
2 5 3 0
5
1 1 3 2 2
1
5
2
2 0
5
2 2 1 5 0
2
4 1
5
4 3 4 1 1
4
0 2 4 0
5
3 4 1 4 2
4
4 2 2 2
1
3
5
5 0 4 5 5
4
2 3 5 4
2
3 4
2
0 5
2
3 4
3
2 2 5
3
1 0 2
2
1 2
4
4 5 1 4
2
3 4
4
2 1 3 4
1
2
4
0 0 2 2
4
5 4 0 1
4
0 1 4 3
5
5 3 5 1 3
3
4 1 1
5
5 2 3 2 4
3
2 0 2
2
1 1
5
3 4 0 2 4
2
1 1
2
5 1
5
0 0 0 3 0
3
5 5 3
3
2 1 3
4
2 4 2 2
4
3 1 2 2
1
1
3
3 0 3
2
0 0
3
2 3 4
3
2 5 4
2
3 1
2
1 0
5
0 4 0 5 2
4
4 4 0 2
4
0 2 0 5
5
5 4 4 1 2
4
2 3 1 1
3
4 4 2
1
1
5
4 2 3 4 1
4
3 1 4 0
1
2
3
0 1 2
5
2 4 3 4 0
1
4
1
1
3
4 0 1
5
2 5 1 5 5
1
4
5
5 2 0 0 2
4
4 5 3 1
5
0 4 0 1 1
1
4
3
5 5 0
4
1 1 0 3
1
2
1
1
3
1 1 1
4
1 1 4 4
5
2 0 1 4 0
2
5 1
5
4 4 1 1 5
2
3 4
1
1
1
3
0 1 0
3
2 5 1
4
4 3 1 3
2
3 3
2
2 2
3
3 5 2
5
5 5 2 2 1
1
2
3
4 1 4
3
5 5 5
1
1
2
2 0
4
2 0 4 5
4
5 5 0 4
4
2 3 5 4
1
4
3
2 1 1
5
3 2 4 4 2
3
3 0 2
2
4 4
4
4 3 1 5
1
1
5
5 2 3 3 0
1
2
5
5 5 1 4 3
4
2 4 5 1
5
3 3 4 3 1
5
1 1 3 2 0
5
4 0 0 0 5
3
2 5 5
3
3 1 2
5
1 1 1 1 1
1
1
5
1 3 0 0 1
3
3 0 5
2
1 4
4
4 2 5 4
5
1 3 0 1 1
1
4
5
3 1 0 2 2
5
2 2 0 5 1
2
4 5
3
0 4 1
4
2 2 1 0
5
2 2 3 3 5
4
0 0 4 2
5
5 2 2 4 3
4
3 3 5 2
5
4 0 1 0 5
1
2
4
4 4 5 1
1
3
3
0 5 5
2
1 2
5
2 2 5 0 1
2
3 4
3
3 3 5
3
3 4 5
5
5 5 0 1 4
5
5 2 4 5 0
5
5 2 4 2 4
3
0 3 4
5
2 4 2 4 1
4
3 3 5 5
4
0 0 0 1
4
4 1 5 2
2
2 1
2
2 1
3
5 3 1
5
5 2 5 0 3
5
2 0 5 0 5
2
4 1
4
0 5 4 3
4
1 2 2 2
4
5 3 3 0
2
1 2
1
1
1
1
1
2
4
3
4
4
4
2
1
5
2
2
5
1
4
2
5
5
3
4
1
5
`

type testCase struct {
	input    string
	expected string
}

// Fenwick tree for prefix sums and find by order
type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+1)}
}

func (b *BIT) Update(i, delta int) {
	for i <= b.n {
		b.tree[i] += delta
		i += i & -i
	}
}

func (b *BIT) Sum(i int) int {
	if i > b.n {
		i = b.n
	}
	s := 0
	for i > 0 {
		s += b.tree[i]
		i -= i & -i
	}
	return s
}

func (b *BIT) FindByOrder(k int) int {
	idx := 0
	bitMask := 1 << (bits.Len(uint(b.n)) - 1)
	for d := bitMask; d > 0; d >>= 1 {
		ni := idx + d
		if ni <= b.n && b.tree[ni] < k {
			idx = ni
			k -= b.tree[ni]
		}
	}
	return idx + 1
}

func solveCase(b []int) string {
	n := len(b)
	ok := true
	for i := 1; i < n; i++ {
		if b[i] < b[i-1] {
			ok = false
			break
		}
	}
	if b[0] > n {
		ok = false
	}
	if !ok {
		return "No"
	}
	kval := n + 1
	p := -1
	l, r := 1, n+1
	nxt := make([]int, n+2)
	vis := make([]bool, n+2)
	bExt := append([]int{0}, b...)
	bExt = append(bExt, 0)
	bExt[kval] = kval
	if bExt[1] >= 1 && bExt[1] <= n {
		vis[bExt[1]] = true
	}
	for i := 1; i <= n; i++ {
		if bExt[i] < bExt[i+1] {
			if bExt[i+1] <= n {
				nxt[i] = bExt[i+1]
				vis[nxt[i]] = true
			} else {
				p = i
			}
		}
	}
	for l <= n && bExt[l] <= n {
		l++
	}
	a := make([]int, n+2)
	check := func(kval, pval int) (int, string) {
		cnt := 0
		bit := NewBIT(n)
		for i := n; i >= 1; i-- {
			if nxt[i] != 0 {
				a[i] = a[nxt[i]]
			} else if i >= kval || i == pval {
				cnt++
				a[i] = cnt
			} else {
				c := bit.Sum(bExt[i])
				if c == 0 {
					return -1, ""
				}
				idx := bit.FindByOrder(c)
				a[i] = a[idx]
				bit.Update(idx, -1)
			}
			if !vis[i] {
				bit.Update(i, 1)
			}
		}
		total := bit.Sum(n)
		if total > 0 {
			idx := bit.FindByOrder(total)
			if idx > bExt[1] {
				return 1, ""
			}
		}
		var sb strings.Builder
		sb.WriteString("Yes\n")
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(a[i]))
		}
		return 0, sb.String()
	}
	for l <= r {
		mid := (l + r) >> 1
		res, out := check(mid, p)
		if res == 0 {
			return strings.TrimSpace(out)
		}
		if res == 1 {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return "No"
}

func loadCases() ([]testCase, error) {
	r := strings.NewReader(testcaseData)
	var t int
	if _, err := fmt.Fscan(r, &t); err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	cases := make([]testCase, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		var n int
		if _, err := fmt.Fscan(r, &n); err != nil {
			return nil, fmt.Errorf("case %d: read n: %w", caseIdx+1, err)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(r, &arr[i]); err != nil {
				return nil, fmt.Errorf("case %d: read value: %w", caseIdx+1, err)
			}
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: solveCase(arr),
		})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load cases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
