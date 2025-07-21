package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type dsuNext struct {
	parent []int
	n      int
}

func newDSUNext(n int) *dsuNext {
	p := make([]int, n+2)
	for i := 0; i <= n+1; i++ {
		p[i] = i
	}
	return &dsuNext{parent: p, n: n}
}

func (d *dsuNext) find(x int) int {
	if x > d.n+1 {
		return d.n + 1
	}
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsuNext) remove(x int) { d.parent[x] = d.find(x + 1) }

func solve(N int, x int64, M int, y int64) int64 {
	d := x - y
	if d < 0 {
		d = -d
	}
	Li := make([]int, N+1)
	Ri := make([]int, N+1)
	cntI := 0
	var P int64
	for i := 1; i <= N; i++ {
		dd := int64(i) - d
		if dd < 0 {
			dd = -dd
		}
		l := int(dd) + 1
		if l < 1 {
			l = 1
		}
		r := int64(i) + d - 1
		if r > int64(M) {
			r = int64(M)
		}
		if l <= int(r) {
			Li[i] = l
			Ri[i] = int(r)
			cntI++
			P += int64(r - int64(l) + 1)
		} else {
			Li[i] = 1
			Ri[i] = 0
		}
	}
	Lj := make([]int, M+1)
	Rj := make([]int, M+1)
	cntJ := 0
	for j := 1; j <= M; j++ {
		dd := int64(j) - d
		if dd < 0 {
			dd = -dd
		}
		l := int(dd) + 1
		if l < 1 {
			l = 1
		}
		r := int64(j) + d - 1
		if r > int64(N) {
			r = int64(N)
		}
		if l <= int(r) {
			Lj[j] = l
			Rj[j] = int(r)
			cntJ++
		} else {
			Lj[j] = 1
			Rj[j] = 0
		}
	}
	dsuI := newDSUNext(N)
	dsuJ := newDSUNext(M)
	for i := 1; i <= N; i++ {
		if Ri[i] < Li[i] {
			dsuI.remove(i)
		}
	}
	for j := 1; j <= M; j++ {
		if Rj[j] < Lj[j] {
			dsuJ.remove(j)
		}
	}
	CCbig := 0
	type node struct {
		left bool
		idx  int
	}
	queue := []node{}
	for i := 1; i <= N; i++ {
		ii := dsuI.find(i)
		if ii > N {
			break
		}
		CCbig++
		queue = append(queue, node{true, ii})
		dsuI.remove(ii)
		for q := 0; q < len(queue); q++ {
			nd := queue[q]
			if nd.left {
				i0 := nd.idx
				for j := dsuJ.find(Li[i0]); j <= Ri[i0]; j = dsuJ.find(j) {
					queue = append(queue, node{false, j})
					dsuJ.remove(j)
				}
			} else {
				j0 := nd.idx
				for i2 := dsuI.find(Lj[j0]); i2 <= Rj[j0]; i2 = dsuI.find(i2) {
					queue = append(queue, node{true, i2})
					dsuI.remove(i2)
				}
			}
		}
		queue = queue[:0]
	}
	Z := int64((N - cntI) + (M - cntJ))
	Cc := Z + int64(CCbig)
	F := 2*P + Cc + 1
	return F
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var N, M int
		var x, y int64
		fmt.Sscan(line, &N, &x, &M, &y)
		exp := solve(N, x, M, y)
		input := fmt.Sprintf("%d %d %d %d\n", N, x, M, y)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		var got int64
		fmt.Sscan(strings.TrimSpace(string(out)), &got)
		if got != exp {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, strings.TrimSpace(string(out)))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
