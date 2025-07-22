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

var r, h float64
var pi = math.Pi

const eps = 1e-6

func isGreaterThan(a, b float64) bool {
	return a-b > eps
}

func isZero(a float64) bool {
	return math.Abs(a) < eps
}

func findCoordinatesOnPlane(x, y, z float64) (theta, d float64) {
	angle := math.Atan2(y, x)
	if angle < 0 {
		angle += 2 * pi
	}
	theta = angle * r / math.Hypot(r, h)
	if !isZero(z) {
		d = math.Sqrt(x*x + y*y + (h-z)*(h-z))
	} else {
		d = math.Hypot(r, h)
	}
	return
}

func euclideanDistance(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2
	return math.Hypot(dx, dy)
}

func findConeSurfaceDistance(x1, y1, z1, x2, y2, z2 float64) float64 {
	t1, d1 := findCoordinatesOnPlane(x1, y1, z1)
	t2, d2 := findCoordinatesOnPlane(x2, y2, z2)
	var aT, bT, aD, bD float64
	if t1 > t2 || (isZero(t1-t2) && d1 > d2) {
		aT, aD = t1, d1
		bT, bD = t2, d2
	} else {
		aT, aD = t2, d2
		bT, bD = t1, d1
	}
	alpha := 2 * pi * r / math.Hypot(r, h)
	if isGreaterThan(aT-bT, alpha/2) {
		aT = bT + (alpha - (aT - bT))
	}
	xA := aD * math.Cos(aT)
	yA := aD * math.Sin(aT)
	xB := bD * math.Cos(bT)
	yB := bD * math.Sin(bT)
	return euclideanDistance(xA, yA, xB, yB)
}

func sortaBinarySearch(left, right, xHi, yHi, zHi, xLo, yLo, zLo float64) float64 {
	x := r * math.Cos(left)
	y := r * math.Sin(left)
	leftVal := euclideanDistance(x, y, xLo, yLo) + findConeSurfaceDistance(xHi, yHi, zHi, x, y, zLo)
	x = r * math.Cos(right)
	y = r * math.Sin(right)
	rightVal := euclideanDistance(x, y, xLo, yLo) + findConeSurfaceDistance(xHi, yHi, zHi, x, y, zLo)
	for !isZero(right - left) {
		mid := (left + right) / 2
		x = r * math.Cos(mid)
		y = r * math.Sin(mid)
		midVal := euclideanDistance(x, y, xLo, yLo) + findConeSurfaceDistance(xHi, yHi, zHi, x, y, zLo)
		if leftVal > rightVal {
			left = mid
			leftVal = midVal
		} else {
			right = mid
			rightVal = midVal
		}
	}
	if leftVal < rightVal {
		return leftVal
	}
	return rightVal
}

func findOptimalSingleBrinkPointDistance(xHi, yHi, zHi, xLo, yLo, zLo float64) float64 {
	aDeg := pi / 180
	minVal := math.Inf(1)
	minDeg := 0.0
	for deg := 0.0; deg < 2*pi; deg += aDeg {
		x := r * math.Cos(deg)
		y := r * math.Sin(deg)
		curr := euclideanDistance(x, y, xLo, yLo) + findConeSurfaceDistance(xHi, yHi, zHi, x, y, zLo)
		if curr < minVal {
			minVal = curr
			minDeg = deg
		}
	}
	best := sortaBinarySearch(minDeg-aDeg, minDeg+aDeg, xHi, yHi, zHi, xLo, yLo, zLo)
	cand := sortaBinarySearch(minDeg-2*aDeg, minDeg+2*aDeg, xHi, yHi, zHi, xLo, yLo, zLo)
	if cand < best {
		best = cand
	}
	return best
}

func sortaBinarySearch2(left, right, x1, y1, z1, x2, y2, z2 float64) float64 {
	x := r * math.Cos(left)
	y := r * math.Sin(left)
	leftVal := findConeSurfaceDistance(x1, y1, z1, x, y, 0) + findOptimalSingleBrinkPointDistance(x2, y2, z2, x, y, 0)
	x = r * math.Cos(right)
	y = r * math.Sin(right)
	rightVal := findConeSurfaceDistance(x1, y1, z1, x, y, 0) + findOptimalSingleBrinkPointDistance(x2, y2, z2, x, y, 0)
	for !isZero(right - left) {
		mid := (left + right) / 2
		x = r * math.Cos(mid)
		y = r * math.Sin(mid)
		midVal := findConeSurfaceDistance(x1, y1, z1, x, y, 0) + findOptimalSingleBrinkPointDistance(x2, y2, z2, x, y, 0)
		if leftVal > rightVal {
			left = mid
			leftVal = midVal
		} else {
			right = mid
			rightVal = midVal
		}
	}
	if leftVal < rightVal {
		return leftVal
	}
	return rightVal
}

func findOptimalDoubleBrinkPointDistance(x1, y1, z1, x2, y2, z2 float64) float64 {
	aDeg := pi / 180
	baseDeg := math.Atan2(y1, x1)
	best := sortaBinarySearch2(baseDeg-aDeg, baseDeg+aDeg, x1, y1, z1, x2, y2, z2)
	for k := 2; k <= 3; k++ {
		cand := sortaBinarySearch2(baseDeg-float64(k)*aDeg, baseDeg+float64(k)*aDeg, x1, y1, z1, x2, y2, z2)
		if cand < best {
			best = cand
		}
	}
	return best
}

func solve(r0, h0, x1, y1, z1, x2, y2, z2 float64) float64 {
	r, h = r0, h0
	pt1OnBase := isZero(z1)
	pt2OnBase := isZero(z2)
	var res float64
	if pt1OnBase {
		if pt2OnBase {
			res = euclideanDistance(x1, y1, x2, y2)
		} else {
			res = findOptimalSingleBrinkPointDistance(x2, y2, z2, x1, y1, z1)
			if isZero(x1*x1 + y1*y1 - r*r) {
				tmp := findConeSurfaceDistance(x1, y1, z1, x2, y2, z2)
				if tmp < res {
					res = tmp
				}
			}
		}
	} else if pt2OnBase {
		res = findOptimalSingleBrinkPointDistance(x1, y1, z1, x2, y2, z2)
		if isZero(x2*x2 + y2*y2 - r*r) {
			tmp := findConeSurfaceDistance(x1, y1, z1, x2, y2, z2)
			if tmp < res {
				res = tmp
			}
		}
	} else {
		res = findConeSurfaceDistance(x1, y1, z1, x2, y2, z2)
		tmp := findOptimalDoubleBrinkPointDistance(x1, y1, z1, x2, y2, z2)
		if tmp < res {
			res = tmp
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Println("could not open testcasesE.txt:", err)
		os.Exit(1)
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
		var r0, h0, x1, y1, z1, x2, y2, z2 float64
		if _, err := fmt.Sscan(line, &r0, &h0, &x1, &y1, &z1, &x2, &y2, &z2); err != nil {
			fmt.Printf("bad test case on line %d\n", idx)
			os.Exit(1)
		}
		expect := solve(r0, h0, x1, y1, z1, x2, y2, z2)
		input := fmt.Sprintf("%.6f %.6f\n%.6f %.6f %.6f\n%.6f %.6f %.6f\n", r0, h0, x1, y1, z1, x2, y2, z2)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var errBuf bytes.Buffer
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		var got float64
		if _, err := fmt.Fscan(bytes.NewReader(out.Bytes()), &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx, out.String())
			os.Exit(1)
		}
		if math.Abs(got-expect) > 1e-6 {
			fmt.Printf("test %d failed: expected %.6f got %.6f\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
