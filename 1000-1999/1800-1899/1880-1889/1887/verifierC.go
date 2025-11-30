package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
4
5 3 -5 -2
3
4 4 -3
2 4 -3
3 3 3
2
-3 4
2
2 2 -3
1 1 1
2
-1 -2
4
2 2 3
2 2 1
2 2 2
1 1 -3
5
3 -1 0 -2 1
2
3 3 3
2 5 2
2
4 4
5
2 2 -3
1 2 1
1 1 1
1 2 -1
2 2 2
1
0
2
1 1 0
1 1 -1
1
2
2
1 1 1
1 1 0
2
4 1
2
1 1 3
1 2 -1
1
5
4
1 1 3
1 1 1
1 1 -3
1 1 -3
4
1 2 -5 0
2
1 4 2
2 2 -2
4
-5 -5 4 1
2
3 4 0
1 4 1
2
3 -3
2
1 1 1
2 2 -1
3
-2 1 5
1
3 3 1
1
-5
3
1 1 3
1 1 -1
1 1 3
2
0 4
1
1 2 -2
3
2 0 5
4
2 2 -3
2 3 0
2 3 1
2 2 -2
2
5 1
5
2 2 -3
1 1 -3
2 2 2
1 2 1
2 2 2
4
1 4 0 5
5
4 4 -2
1 4 3
3 4 3
1 2 1
1 4 -3
3
1 4 -2
5
1 3 1
1 1 1
2 2 -1
1 2 -1
2 3 3
4
-5 5 -2 5
4
3 3 -3
4 4 3
4 4 2
4 4 -3
4
4 4 5 -3
5
4 4 1
3 4 1
2 4 0
4 4 -1
4 4 -2
2
-1 -2
5
2 2 -2
2 2 -1
1 1 3
2 2 1
1 2 0
2
-3 4
3
1 1 1
2 2 -1
2 2 1
2
0 1
3
1 1 -2
1 1 3
2 2 -1
5
-1 2 1 0 5
4
1 1 -3
3 3 3
3 5 2
5 5 -1
5
1 2 0 -2 5
1
1 1 1
1
5
5
1 1 1
1 1 -3
1 1 -3
1 1 3
1 1 -1
2
-3 4
2
1 1 -2
2 2 3
3
-3 4 -1
2
1 1 3
3 3 -2
2
4 -5
3
2 2 3
2 2 2
2 2 -3
2
-2 4
1
2 2 -2
4
3 -5 3 0
2
2 3 -1
4 4 2
3
5 -2 -4
3
2 3 -1
3 3 -2
2 3 -3
2
1 -4
5
1 1 2
2 2 0
2 2 3
1 1 -1
2 2 0
5
-5 0 1 1 3
2
5 5 2
1 3 1
4
3 4 1 -2
4
2 3 -3
1 4 -3
4 4 -3
2 2 -2
1
4
4
1 1 3
1 1 -2
1 1 -2
1 1 1
2
3 4
5
2 2 3
2 2 2
1 2 -1
1 1 -2
2 2 1
2
-5 -3
2
2 2 -3
2 2 0
3
5 -1 5
1
3 3 3
4
0 4 4 5
1
2 2 0
2
-5 2
3
1 2 2
1 1 1
2 2 2
1
0
2
1 1 1
1 1 0
1
-3
4
1 1 2
1 1 3
1 1 -2
1 1 -1
2
-4 -3
2
2 2 -2
1 2 3
4
4 -4 5 0
1
1 4 -2
5
-5 -5 5 4 -4
4
3 5 1
3 5 1
4 4 -2
1 5 2
5
-3 2 5 5 -2
1
2 3 2
5
5 -3 5 -3 4
5
2 2 -3
3 3 -1
5 5 1
3 5 0
1 1 -2
4
-5 0 0 -1
4
4 4 -1
2 2 3
1 2 -2
2 3 -1
1
-1
4
1 1 2
1 1 -1
1 1 2
1 1 3
3
-1 -4 0
2
2 2 0
3 3 -3
3
3 4 1
4
3 3 -2
1 2 -3
2 3 -2
2 3 -1
5
1 -3 4 2 -1
2
1 2 -1
3 5 2
1
5
1
1 1 -1
5
-2 0 1 -3 5
1
2 4 3
4
2 -3 5 0
4
4 4 -1
1 1 1
4 4 -3
1 4 2
4
-2 1 5 -4
4
1 3 1
2 3 -1
1 3 -3
2 2 -3
3
0 2 -2
2
2 3 -1
3 3 0
4
-5 -4 -5 -4
3
2 3 1
1 3 -3
1 3 2
1
-2
1
1 1 -2
2
2 -1
2
2 2 1
2 2 -3
1
-1
1
1 1 -1
4
-2 2 5 4
5
3 3 0
1 4 1
2 3 -1
3 4 2
4 4 -1
3
-3 -3 -2
4
2 2 -1
2 3 2
1 1 -3
2 2 -2
4
4 0 3 -3
1
4 4 0
4
0 -5 -2 4
2
1 4 -2
1 4 1
3
-2 -4 1
1
3 3 2
3
-5 -2 4
3
2 3 2
3 3 -3
1 1 -1
4
5 1 1 5
2
3 4 0
1 4 2
3
5 2 -3
4
3 3 0
3 3 0
3 3 -1
2 2 -1
1
1
4
1 1 -3
1 1 -2
1 1 1
1 1 -2
4
-4 -4 -3 4
1
2 4 -2
1
5
5
1 1 -2
1 1 0
1 1 -2
1 1 -3
1 1 3
3
-1 -1 1
1
3 3 -1
4
-4 -5 5 -1
2
4 4 1
2 3 -3
4
-5 -5 -4 0
1
3 4 -3
1
-5
5
1 1 3
1 1 2
1 1 -3
1 1 -3
1 1 -2
1
2
4
1 1 0
1 1 1
1 1 3
1 1 1
2
-4 5
5
1 1 -3
2 2 1
2 2 -1
1 1 -2
2 2 2
3
0 3 0
1
1 1 -2
3
5 -2 5
2
1 1 3
2 3 2
2
-1 -1
2
1 1 2
2 2 -2
1
3
4
1 1 -1
1 1 1
1 1 3
1 1 -1
4
3 0 -5 5
3
2 3 -3
4 4 3
2 3 -3
3
-4 0 5
4
1 3 -1
1 3 3
3 3 -1
2 2 -2
4
-2 4 4 4
1
4 4 3
1
-1
1
1 1 -1
2
1 -2
5
1 1 2
1 2 1
2 2 3
2 2 -3
1 2 -1
2
2 -2
4
1 1 1
2 2 3
2 2 -2
2 2 0
2
-1 -4
2
1 2 2
1 1 -2
3
-3 1 -1
5
3 3 -1
3 3 0
2 2 -1
3 3 -1
2 3 0
5
3 -5 2 5 3
5
2 5 -1
4 4 -3
4 4 3
3 5 3
1 1 0
4
3 0 3 2
5
4 4 1
2 4 3
1 3 -2
3 4 -3
2 2 -2
5
2 5 -3 2 0
4
2 5 3
3 4 -1
5 5 3
1 3 1
3
1 3 5
5
1 2 -2
3 3 -2
1 2 -2
1 1 0
2 2 1
5
-5 -3 5 2 4
2
5 5 -3
3 4 -2
4
0 -5 0 -4
1
1 3 2
3
4 0 -3
1
3 3 0
4
5 -2 -4 1
3
4 4 2
4 4 -2
4 4 3`

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type op struct {
	l, r int
	x    int64
}

type SegTree struct {
	n    int
	min  []int64
	max  []int64
	add  []int64
	time []int
	cur  int
}

func NewSegTree(n int) *SegTree {
	size := 4 * (n + 2)
	return &SegTree{
		n:    n,
		min:  make([]int64, size),
		max:  make([]int64, size),
		add:  make([]int64, size),
		time: make([]int, size),
		cur:  1,
	}
}

func (st *SegTree) touch(i int) {
	if st.time[i] != st.cur {
		st.time[i] = st.cur
		st.min[i] = 0
		st.max[i] = 0
		st.add[i] = 0
	}
}

func (st *SegTree) push(i int) {
	st.touch(i)
	if st.add[i] != 0 {
		val := st.add[i]
		left, right := i<<1, i<<1|1
		st.touch(left)
		st.touch(right)
		st.add[left] += val
		st.min[left] += val
		st.max[left] += val
		st.add[right] += val
		st.min[right] += val
		st.max[right] += val
		st.add[i] = 0
	}
}

func (st *SegTree) pull(i int) {
	left, right := i<<1, i<<1|1
	if st.min[left] < st.min[right] {
		st.min[i] = st.min[left]
	} else {
		st.min[i] = st.min[right]
	}
	if st.max[left] > st.max[right] {
		st.max[i] = st.max[left]
	} else {
		st.max[i] = st.max[right]
	}
}

func (st *SegTree) rangeAdd(i, l, r, ql, qr int, val int64) {
	st.touch(i)
	if ql <= l && r <= qr {
		st.add[i] += val
		st.min[i] += val
		st.max[i] += val
		return
	}
	st.push(i)
	mid := (l + r) >> 1
	if ql <= mid {
		st.rangeAdd(i<<1, l, mid, ql, qr, val)
	}
	if qr > mid {
		st.rangeAdd(i<<1|1, mid+1, r, ql, qr, val)
	}
	st.pull(i)
}

func (st *SegTree) firstNonZero(i, l, r int) (int, int64, bool) {
	st.touch(i)
	if st.min[i] == 0 && st.max[i] == 0 {
		return 0, 0, false
	}
	if l == r {
		return l, st.min[i], true
	}
	st.push(i)
	mid := (l + r) >> 1
	if idx, val, ok := st.firstNonZero(i<<1, l, mid); ok {
		return idx, val, true
	}
	return st.firstNonZero(i<<1|1, mid+1, r)
}

func (st *SegTree) Clear() {
	st.cur++
	st.touch(1)
}

func expectedCase(n int, a []int64, ops []op) string {
	seg := NewSegTree(n)
	best := 0
	for j, operation := range ops {
		seg.rangeAdd(1, 1, n, operation.l, operation.r, operation.x)
		if _, val, ok := seg.firstNonZero(1, 1, n); ok {
			if val < 0 {
				best = j + 1
				seg.Clear()
			}
		}
	}
	diff := make([]int64, n+2)
	for i := 0; i < best; i++ {
		d := ops[i]
		diff[d.l] += d.x
		diff[d.r+1] -= d.x
	}
	cur := int64(0)
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		cur += diff[i]
		a[i] += cur
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(a[i], 10))
	}
	return sb.String()
}

func loadCases() ([]string, []string) {
	tokens := strings.Fields(testcasesRaw)
	if len(tokens) == 0 {
		fmt.Println("no embedded testcases")
		os.Exit(1)
	}
	t, err := strconv.Atoi(tokens[0])
	if err != nil {
		fmt.Println("invalid testcase count")
		os.Exit(1)
	}
	var inputs []string
	var expects []string
	pos := 1
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos >= len(tokens) {
			fmt.Printf("case %d incomplete\n", caseIdx+1)
			os.Exit(1)
		}
		n, errN := strconv.Atoi(tokens[pos])
		if errN != nil {
			fmt.Printf("case %d invalid n\n", caseIdx+1)
			os.Exit(1)
		}
		pos++
		a := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			if pos >= len(tokens) {
				fmt.Printf("case %d missing array values\n", caseIdx+1)
				os.Exit(1)
			}
			val, errV := strconv.ParseInt(tokens[pos], 10, 64)
			if errV != nil {
				fmt.Printf("case %d invalid array value\n", caseIdx+1)
				os.Exit(1)
			}
			a[i] = val
			pos++
		}
		if pos >= len(tokens) {
			fmt.Printf("case %d missing q\n", caseIdx+1)
			os.Exit(1)
		}
		q, errQ := strconv.Atoi(tokens[pos])
		if errQ != nil {
			fmt.Printf("case %d invalid q\n", caseIdx+1)
			os.Exit(1)
		}
		pos++
		ops := make([]op, q)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(strconv.Itoa(n))
		input.WriteByte('\n')
		for i := 1; i <= n; i++ {
			if i > 1 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatInt(a[i], 10))
		}
		input.WriteByte('\n')
		input.WriteString(strconv.Itoa(q))
		input.WriteByte('\n')
		for i := 0; i < q; i++ {
			if pos+2 >= len(tokens) {
				fmt.Printf("case %d missing ops\n", caseIdx+1)
				os.Exit(1)
			}
			l, _ := strconv.Atoi(tokens[pos])
			r, _ := strconv.Atoi(tokens[pos+1])
			x, _ := strconv.ParseInt(tokens[pos+2], 10, 64)
			pos += 3
			ops[i] = op{l, r, x}
			input.WriteString(fmt.Sprintf("%d %d %d\n", l, r, x))
		}
		inputs = append(inputs, input.String())
		expects = append(expects, expectedCase(n, a, ops))
	}
	return inputs, expects
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	inputs, expects := loadCases()
	for idx, input := range inputs {
		got, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expects[idx] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expects[idx], got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
