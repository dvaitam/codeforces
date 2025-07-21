package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type Point struct{ x, y int }

func solveCase(points []Point, d int) string {
	n := len(points)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				ab := abs(points[i].x-points[j].x) + abs(points[i].y-points[j].y)
				ac := abs(points[i].x-points[k].x) + abs(points[i].y-points[k].y)
				bc := abs(points[j].x-points[k].x) + abs(points[j].y-points[k].y)
				if ab == d && ac == d && bc == d {
					return fmt.Sprintf("%d %d %d", i+1, j+1, k+1)
				}
			}
		}
	}
	return "0 0 0"
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		d, _ := strconv.Atoi(scan.Text())
		pts := make([]Point, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			x, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			y, _ := strconv.Atoi(scan.Text())
			pts[j] = Point{x, y}
		}
		expected[i] = solveCase(pts, d)
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
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		a := outScan.Text()
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		b := outScan.Text()
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		c := outScan.Text()
		got := fmt.Sprintf("%s %s %s", a, b, c)
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
