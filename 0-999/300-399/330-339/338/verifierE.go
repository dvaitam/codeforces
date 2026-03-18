package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, input := range tests {
		want, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if normalize(got) != normalize(want) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nInput:\n%sExpected:\n%sGot:\n%s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

const embeddedRefSource = `package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

type Scanner struct {
	b   []byte
	pos int
}

func NewScanner() *Scanner {
	b, _ := io.ReadAll(os.Stdin)
	return &Scanner{b: b, pos: 0}
}

func (s *Scanner) nextInt() int {
	for s.pos < len(s.b) && (s.b[s.pos] < '0' || s.b[s.pos] > '9') {
		s.pos++
	}
	if s.pos >= len(s.b) {
		return 0
	}
	res := 0
	for s.pos < len(s.b) && s.b[s.pos] >= '0' && s.b[s.pos] <= '9' {
		res = res*10 + int(s.b[s.pos]-'0')
		s.pos++
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var tree []int
var lazy []int

func build(node, l, r int) {
	if l == r {
		tree[node] = -l
		return
	}
	mid := (l + r) / 2
	build(node*2, l, mid)
	build(node*2+1, mid+1, r)
	tree[node] = min(tree[node*2], tree[node*2+1])
}

func push(node int) {
	if lazy[node] != 0 {
		tree[node*2] += lazy[node]
		lazy[node*2] += lazy[node]
		tree[node*2+1] += lazy[node]
		lazy[node*2+1] += lazy[node]
		lazy[node] = 0
	}
}

func update(node, l, r, ql, qr, val int) {
	if ql <= l && r <= qr {
		tree[node] += val
		lazy[node] += val
		return
	}
	push(node)
	mid := (l + r) / 2
	if ql <= mid {
		update(node*2, l, mid, ql, qr, val)
	}
	if qr > mid {
		update(node*2+1, mid+1, r, ql, qr, val)
	}
	tree[node] = min(tree[node*2], tree[node*2+1])
}

func getIdx(x int, c []int, m int) int {
	left, right := 1, m
	ans := m + 1
	for left <= right {
		mid := (left + right) / 2
		if c[mid] <= x {
			ans = mid
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return ans
}

func main() {
	sc := NewScanner()
	n := sc.nextInt()
	if n == 0 {
		return
	}
	m := sc.nextInt()
	h := sc.nextInt()

	b := make([]int, m)
	for i := 0; i < m; i++ {
		b[i] = sc.nextInt()
	}

	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = sc.nextInt()
	}

	sort.Ints(b)
	c := make([]int, m+1)
	for i := 1; i <= m; i++ {
		c[i] = h - b[i-1]
	}

	tree = make([]int, 4*m+5)
	lazy = make([]int, 4*m+5)

	build(1, 1, m)

	ans := 0
	for i := 0; i < m; i++ {
		idx := getIdx(a[i], c, m)
		if idx <= m {
			update(1, 1, m, idx, m, 1)
		}
	}

	if tree[1] >= 0 {
		ans++
	}

	for i := m; i < n; i++ {
		idxOut := getIdx(a[i-m], c, m)
		if idxOut <= m {
			update(1, 1, m, idxOut, m, -1)
		}

		idxIn := getIdx(a[i], c, m)
		if idxIn <= m {
			update(1, 1, m, idxIn, m, 1)
		}

		if tree[1] >= 0 {
			ans++
		}
	}

	fmt.Println(ans)
}
`

func buildReference() (string, error) {
	outPath := fmt.Sprintf("%s/ref338E_%d.bin", os.TempDir(), time.Now().UnixNano())
	tmpGo := outPath + ".go"
	if err := os.WriteFile(tmpGo, []byte(embeddedRefSource), 0644); err != nil {
		return "", err
	}
	cmd := exec.Command("go", "build", "-o", outPath, tmpGo)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmpGo)
		return "", fmt.Errorf("Go build failed: %v\n%s", err, string(out))
	}
	os.Remove(tmpGo)
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func normalize(out string) string {
	return strings.TrimSpace(out)
}

func generateTests() []string {
	var tests []string
	tests = append(tests, formatTest(1, 1, 0, []int64{3}, []int64{3}))
	tests = append(tests, formatTest(2, 1, 5, []int64{10}, []int64{15, 20}))
	tests = append(tests, formatTest(5, 3, 2, []int64{1, 3, 5}, []int64{1, 2, 3, 4, 5}))
	tests = append(tests, formatTest(3, 3, 10, []int64{4, 7, 11}, []int64{5, 8, 12}))
	tests = append(tests, formatTest(4, 2, 0, []int64{1000000000, 1}, []int64{1000000000, 2, 1000000000, 3}))
	tests = append(tests, formatIncreasing(2000, 500))
	tests = append(tests, formatConstant(1500, 800))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 160 {
		var n int
		if len(tests)%20 == 0 {
			n = 50000 + rng.Intn(20000)
		} else {
			n = rng.Intn(80) + 1
		}
		m := rng.Intn(n) + 1
		h := rng.Int63n(1_000_000_000) + 1
		b := make([]int64, m)
		a := make([]int64, n)
		for i := 0; i < m; i++ {
			b[i] = rng.Int63n(1_000_000_000) + 1
		}
		for i := 0; i < n; i++ {
			val := rng.Int63n(1_000_000_000) + 1
			if rng.Intn(6) == 0 && i > 0 {
				diff := int64(rng.Intn(2000)) - 1000
				val = a[i-1] + diff
				if val < 1 {
					val = 1
				}
				if val > 1_000_000_000 {
					val = 1_000_000_000
				}
			}
			a[i] = val
		}
		tests = append(tests, formatTest(n, m, h, b, a))
	}
	return tests
}

func formatTest(n, m int, h int64, b, a []int64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, h))
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func formatIncreasing(n, m int) string {
	if m > n {
		m = n
	}
	h := int64(5)
	b := make([]int64, m)
	a := make([]int64, n)
	for i := 0; i < m; i++ {
		b[i] = int64(i * 3)
	}
	for i := 0; i < n; i++ {
		a[i] = int64(i * 3)
	}
	return formatTest(n, m, h, b, a)
}

func formatConstant(n, m int) string {
	if m > n {
		m = n
	}
	h := int64(0)
	b := make([]int64, m)
	a := make([]int64, n)
	for i := 0; i < m; i++ {
		b[i] = 1
	}
	for i := 0; i < n; i++ {
		a[i] = 1
	}
	return formatTest(n, m, h, b, a)
}
