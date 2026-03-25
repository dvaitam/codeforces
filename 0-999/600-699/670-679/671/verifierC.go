package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseC struct {
	n   int
	arr []int
}

func genTestsC() []testCaseC {
	rand.Seed(3)
	tests := make([]testCaseC, 100)
	for i := range tests {
		n := rand.Intn(6) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(10) + 1
		}
		tests[i] = testCaseC{n, arr}
	}
	return tests
}

// ---------- embedded correct solver (from cf_t25_671_C.go) ----------

type stNode struct {
	sum  int64
	max  int
	min  int
	lazy int
}

var stTree []stNode

func stMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func stMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func stApplySet(node, l, r, X int) {
	stTree[node].sum = int64(X) * int64(r-l+1)
	stTree[node].max = X
	stTree[node].min = X
	stTree[node].lazy = X
}

func stPushDown(node, l, r int) {
	if stTree[node].lazy != -1 {
		mid := (l + r) / 2
		stApplySet(node*2, l, mid, stTree[node].lazy)
		stApplySet(node*2+1, mid+1, r, stTree[node].lazy)
		stTree[node].lazy = -1
	}
}

func stPullUp(node int) {
	stTree[node].sum = stTree[node*2].sum + stTree[node*2+1].sum
	stTree[node].max = stMax(stTree[node*2].max, stTree[node*2+1].max)
	stTree[node].min = stMin(stTree[node*2].min, stTree[node*2+1].min)
}

func stUpdateMin(node, l, r, ql, qr, X int) {
	if ql > qr {
		return
	}
	if l > qr || r < ql || stTree[node].max <= X {
		return
	}
	if ql <= l && r <= qr && stTree[node].min >= X {
		stApplySet(node, l, r, X)
		return
	}
	stPushDown(node, l, r)
	mid := (l + r) / 2
	stUpdateMin(node*2, l, mid, ql, qr, X)
	stUpdateMin(node*2+1, mid+1, r, ql, qr, X)
	stPullUp(node)
}

func stBuild(node, l, r int) {
	stTree[node].lazy = -1
	if l == r {
		stTree[node].sum = int64(l)
		stTree[node].max = l
		stTree[node].min = l
		return
	}
	mid := (l + r) / 2
	stBuild(node*2, l, mid)
	stBuild(node*2+1, mid+1, r)
	stPullUp(node)
}

func solveC(tc testCaseC) int64 {
	n := tc.n
	arr := tc.arr
	maxA := 0
	for _, v := range arr {
		if v > maxA {
			maxA = v
		}
	}

	pos := make([]int, maxA+1)
	for i := 0; i < n; i++ {
		pos[arr[i]] = i + 1
	}

	stTree = make([]stNode, 4*n+1)
	if n > 0 {
		stBuild(1, 1, n)
	}

	var totalAns int64

	for g := maxA; g >= 1; g-- {
		p1, p2 := n+1, n+1
		pk_1, pk := 0, 0

		for m := 1; m*g <= maxA; m++ {
			p := pos[m*g]
			if p != 0 {
				if p < p1 {
					p2 = p1
					p1 = p
				} else if p < p2 {
					p2 = p
				}
				if p > pk {
					pk_1 = pk
					pk = p
				} else if p > pk_1 {
					pk_1 = p
				}
			}
		}

		if pk_1 != 0 {
			stUpdateMin(1, 1, n, 1, pk_1-1, 0)
			stUpdateMin(1, 1, n, pk_1, pk-1, p1)
			stUpdateMin(1, 1, n, pk, n, p2)
		}

		if n > 0 {
			U := stTree[1].sum
			totalAns += int64(n)*int64(n+1)/2 - U
		}
	}

	return totalAns
}

// ---------- end embedded solver ----------

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()
	for i, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n", tc.n)
		for j, v := range tc.arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		expect := solveC(tc)
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output %q\n", i+1, out)
			os.Exit(1)
		}
		if val != expect {
			fmt.Fprintf(os.Stderr, "test %d: expected %d got %d\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
