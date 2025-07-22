package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Println("could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	type caseData struct {
		n    int
		x, y []uint32
	}
	cases := make([]caseData, t)
	for idx := 0; idx < t; idx++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		x := make([]uint32, n)
		y := make([]uint32, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			v, _ := strconv.ParseUint(scan.Text(), 10, 32)
			x[i] = uint32(v)
		}
		for i := 0; i < n; i++ {
			scan.Scan()
			v, _ := strconv.ParseUint(scan.Text(), 10, 32)
			y[i] = uint32(v)
		}
		cases[idx] = caseData{n, x, y}
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
		if !outScan.Scan() {
			fmt.Printf("missing output for case %d\n", idx+1)
			os.Exit(1)
		}
		tok := outScan.Text()
		if tok == "-1" {
			fmt.Printf("case %d: reported impossible but solution exists\n", idx+1)
			os.Exit(1)
		}
		m, err := strconv.Atoi(tok)
		if err != nil {
			fmt.Printf("case %d: bad m\n", idx+1)
			os.Exit(1)
		}
		if m < 0 || m > 1000000 {
			fmt.Printf("case %d: invalid m %d\n", idx+1, m)
			os.Exit(1)
		}
		arr := make([]uint32, len(c.x))
		copy(arr, c.x)
		for op := 0; op < m; op++ {
			if !outScan.Scan() {
				fmt.Printf("case %d: missing op %d\n", idx+1, op+1)
				os.Exit(1)
			}
			i, _ := strconv.Atoi(outScan.Text())
			if !outScan.Scan() {
				fmt.Printf("case %d: missing op %d\n", idx+1, op+1)
				os.Exit(1)
			}
			j, _ := strconv.Atoi(outScan.Text())
			if i < 1 || i > c.n || j < 1 || j > c.n {
				fmt.Printf("case %d: invalid indices\n", idx+1)
				os.Exit(1)
			}
			arr[i-1] ^= arr[j-1]
		}
		for i := 0; i < c.n; i++ {
			if arr[i] != c.y[i] {
				fmt.Printf("case %d failed: array mismatch\n", idx+1)
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
