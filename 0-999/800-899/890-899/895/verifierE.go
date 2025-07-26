package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Node struct {
	x, t, s float64
	sz      int
}

var (
	tree []Node
	D    []float64
	N    int
)

func (nd *Node) upd(tt, ts float64) {
	nd.x = tt*nd.x + ts*float64(nd.sz)
	nd.t *= tt
	nd.s = nd.s*tt + ts
}

func build(i, l, r int) {
	tree[i].t = 1
	tree[i].s = 0
	tree[i].sz = r - l
	if r-l == 1 {
		tree[i].x = D[l]
		return
	}
	m := (l + r) >> 1
	build(i<<1, l, m)
	build(i<<1|1, m, r)
	tree[i].x = tree[i<<1].x + tree[i<<1|1].x
}

func push(i int) {
	tt := tree[i].t
	ts := tree[i].s
	if tt != 1 || ts != 0 {
		tree[i<<1].upd(tt, ts)
		tree[i<<1|1].upd(tt, ts)
		tree[i].t = 1
		tree[i].s = 0
	}
}

func query(i, l, r, ql, qr int) float64 {
	if ql <= l && r <= qr {
		return tree[i].x
	}
	push(i)
	m := (l + r) >> 1
	var ans float64
	if ql < m && qr > l {
		ans += query(i<<1, l, m, ql, qr)
	}
	if qr > m && ql < r {
		ans += query(i<<1|1, m, r, ql, qr)
	}
	return ans
}

func modify(i, l, r, ql, qr int, tt, ts float64) {
	if ql <= l && r <= qr {
		tree[i].upd(tt, ts)
		return
	}
	push(i)
	m := (l + r) >> 1
	if ql < m && qr > l {
		modify(i<<1, l, m, ql, qr, tt, ts)
	}
	if qr > m && ql < r {
		modify(i<<1|1, m, r, ql, qr, tt, ts)
	}
	tree[i].x = tree[i<<1].x + tree[i<<1|1].x
}

func solveCase(n, q int, arr []int, ops []int, args [][]int) []string {
	N = n
	D = make([]float64, n)
	for i := 0; i < n; i++ {
		D[i] = float64(arr[i])
	}
	tree = make([]Node, 4*n+5)
	build(1, 0, n)
	var res []string
	for i := 0; i < q; i++ {
		if ops[i] == 1 {
			l1, r1, l2, r2 := args[i][0]-1, args[i][1], args[i][2]-1, args[i][3]
			b0 := r1
			d0 := r2
			len1 := b0 - l1
			len2 := d0 - l2
			sum2 := query(1, 0, n, l2, d0)
			sum1 := query(1, 0, n, l1, b0)
			s1 := sum2 / float64(len2) / float64(len1)
			s2 := sum1 / float64(len2) / float64(len1)
			t1 := 1 - 1/float64(len1)
			t2 := 1 - 1/float64(len2)
			modify(1, 0, n, l1, b0, t1, s1)
			modify(1, 0, n, l2, d0, t2, s2)
		} else {
			l, r := args[i][0]-1, args[i][1]
			sum := query(1, 0, n, l, r)
			res = append(res, fmt.Sprintf("%.7f", sum))
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		pos := 0
		n, _ := strconv.Atoi(fields[pos])
		pos++
		q, _ := strconv.Atoi(fields[pos])
		pos++
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[pos])
			pos++
			arr[i] = v
		}
		ops := make([]int, q)
		args := make([][]int, q)
		for i := 0; i < q; i++ {
			op, _ := strconv.Atoi(fields[pos])
			pos++
			ops[i] = op
			if op == 1 {
				l1, _ := strconv.Atoi(fields[pos])
				pos++
				r1, _ := strconv.Atoi(fields[pos])
				pos++
				l2, _ := strconv.Atoi(fields[pos])
				pos++
				r2, _ := strconv.Atoi(fields[pos])
				pos++
				args[i] = []int{l1, r1, l2, r2}
			} else {
				l, _ := strconv.Atoi(fields[pos])
				pos++
				r, _ := strconv.Atoi(fields[pos])
				pos++
				args[i] = []int{l, r}
			}
		}
		if pos != len(fields) {
			fmt.Println("bad test case length")
			os.Exit(1)
		}
		expected := solveCase(n, q, arr, ops, args)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		for i := 0; i < q; i++ {
			if ops[i] == 1 {
				fmt.Fprintf(&input, "1 %d %d %d %d\n", args[i][0], args[i][1], args[i][2], args[i][3])
			} else {
				fmt.Fprintf(&input, "2 %d %d\n", args[i][0], args[i][1])
			}
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		outScan := bufio.NewScanner(bytes.NewReader(out.Bytes()))
		for i := 0; i < len(expected); i++ {
			if !outScan.Scan() {
				fmt.Printf("test %d: missing output line %d\n", idx, i+1)
				os.Exit(1)
			}
			got := strings.TrimSpace(outScan.Text())
			if got != expected[i] {
				fmt.Printf("test %d failed on query %d: expected %s got %s\n", idx, i+1, expected[i], got)
				os.Exit(1)
			}
		}
		if outScan.Scan() {
			fmt.Printf("test %d: extra output detected\n", idx)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
