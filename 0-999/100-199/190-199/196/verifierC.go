package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type point struct{ x, y int }

type testCase struct {
	n     int
	edges [][2]int
	pts   []point
}

func orientation(a, b, c point) int64 {
	return int64(b.x-a.x)*int64(c.y-a.y) - int64(b.y-a.y)*int64(c.x-a.x)
}

func segmentsCross(a1, a2, b1, b2 point) bool {
	o1 := orientation(a1, a2, b1)
	o2 := orientation(a1, a2, b2)
	o3 := orientation(b1, b2, a1)
	o4 := orientation(b1, b2, a2)
	if o1 == 0 && o2 == 0 {
		// Collinear - overlapping check
		minax := a1.x
		maxax := a2.x
		if minax > maxax {
			minax, maxax = maxax, minax
		}
		minbx := b1.x
		maxbx := b2.x
		if minbx > maxbx {
			minbx, maxbx = maxbx, minbx
		}
		minay := a1.y
		maxay := a2.y
		if minay > maxay {
			minay, maxay = maxay, minay
		}
		minby := b1.y
		maxby := b2.y
		if minby > maxby {
			minby, maxby = maxby, minby
		}
		if maxax < minbx || maxbx < minax || maxay < minby || maxby < minay {
			return false
		}
		return true
	}
	return (o1 > 0) != (o2 > 0) && (o3 > 0) != (o4 > 0)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscanf(scan.Text(), "%d", &t)
	cases := make([]testCase, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		fmt.Sscanf(scan.Text(), "%d", &cases[i].n)
		cases[i].edges = make([][2]int, cases[i].n-1)
		for j := 0; j < cases[i].n-1; j++ {
			scan.Scan()
			fmt.Sscanf(scan.Text(), "%d", &cases[i].edges[j][0])
			scan.Scan()
			fmt.Sscanf(scan.Text(), "%d", &cases[i].edges[j][1])
		}
		cases[i].pts = make([]point, cases[i].n)
		for j := 0; j < cases[i].n; j++ {
			scan.Scan()
			fmt.Sscanf(scan.Text(), "%d", &cases[i].pts[j].x)
			scan.Scan()
			fmt.Sscanf(scan.Text(), "%d", &cases[i].pts[j].y)
		}
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for idx, tc := range cases {
		n := tc.n
		perm := make([]int, n)
		for i := 0; i < n; i++ {
			if !outScan.Scan() {
				fmt.Printf("missing output for test %d\n", idx+1)
				os.Exit(1)
			}
			fmt.Sscanf(outScan.Text(), "%d", &perm[i])
			if perm[i] < 1 || perm[i] > n {
				fmt.Printf("test %d failed: value out of range\n", idx+1)
				os.Exit(1)
			}
		}
		seen := make([]bool, n+1)
		for _, v := range perm {
			seen[v] = !seen[v]
			if !seen[v] {
				fmt.Printf("test %d failed: duplicate vertices\n", idx+1)
				os.Exit(1)
			}
		}
		// Build mapping from vertex to point coordinates
		vertPt := make([]point, n+1)
		for i, v := range perm {
			vertPt[v] = tc.pts[i]
		}
		// check edges for crossing
		for i := 0; i < len(tc.edges); i++ {
			for j := i + 1; j < len(tc.edges); j++ {
				u1, v1 := tc.edges[i][0], tc.edges[i][1]
				u2, v2 := tc.edges[j][0], tc.edges[j][1]
				if u1 == u2 || u1 == v2 || v1 == u2 || v1 == v2 {
					continue
				}
				if segmentsCross(vertPt[u1], vertPt[v1], vertPt[u2], vertPt[v2]) {
					fmt.Printf("test %d failed: edges cross\n", idx+1)
					os.Exit(1)
				}
			}
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
