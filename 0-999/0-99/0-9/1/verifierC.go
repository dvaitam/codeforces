package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
)

func gcd(x, y float64) float64 {
	for y > 1e-7 {
		x, y = y, math.Mod(x, y)
	}
	return x
}

func expectedArea(ax, ay, bx, by, cx, cy float64) float64 {
	a := math.Hypot(ax-bx, ay-by)
	b := math.Hypot(ax-cx, ay-cy)
	c := math.Hypot(bx-cx, by-cy)
	p := (a + b + c) / 2.0
	s := math.Sqrt(p * (p - a) * (p - b) * (p - c))
	R := a * b * c / (4.0 * s)
	A := math.Acos((b*b + c*c - a*a) / (2 * b * c))
	B := math.Acos((a*a + c*c - b*b) / (2 * a * c))
	C := math.Acos((a*a + b*b - c*c) / (2 * a * b))
	n := math.Pi / gcd(A, gcd(B, C))
	return R * R * math.Sin(2*math.Pi/n) * n / 2
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var ax, ay, bx, by, cx, cy float64
		fmt.Sscan(line, &ax, &ay, &bx, &by, &cx, &cy)
		exp := expectedArea(ax, ay, bx, by, cx, cy)
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("%f %f %f %f %f %f\n", ax, ay, bx, by, cx, cy))
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		var got float64
		outStr := strings.TrimSpace(outBuf.String())
		fmt.Sscan(outStr, &got)
		if math.Abs(got-exp) > 1e-4*math.Max(1.0, math.Abs(exp)) {
			fmt.Printf("Test %d failed: expected %f got %s\n", idx, exp, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
