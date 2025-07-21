package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func heightIntersect(h0, r0, R0, h1, r1, R1 float64) float64 {
	res1 := (r1 - r0) / (R0 - r0) * h0
	res2 := (R0 - r1) / (R1 - r1) * h1
	res3 := (R1 - r0) / (R0 - r0) * h0
	if res2 > h1 {
		res2 = -1000
	} else {
		res2 = h0 - res2
	}
	if res3 > h0 {
		res3 = h0
	}
	if res3 < 0 {
		res3 = 0
	}
	res3 -= h1
	if r1 >= R0 {
		return h0
	}
	if res1 < 0 {
		res1 = 0
	}
	if res2 > h0 {
		res2 = h0
	}
	if res1 > res2 && res1 > res3 {
		return res1
	}
	if res3 > res2 {
		return res3
	}
	return res2
}

func solveCase(h, r, R []float64) float64 {
	n := len(h)
	sh := make([]float64, n)
	ch := make([]float64, n)
	cr := make([]float64, n)
	cR := make([]float64, n)
	curr := 1
	sh[0] = 0
	ch[0] = h[0]
	cr[0] = r[0]
	cR[0] = R[0]
	for i := 1; i < n; i++ {
		mh := 0.0
		for j := 0; j < curr; j++ {
			hh := heightIntersect(ch[j], cr[j], cR[j], h[i], r[i], R[i]) + sh[j]
			if hh > mh {
				mh = hh
			}
		}
		prog := 0
		for j := 0; j < curr; j++ {
			if sh[j]+ch[j] > mh {
				sh[prog] = sh[j]
				ch[prog] = ch[j]
				cr[prog] = cr[j]
				cR[prog] = cR[j]
				prog++
			}
		}
		sh[prog] = mh
		ch[prog] = h[i]
		cr[prog] = r[i]
		cR[prog] = R[i]
		curr = prog + 1
	}
	res := 0.0
	for i := 0; i < curr; i++ {
		if sh[i]+ch[i] > res {
			res = sh[i] + ch[i]
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
		fmt.Println("bad file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		var n int
		fmt.Sscan(scan.Text(), &n)
		h := make([]float64, n)
		r := make([]float64, n)
		R := make([]float64, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &h[j])
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &r[j])
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &R[j])
		}
		expected := solveCase(h, r, R)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n", n)
		for j := 0; j < n; j++ {
			fmt.Fprintf(&input, "%g %g %g\n", h[j], r[j], R[j])
		}
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		var got float64
		fmt.Sscan(strings.TrimSpace(string(out)), &got)
		if diff := got - expected; diff < -1e-6 || diff > 1e-6 {
			fmt.Printf("case %d failed: expected %.6f got %.6f\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
