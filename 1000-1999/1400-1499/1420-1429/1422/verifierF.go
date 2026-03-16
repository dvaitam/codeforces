package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MOD = 1000000007
const MAXA = 200005

var spfG [MAXA]int32
var invG [MAXA]int32

type Node struct {
	ls, rs int32
	val    int32
}

var treeG []Node

type upd struct {
	pos int32
	val int32
}

func updateTreeBatch(node int32, l, r int32, upds []upd) int32 {
	newNode := int32(len(treeG))
	treeG = append(treeG, treeG[node])

	mul := int64(1)
	for i := 0; i < len(upds); i++ {
		mul = (mul * int64(upds[i].val)) % MOD
	}
	treeG[newNode].val = int32((int64(treeG[newNode].val) * mul) % MOD)

	if l == r {
		return newNode
	}

	mid := (l + r) >> 1
	split := 0
	for split < len(upds) && upds[split].pos <= mid {
		split++
	}

	if split > 0 {
		treeG[newNode].ls = updateTreeBatch(treeG[node].ls, l, mid, upds[:split])
	}
	if split < len(upds) {
		treeG[newNode].rs = updateTreeBatch(treeG[node].rs, mid+1, r, upds[split:])
	}

	return newNode
}

func queryTree(node int32, l, r int32, ql, qr int32) int32 {
	if node == 0 {
		return 1
	}
	if ql <= l && r <= qr {
		return treeG[node].val
	}
	mid := (l + r) >> 1
	res := int32(1)
	if ql <= mid {
		res = int32((int64(res) * int64(queryTree(treeG[node].ls, l, mid, ql, qr))) % MOD)
	}
	if qr > mid {
		res = int32((int64(res) * int64(queryTree(treeG[node].rs, mid+1, r, ql, qr))) % MOD)
	}
	return res
}

func solve(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 1024*1024*10)

	scanInt := func() int32 {
		scanner.Scan()
		b := scanner.Bytes()
		res := int32(0)
		for _, v := range b {
			res = res*10 + int32(v-'0')
		}
		return res
	}

	n := scanInt()

	a := make([]int32, n+1)
	for i := int32(1); i <= n; i++ {
		a[i] = scanInt()
	}

	treeG = make([]Node, 1, 5000000)
	treeG[0] = Node{0, 0, 1}

	roots := make([]int32, n+1)
	roots[0] = 0

	prev := make([]int32, MAXA)

	for i := int32(1); i <= n; i++ {
		x := a[i]
		var upds []upd

		for x > 1 {
			p := spfG[x]
			pk := int32(1)
			for x%p == 0 {
				x /= p
				pk *= p

				upds = append(upds, upd{pos: i, val: p})
				if prev[pk] != 0 {
					upds = append(upds, upd{pos: prev[pk], val: invG[p]})
				}
				prev[pk] = i
			}
		}

		combined := make([]upd, 0, len(upds))
		for j := 0; j < len(upds); j++ {
			u := upds[j]
			found := false
			for k := 0; k < len(combined); k++ {
				if combined[k].pos == u.pos {
					combined[k].val = int32((int64(combined[k].val) * int64(u.val)) % MOD)
					found = true
					break
				}
			}
			if !found {
				combined = append(combined, u)
			}
		}

		for j := 1; j < len(combined); j++ {
			k := j
			for k > 0 && combined[k-1].pos > combined[k].pos {
				combined[k-1], combined[k] = combined[k], combined[k-1]
				k--
			}
		}

		currRoot := roots[i-1]
		if len(combined) > 0 {
			currRoot = updateTreeBatch(currRoot, 1, n, combined)
		}
		roots[i] = currRoot
	}

	// Reset prev for next call
	for i := range prev {
		prev[i] = 0
	}

	q := scanInt()
	last := int32(0)

	var out strings.Builder
	for i := int32(0); i < q; i++ {
		x := scanInt()
		y := scanInt()

		l := int32((int64(last)+int64(x))%int64(n)) + 1
		r := int32((int64(last)+int64(y))%int64(n)) + 1
		if l > r {
			l, r = r, l
		}

		ans := queryTree(roots[r], 1, n, l, r)
		out.WriteString(fmt.Sprint(ans))
		if i+1 < q {
			out.WriteByte('\n')
		}
		last = ans
	}
	return out.String()
}

type test struct {
	input    string
	expected string
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(47))
	tests := []test{}
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		q := rng.Intn(5) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(10)+1)
		}
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "%d\n", q)
		for i := 0; i < q; i++ {
			x := rng.Intn(n)
			y := rng.Intn(n)
			fmt.Fprintf(&sb, "%d %d\n", x, y)
		}
		inp := sb.String()

		// Reset global state for solve
		treeG = nil
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	// Initialize spf and inv tables
	for i := int32(2); i < MAXA; i++ {
		if spfG[i] == 0 {
			for j := i; j < MAXA; j += i {
				if spfG[j] == 0 {
					spfG[j] = i
				}
			}
		}
	}
	invG[1] = 1
	for i := int32(2); i < MAXA; i++ {
		invG[i] = int32((int64(MOD) - int64(MOD/i)) * int64(invG[MOD%i]) % MOD)
	}

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotFields := strings.Fields(got)
		expFields := strings.Fields(t.expected)
		match := len(gotFields) == len(expFields)
		if match {
			for k := range gotFields {
				if gotFields[k] != expFields[k] {
					match = false
					break
				}
			}
		}
		if !match {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
