package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func expected(a, x, y int64) int64 {
	a2 := 2 * a
	x2 := 2 * x
	y2 := 2 * y
	if y2 <= 0 || y2%a2 == 0 {
		return -1
	}
	row := y2/a2 + 1
	var w int64
	if row == 1 || row%2 == 0 {
		w = 1
	} else {
		w = 2
	}
	var pos int64
	switch w {
	case 1:
		if x2 <= -a || x2 >= a {
			return -1
		}
		pos = 1
	case 2:
		if x2 <= -a2 || x2 >= a2 || x2 == 0 {
			return -1
		}
		if x2 < 0 {
			pos = 1
		} else {
			pos = 2
		}
	}
	L := row - 1
	var sumPrev int64
	switch {
	case L <= 0:
		sumPrev = 0
	case L == 1:
		sumPrev = 1
	default:
		countEven := L / 2
		countOdd := (L+1)/2 - 1
		sumPrev = 1 + countEven + countOdd*2
	}
	return sumPrev + pos
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("failed to open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanWords)
	tests := 0
	for {
		var a, x, y int64
		if !sc.Scan() {
			break
		}
		fmt.Sscan(sc.Text(), &a)
		if !sc.Scan() {
			fmt.Println("unexpected EOF")
			os.Exit(1)
		}
		fmt.Sscan(sc.Text(), &x)
		if !sc.Scan() {
			fmt.Println("unexpected EOF")
			os.Exit(1)
		}
		fmt.Sscan(sc.Text(), &y)
		expect := expected(a, x, y)
		input := fmt.Sprintf("%d %d %d\n", a, x, y)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.Output()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", tests+1, err)
			os.Exit(1)
		}
		var got int64
		fmt.Sscan(string(bytes.TrimSpace(out)), &got)
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d (a=%d x=%d y=%d)\n", tests+1, expect, got, a, x, y)
			os.Exit(1)
		}
		tests++
	}
	if err := sc.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", tests)
}
