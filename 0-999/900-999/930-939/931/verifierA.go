package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(a, b int) int {
	if a > b {
		a, b = b, a
	}
	diff := b - a
	x := diff / 2
	y := diff - x
	return x*(x+1)/2 + y*(y+1)/2
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not open testcasesA.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("empty test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Printf("test %d missing a\n", i+1)
			os.Exit(1)
		}
		var a int
		fmt.Sscan(scan.Text(), &a)
		if !scan.Scan() {
			fmt.Printf("test %d missing b\n", i+1)
			os.Exit(1)
		}
		var b int
		fmt.Sscan(scan.Text(), &b)
		exp := expected(a, b)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", a, b)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		var got int
		if _, err := fmt.Sscan(gotStr, &got); err != nil {
			fmt.Printf("case %d: cannot parse output: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed: expected %d got %d\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
