package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

type Point struct {
	x, y int64
}

type testG struct {
	ants    []Point
	queries [][]int64
}

func genTestsG() []testG {
	rand.Seed(122007)
	tests := make([]testG, 100)
	for i := range tests {
		n := rand.Intn(3) + 2
		ants := make([]Point, n)
		for j := range ants {
			ants[j] = Point{int64(rand.Intn(11) - 5), int64(rand.Intn(11) - 5)}
		}
		m := rand.Intn(3) + 1
		queries := make([][]int64, m)
		for q := 0; q < m; q++ {
			arr := make([]int64, n)
			for j := range arr {
				val := rand.Intn(25)
				arr[j] = int64(val)
			}
			sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
			queries[q] = arr
		}
		tests[i] = testG{ants: ants, queries: queries}
	}
	return tests
}

func int64sqrt(x int64) float64 { return math.Sqrt(float64(x)) }

func circleIntersections(p1 Point, r1 float64, p2 Point, r2 float64) []Point {
	dx := float64(p2.x - p1.x)
	dy := float64(p2.y - p1.y)
	d := math.Hypot(dx, dy)
	if d > r1+r2+1e-6 || d < math.Abs(r1-r2)-1e-6 || d == 0 {
		return nil
	}
	a := (r1*r1 - r2*r2 + d*d) / (2 * d)
	h2 := r1*r1 - a*a
	if h2 < -1e-6 {
		return nil
	}
	if h2 < 0 {
		h2 = 0
	}
	h := math.Sqrt(h2)
	xm := float64(p1.x) + a*dx/d
	ym := float64(p1.y) + a*dy/d
	rx := -dy * (h / d)
	ry := dx * (h / d)
	points := []Point{}
	x1 := math.Round(xm + rx)
	y1 := math.Round(ym + ry)
	if math.Abs((xm+rx)-x1) < 1e-6 && math.Abs((ym+ry)-y1) < 1e-6 {
		points = append(points, Point{int64(x1), int64(y1)})
	}
	x2 := math.Round(xm - rx)
	y2 := math.Round(ym - ry)
	if math.Abs((xm-rx)-x2) < 1e-6 && math.Abs((ym-ry)-y2) < 1e-6 {
		p := Point{int64(x2), int64(y2)}
		if len(points) == 0 || points[0] != p {
			points = append(points, p)
		}
	}
	return points
}

func checkPoint(p Point, ants []Point, d []int64) bool {
	n := len(ants)
	arr := make([]int64, n)
	for i, a := range ants {
		dx := p.x - a.x
		dy := p.y - a.y
		arr[i] = dx*dx + dy*dy
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	for i := 0; i < n; i++ {
		if arr[i] != d[i] {
			return false
		}
	}
	return true
}

func solveQuery(ants []Point, d []int64) []Point {
	n := len(ants)
	dmin := d[0]
	dmax := d[n-1]
	var ans []Point
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			pts := circleIntersections(ants[i], int64sqrt(dmin), ants[j], int64sqrt(dmax))
			for _, p := range pts {
				if checkPoint(p, ants, d) {
					ans = append(ans, p)
				}
			}
		}
	}
	if len(ans) > 1 {
		sort.Slice(ans, func(i, j int) bool {
			if ans[i].x == ans[j].x {
				return ans[i].y < ans[j].y
			}
			return ans[i].x < ans[j].x
		})
		uniq := ans[:1]
		for k := 1; k < len(ans); k++ {
			if ans[k] != ans[k-1] {
				uniq = append(uniq, ans[k])
			}
		}
		ans = uniq
	}
	return ans
}

func solveG(tc testG) [][]Point {
	res := make([][]Point, len(tc.queries))
	for i, q := range tc.queries {
		res[i] = solveQuery(tc.ants, q)
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsG()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, len(tc.ants))
		for _, p := range tc.ants {
			fmt.Fprintf(&input, "%d %d\n", p.x, p.y)
		}
		fmt.Fprintln(&input, len(tc.queries))
		for _, q := range tc.queries {
			for i, v := range q {
				if i > 0 {
					input.WriteByte(' ')
				}
				fmt.Fprint(&input, v)
			}
			input.WriteByte('\n')
		}
	}

	expected := make([][][]Point, len(tests))
	for i, tc := range tests {
		expected[i] = solveG(tc)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, out.String())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i, queries := range expected {
		for j, exp := range queries {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d query %d\n", i+1, j+1)
				os.Exit(1)
			}
			cnt, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d query %d\n", i+1, j+1)
				os.Exit(1)
			}
			if cnt != len(exp) {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d query %d\n", i+1, j+1)
				os.Exit(1)
			}
			for k := 0; k < cnt; k++ {
				if !scanner.Scan() {
					fmt.Fprintf(os.Stderr, "wrong output format on test %d query %d\n", i+1, j+1)
					os.Exit(1)
				}
				x, err := strconv.ParseInt(scanner.Text(), 10, 64)
				if err != nil {
					fmt.Fprintf(os.Stderr, "non-integer output on test %d query %d\n", i+1, j+1)
					os.Exit(1)
				}
				if !scanner.Scan() {
					fmt.Fprintf(os.Stderr, "wrong output format on test %d query %d\n", i+1, j+1)
					os.Exit(1)
				}
				y, err := strconv.ParseInt(scanner.Text(), 10, 64)
				if err != nil {
					fmt.Fprintf(os.Stderr, "non-integer output on test %d query %d\n", i+1, j+1)
					os.Exit(1)
				}
				if x != exp[k].x || y != exp[k].y {
					fmt.Fprintf(os.Stderr, "wrong answer on test %d query %d\n", i+1, j+1)
					os.Exit(1)
				}
			}
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
