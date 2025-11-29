package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesRaw = `3 3 10 8 0 7 3 10 0 2 1 5 7 3 6 8 1
2 1 3 6 4 2 6 2 1 2
4 2 2 0 0 3 3 2 2 4 5 3 8 10 10 3 2 3
4 3 0 5 6 2 2 4 1 5 4 9 9 0 9 10 5 1 4 5
3 4 5 2 7 7 2 0 4 0 5 6 0 8 6 5 6 9 0`

type query struct {
	t    int
	x, y int64
	L, R int
}

type testCase struct {
	n int
	m int
	a []int64
	b []int64
	c []int64
	d []int64
	e []int64
}

// Embedded solver logic from 1139F.go.
func solve(tc testCase) []int {
	q := make([]query, 0, tc.n+tc.m)
	vals := make([]int64, 0, 2*(tc.n+tc.m))
	for i := 0; i < tc.n; i++ {
		x := tc.a[i] - tc.c[i]
		y := tc.a[i] + tc.c[i]
		q = append(q, query{t: i + 1, x: x, y: y, L: int(tc.a[i]), R: int(tc.b[i])})
		vals = append(vals, tc.a[i], tc.b[i])
	}
	for i := 0; i < tc.m; i++ {
		x := tc.d[i] - tc.e[i]
		y := tc.d[i] + tc.e[i]
		q = append(q, query{t: -(i + 1), x: x, y: y, L: int(tc.d[i]), R: int(tc.d[i])})
		vals = append(vals, tc.d[i], tc.d[i])
	}

	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	uni := vals[:0]
	for i, v := range vals {
		if i == 0 || v != vals[i-1] {
			uni = append(uni, v)
		}
	}
	for idx := range q {
		origL := q[idx].L
		origR := q[idx].R
		li := sort.Search(len(uni), func(i int) bool { return uni[i] >= int64(origL) })
		ri := sort.Search(len(uni), func(i int) bool { return uni[i] >= int64(origR) })
		q[idx].L = li + 1
		q[idx].R = ri + 1
	}

	sort.Slice(q, func(i, j int) bool {
		if q[i].x != q[j].x {
			return q[i].x < q[j].x
		}
		if q[i].y != q[j].y {
			return q[i].y < q[j].y
		}
		return q[i].t > q[j].t
	})

	tmp := make([]query, len(q))
	bit := make([]int, len(uni)+5)
	ans := make([]int, tc.m+1)

	var bitUpdate = func(i, v int) {
		for ; i < len(bit); i += i & -i {
			bit[i] += v
		}
	}
	var bitQuery = func(i int) int {
		s := 0
		for ; i > 0; i -= i & -i {
			s += bit[i]
		}
		return s
	}

	var cdq func(int, int)
	cdq = func(l, r int) {
		if l >= r {
			return
		}
		m := (l + r) >> 1
		cdq(l, m)
		cdq(m+1, r)
		i, j := l, m+1
		k := l
		for i <= m || j <= r {
			if j > r || (i <= m && (q[i].y < q[j].y || (q[i].y == q[j].y && q[i].t > 0))) {
				tmp[k] = q[i]
				if q[i].t > 0 {
					bitUpdate(q[i].L, 1)
					bitUpdate(q[i].R+1, -1)
				}
				i++
			} else {
				tmp[k] = q[j]
				if q[j].t < 0 {
					idx := -q[j].t
					ans[idx] += bitQuery(q[j].L)
				}
				j++
			}
			k++
		}
		for p := l; p <= m; p++ {
			if q[p].t > 0 {
				bitUpdate(q[p].L, -1)
				bitUpdate(q[p].R+1, 1)
			}
		}
		for p := l; p <= r; p++ {
			q[p] = tmp[p]
		}
	}

	cdq(0, len(q)-1)
	return ans[1:]
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range tc.c {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range tc.d {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range tc.e {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	input := sb.String()

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
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotFields := strings.Fields(out.String())
	if len(gotFields) != tc.m {
		return fmt.Errorf("expected %d outputs got %d", tc.m, len(gotFields))
	}
	want := solve(tc)
	for i := 0; i < tc.m; i++ {
		v, err := strconv.Atoi(gotFields[i])
		if err != nil {
			return fmt.Errorf("invalid integer output %q", gotFields[i])
		}
		if v != want[i] {
			return fmt.Errorf("at position %d expected %d got %d", i+1, want[i], v)
		}
	}
	return nil
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	tests := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("invalid test line: %q", line)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		if len(fields) != 2+3*n+2*m {
			return nil, fmt.Errorf("invalid counts for n=%d m=%d", n, m)
		}
		tc := testCase{n: n, m: m}
		tc.a = make([]int64, n)
		tc.b = make([]int64, n)
		tc.c = make([]int64, n)
		tc.d = make([]int64, m)
		tc.e = make([]int64, m)
		offset := 2
		for i := 0; i < n; i++ {
			tc.a[i], err = strconv.ParseInt(fields[offset+i], 10, 64)
			if err != nil {
				return nil, err
			}
		}
		offset += n
		for i := 0; i < n; i++ {
			tc.b[i], err = strconv.ParseInt(fields[offset+i], 10, 64)
			if err != nil {
				return nil, err
			}
		}
		offset += n
		for i := 0; i < n; i++ {
			tc.c[i], err = strconv.ParseInt(fields[offset+i], 10, 64)
			if err != nil {
				return nil, err
			}
		}
		offset += n
		for i := 0; i < m; i++ {
			tc.d[i], err = strconv.ParseInt(fields[offset+i], 10, 64)
			if err != nil {
				return nil, err
			}
		}
		offset += m
		for i := 0; i < m; i++ {
			tc.e[i], err = strconv.ParseInt(fields[offset+i], 10, 64)
			if err != nil {
				return nil, err
			}
		}
		tests = append(tests, tc)
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
