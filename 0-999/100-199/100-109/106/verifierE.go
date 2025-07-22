package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const eps = 1e-7

type point struct{ x, y, z float64 }

var (
	p      []point
	outer  [4]point
	center point
	radius float64
	nouter int
)

func dissqr(a, b point) float64 {
	dx := a.x - b.x
	dy := a.y - b.y
	dz := a.z - b.z
	return dx*dx + dy*dy + dz*dz
}
func dot(a, b point) float64       { return a.x*b.x + a.y*b.y + a.z*b.z }
func det2(m [2][2]float64) float64 { return m[0][0]*m[1][1] - m[0][1]*m[1][0] }
func det3(m [3][3]float64) float64 {
	return m[0][0]*m[1][1]*m[2][2] + m[0][1]*m[1][2]*m[2][0] + m[0][2]*m[2][1]*m[1][0] - m[0][2]*m[1][1]*m[2][0] - m[0][1]*m[1][0]*m[2][2] - m[0][0]*m[1][2]*m[2][1]
}

func ball() {
	center = point{}
	radius = 0
	switch nouter {
	case 1:
		center = outer[0]
	case 2:
		center.x = (outer[0].x + outer[1].x) / 2
		center.y = (outer[0].y + outer[1].y) / 2
		center.z = (outer[0].z + outer[1].z) / 2
		radius = dissqr(center, outer[0])
	case 3:
		var q [2]point
		for i := 0; i < 2; i++ {
			q[i] = point{outer[i+1].x - outer[0].x, outer[i+1].y - outer[0].y, outer[i+1].z - outer[0].z}
		}
		var m [2][2]float64
		var sol [2]float64
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				m[i][j] = dot(q[i], q[j]) * 2
			}
			sol[i] = dot(q[i], q[i])
		}
		d := det2(m)
		if math.Abs(d) < eps {
			return
		}
		var L [2]float64
		L[0] = (sol[0]*m[1][1] - sol[1]*m[0][1]) / d
		L[1] = (sol[1]*m[0][0] - sol[0]*m[1][0]) / d
		center = point{outer[0].x + q[0].x*L[0] + q[1].x*L[1], outer[0].y + q[0].y*L[0] + q[1].y*L[1], outer[0].z + q[0].z*L[0] + q[1].z*L[1]}
		radius = dissqr(center, outer[0])
	case 4:
		var q [3]point
		var sol [3]float64
		var m [3][3]float64
		for i := 0; i < 3; i++ {
			q[i] = point{outer[i+1].x - outer[0].x, outer[i+1].y - outer[0].y, outer[i+1].z - outer[0].z}
			sol[i] = dot(q[i], q[i])
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				m[i][j] = dot(q[i], q[j]) * 2
			}
		}
		d := det3(m)
		if math.Abs(d) < eps {
			return
		}
		var L [3]float64
		for j := 0; j < 3; j++ {
			for i := 0; i < 3; i++ {
				m[i][j] = sol[i]
			}
			L[j] = det3(m) / d
			for i := 0; i < 3; i++ {
				m[i][j] = dot(q[i], q[j]) * 2
			}
		}
		center = outer[0]
		for i := 0; i < 3; i++ {
			center.x += q[i].x * L[i]
			center.y += q[i].y * L[i]
			center.z += q[i].z * L[i]
		}
		radius = dissqr(center, outer[0])
	}
}

func minball(n int) {
	ball()
	if nouter < 4 {
		for i := 0; i < n; i++ {
			if dissqr(center, p[i])-radius > eps {
				outer[nouter] = p[i]
				nouter++
				minball(i)
				nouter--
				if i > 0 {
					tmp := p[i]
					copy(p[1:i+1], p[0:i])
					p[0] = tmp
				}
			}
		}
	}
}

func solve(points []point) (point, float64) {
	p = make([]point, len(points))
	copy(p, points)
	nouter = 0
	minball(len(p))
	return center, math.Sqrt(radius)
}

func parseTests(path string) ([][]point, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	var res [][]point
	for {
		if !scan.Scan() {
			break
		}
		line := strings.TrimSpace(scan.Text())
		if line == "" {
			continue
		}
		n, _ := strconv.Atoi(line)
		pts := make([]point, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			var x, y, z float64
			fmt.Sscan(scan.Text(), &x, &y, &z)
			pts[i] = point{x, y, z}
		}
		res = append(res, pts)
		scan.Scan() // blank line
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, pts := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(pts)))
		for _, pt := range pts {
			sb.WriteString(fmt.Sprintf("%.2f %.2f %.2f\n", pt.x, pt.y, pt.z))
		}
		_, radExp := solve(pts)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(got)
		if len(fields) != 3 {
			fmt.Fprintf(os.Stderr, "case %d: expected three numbers\n", idx+1)
			os.Exit(1)
		}
		gx, _ := strconv.ParseFloat(fields[0], 64)
		gy, _ := strconv.ParseFloat(fields[1], 64)
		gz, _ := strconv.ParseFloat(fields[2], 64)
		radGot := 0.0
		for _, pt := range pts {
			d := math.Sqrt(dissqr(point{gx, gy, gz}, pt))
			if d > radGot {
				radGot = d
			}
		}
		if math.Abs(radGot-radExp) > 1e-4 {
			fmt.Fprintf(os.Stderr, "case %d failed: radius mismatch\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
