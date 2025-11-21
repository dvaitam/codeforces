package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSourceB2 = "690B2.go"
	refBinaryB2 = "ref690B2.bin"
	totalTests  = 60
)

type point struct {
	x int
	y int
}

type testCase struct {
	n    int
	grid []string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := generateTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}

		refPoly, err := parseOutput(refOut)
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}
		candPoly, err := parseOutput(candOut)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if len(refPoly) != len(candPoly) {
			fmt.Printf("test %d failed: expected polygon with %d vertices, got %d\n", idx+1, len(refPoly), len(candPoly))
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}

		match := true
		for i := range refPoly {
			if refPoly[i] != candPoly[i] {
				match = false
				break
			}
		}
		if !match {
			fmt.Printf("test %d failed: polygon vertices differ\n", idx+1)
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryB2, refSourceB2)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryB2), nil
}

func runProgram(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for _, row := range tc.grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	sb.WriteString("0\n")
	return sb.String()
}

func parseOutput(out string) ([]point, error) {
	if strings.TrimSpace(out) == "" {
		return nil, fmt.Errorf("empty output")
	}
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	v, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid vertex count %q: %v", fields[0], err)
	}
	if v < 0 {
		return nil, fmt.Errorf("negative vertex count")
	}
	expected := 1 + 2*v
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]point, v)
	for i := 0; i < v; i++ {
		x, err := strconv.Atoi(fields[1+2*i])
		if err != nil {
			return nil, fmt.Errorf("invalid x coordinate: %v", err)
		}
		y, err := strconv.Atoi(fields[1+2*i+1])
		if err != nil {
			return nil, fmt.Errorf("invalid y coordinate: %v", err)
		}
		res[i] = point{x: x, y: y}
	}
	return res, nil
}

func generateTests() []testCase {
	tests := []testCase{
		buildSimpleSquare(6),
		buildSimpleTriangle(7),
		buildSimplePentagon(10),
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-5 {
		n := rnd.Intn(40) + 10
		tc := randomTestCase(n, rnd)
		tests = append(tests, tc)
	}
	tests = append(tests,
		randomTestCase(60, rand.New(rand.NewSource(1))),
		randomTestCase(100, rand.New(rand.NewSource(2))),
		randomTestCase(200, rand.New(rand.NewSource(3))),
		randomTestCase(400, rand.New(rand.NewSource(4))),
		randomTestCase(500, rand.New(rand.NewSource(5))),
	)
	return tests
}

func buildSimpleSquare(n int) testCase {
	poly := []point{{2, 2}, {2, n - 2}, {n - 2, n - 2}, {n - 2, 2}}
	return buildTestCase(n, poly)
}

func buildSimpleTriangle(n int) testCase {
	poly := []point{{2, 2}, {n - 2, 2}, {n/2 + 1, n - 2}}
	return buildTestCase(n, poly)
}

func buildSimplePentagon(n int) testCase {
	poly := []point{
		{2, 2},
		{2, n - 3},
		{n/2 + 1, n - 2},
		{n - 2, n/2 + 1},
		{n - 3, 2},
	}
	return buildTestCase(n, poly)
}

func randomTestCase(n int, rnd *rand.Rand) testCase {
	for {
		points := make([]point, 0, 40)
		count := rnd.Intn(20) + 6
		for i := 0; i < count; i++ {
			points = append(points, point{
				x: rnd.Intn(n-3) + 2,
				y: rnd.Intn(n-3) + 2,
			})
		}
		hull := convexHull(points)
		if len(hull) >= 3 {
			return buildTestCase(n, hull)
		}
	}
}

func buildTestCase(n int, poly []point) testCase {
	poly = normalizePolygon(poly)
	grid := buildGrid(n, poly)
	return testCase{n: n, grid: grid}
}

func buildGrid(n int, poly []point) []string {
	rows := make([]string, n)
	for row := 0; row < n; row++ {
		y := n - row
		var sb strings.Builder
		for x := 1; x <= n; x++ {
			count := 0
			corners := []point{
				{x - 1, y - 1},
				{x, y - 1},
				{x - 1, y},
				{x, y},
			}
			for _, c := range corners {
				if pointInConvex(c, poly) {
					count++
				}
			}
			sb.WriteByte(byte('0' + count))
		}
		rows[row] = sb.String()
	}
	return rows
}

func normalizePolygon(poly []point) []point {
	if len(poly) == 0 {
		return poly
	}
	area := int64(0)
	for i := 0; i < len(poly); i++ {
		j := (i + 1) % len(poly)
		area += int64(poly[i].x*poly[j].y - poly[j].x*poly[i].y)
	}
	if area > 0 {
		for i, j := 0, len(poly)-1; i < j; i, j = i+1, j-1 {
			poly[i], poly[j] = poly[j], poly[i]
		}
	} else if area == 0 {
		// ensure strict convex order by re-running hull
		poly = convexHull(poly)
	}
	start := 0
	for i := 1; i < len(poly); i++ {
		if poly[i].x < poly[start].x || (poly[i].x == poly[start].x && poly[i].y < poly[start].y) {
			start = i
		}
	}
	res := make([]point, len(poly))
	for i := 0; i < len(poly); i++ {
		idx := (start - i + len(poly)) % len(poly)
		res[i] = poly[idx]
	}
	return res
}

func convexHull(points []point) []point {
	if len(points) <= 1 {
		cp := make([]point, len(points))
		copy(cp, points)
		return cp
	}
	pts := make([]point, len(points))
	copy(pts, points)
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].x == pts[j].x {
			return pts[i].y < pts[j].y
		}
		return pts[i].x < pts[j].x
	})
	lower := make([]point, 0, len(pts))
	for _, p := range pts {
		for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], p) <= 0 {
			lower = lower[:len(lower)-1]
		}
		lower = append(lower, p)
	}
	upper := make([]point, 0, len(pts))
	for i := len(pts) - 1; i >= 0; i-- {
		p := pts[i]
		for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
			upper = upper[:len(upper)-1]
		}
		upper = append(upper, p)
	}
	hull := append(lower, upper[1:len(upper)-1]...)
	return hull
}

func cross(a, b, c point) int {
	return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

func pointInConvex(p point, poly []point) bool {
	if len(poly) == 0 {
		return false
	}
	sign := 0
	for i := 0; i < len(poly); i++ {
		j := (i + 1) % len(poly)
		c := cross(poly[i], poly[j], p)
		if c == 0 {
			continue
		}
		if sign == 0 {
			if c > 0 {
				sign = 1
			} else {
				sign = -1
			}
		} else if (c > 0 && sign < 0) || (c < 0 && sign > 0) {
			return false
		}
	}
	return true
}

func printInput(in string) {
	fmt.Println("Input used:")
	fmt.Print(in)
}
