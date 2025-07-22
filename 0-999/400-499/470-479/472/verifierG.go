package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
)

func countAnd(aBlk, bBlk []uint64, p1, p2, length int) int {
	const W = 64
	wa := p1 / W
	ba := p1 & (W - 1)
	wb := p2 / W
	bb := p2 & (W - 1)
	words := (ba + length + W - 1) / W
	var cnt int
	for i := 0; i < words; i++ {
		aw := aBlk[wa+i] >> uint(ba)
		if ba != 0 {
			aw |= aBlk[wa+i+1] << uint(W-ba)
		}
		bw := bBlk[wb+i] >> uint(bb)
		if bb != 0 {
			bw |= bBlk[wb+i+1] << uint(W-bb)
		}
		if i == words-1 {
			r := (ba + length) & (W - 1)
			if r != 0 {
				mask := uint64((1 << uint(r)) - 1)
				aw &= mask
				bw &= mask
			}
		}
		cnt += bits.OnesCount64(aw & bw)
	}
	return cnt
}

func solveCase(a, b string, queries [][3]int) []int {
	na, nb := len(a), len(b)
	preA := make([]int, na+1)
	preB := make([]int, nb+1)
	for i := 0; i < na; i++ {
		preA[i+1] = preA[i]
		if a[i] == '1' {
			preA[i+1]++
		}
	}
	for i := 0; i < nb; i++ {
		preB[i+1] = preB[i]
		if b[i] == '1' {
			preB[i+1]++
		}
	}
	const W = 64
	naBlk := (na + W - 1) / W
	nbBlk := (nb + W - 1) / W
	aBlk := make([]uint64, naBlk+1)
	bBlk := make([]uint64, nbBlk+1)
	for i := 0; i < na; i++ {
		if a[i] == '1' {
			aBlk[i/W] |= 1 << uint(i&(W-1))
		}
	}
	for i := 0; i < nb; i++ {
		if b[i] == '1' {
			bBlk[i/W] |= 1 << uint(i&(W-1))
		}
	}
	ans := make([]int, len(queries))
	for idx, q := range queries {
		p1, p2, l := q[0], q[1], q[2]
		sa1 := preA[p1+l] - preA[p1]
		sb1 := preB[p2+l] - preB[p2]
		inter := countAnd(aBlk, bBlk, p1, p2, l)
		ans[idx] = sa1 + sb1 - 2*inter
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesG.txt")
	if err != nil {
		fmt.Println("could not read testcasesG.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t := 0
	fmt.Sscan(scan.Text(), &t)
	cases := make([]struct {
		a, b    string
		queries [][3]int
	}, t)
	for ci := 0; ci < t; ci++ {
		scan.Scan()
		a := scan.Text()
		scan.Scan()
		b := scan.Text()
		scan.Scan()
		var q int
		fmt.Sscan(scan.Text(), &q)
		qs := make([][3]int, q)
		for i := 0; i < q; i++ {
			scan.Scan()
			fmt.Sscan(scan.Text(), &qs[i][0], &qs[i][1], &qs[i][2])
		}
		cases[ci] = struct {
			a, b    string
			queries [][3]int
		}{a, b, qs}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for idx, c := range cases {
		expected := solveCase(c.a, c.b, c.queries)
		for i := 0; i < len(c.queries); i++ {
			if !outScan.Scan() {
				fmt.Printf("missing output for case %d query %d\n", idx+1, i+1)
				os.Exit(1)
			}
			var got int
			fmt.Sscan(outScan.Text(), &got)
			if got != expected[i] {
				fmt.Printf("case %d query %d failed: expected %d got %d\n", idx+1, i+1, expected[i], got)
				os.Exit(1)
			}
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
