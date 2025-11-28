package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const eps = 1e-4 // Tolerance for floating point comparisons

type testCase struct {
	ax, ay, bx, by, cx, cy float64
}

// calculateExpectedArea uses the robust iterative approach to determine the expected area.
func calculateExpectedArea(ax, ay, bx, by, cx, cy float64) float64 {
	a := math.Hypot(ax-bx, ay-by)
	b := math.Hypot(bx-cx, by-cy)
	c := math.Hypot(cx-ax, cy-ay)

	p := (a + b + c) / 2
	s := math.Sqrt(p * (p - a) * (p - b) * (p - c))

	// Check for degenerate triangle (collinear points or very small area)
	// The problem statement guarantees a valid polygon, implying non-collinear points.
	// If s is very small, R can become very large. Handle this by returning 0 if area is effectively zero.
	if s < 1e-9 { // Use a smaller tolerance for area to distinguish from actual non-degenerate triangles
		return 0.0
	}

	R := (a * b * c) / (4 * s)

	angles := make([]float64, 3)
	sides := []float64{a, b, c}
	for i, side := range sides {
		val := 1.0 - (side*side)/(2.0*R*R)
		// Clamp values to avoid domain errors for Acos due to floating point inaccuracies
		if val < -1.0 {
			val = -1.0
		}
		if val > 1.0 {
			val = 1.0
		}
		angles[i] = math.Acos(val)
	}

	for n := 3; n <= 100; n++ {
		delta := 2.0 * math.Pi / float64(n)
		allInt := true
		for _, theta := range angles {
			// Check if theta is close to an integer multiple of delta
			quotient := theta / delta
			diff := math.Abs(quotient - math.Round(quotient))
			if diff > eps {
				allInt = false
				break
			}
		}
		if allInt {
			return (float64(n) / 2.0) * R * R * math.Sin(2.0*math.Pi/float64(n))
		}
	}
	return -1.0 // Should not happen based on problem constraints (n <= 100), indicates an error in logic or extreme floating point issues
}

// generateDistinctRandoms generates 'count' distinct random integers less than 'max'
func generateDistinctRandoms(max int, count int) []int {
	if count > max {
		panic("count cannot be greater than max for distinct randoms")
	}
	seen := make(map[int]bool)
	result := make([]int, count)
	for i := 0; i < count; {
		r := rand.Intn(max)
		if !seen[r] {
			seen[r] = true
			result[i] = r
			i++
		}
	}
	return result
}

func main() {
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	var generatedTestCases []testCase

	// Add some basic hardcoded test cases first
	generatedTestCases = append(generatedTestCases, testCase{ax: 0, ay: 0, bx: 1, by: 1, cx: 0, cy: 2})
	generatedTestCases = append(generatedTestCases, testCase{ax: 0, ay: 0, bx: 1, by: 0, cx: 0.5, cy: math.Sqrt(3) / 2})

	// Generate 98 additional random test cases (total 100)
	for i := 0; i < 98; i++ {
		R := rand.Float64()*499.0 + 1.0 // Circumradius between 1.0 and 500.0
		n := rand.Intn(98) + 3           // n between 3 and 100

		dx := rand.Float64()*200.0 - 100.0 // Center offset between -100 and 100
		dy := rand.Float64()*200.0 - 100.0

		vertexIndices := generateDistinctRandoms(n, 3)
		k1, k2, k3 := float64(vertexIndices[0]), float64(vertexIndices[1]), float64(vertexIndices[2])

		// Calculate coordinates for vertices
		ax := R*math.Cos(k1*2.0*math.Pi/float64(n)) + dx
		ay := R*math.Sin(k1*2.0*math.Pi/float64(n)) + dy
		bx := R*math.Cos(k2*2.0*math.Pi/float64(n)) + dx
		by := R*math.Sin(k2*2.0*math.Pi/float64(n)) + dy
		cx := R*math.Cos(k3*2.0*math.Pi/float64(n)) + dx
		cy := R*math.Sin(k3*2.0*math.Pi/float64(n)) + dy

		// Add small noise to coordinates to simulate input precision
		noise := (rand.Float64()*2 - 1) * 1e-7 // +/- 1e-7 noise
		ax += noise; ay += noise
		bx += noise; by += noise
		cx += noise; cy += noise

		generatedTestCases = append(generatedTestCases, testCase{ax: ax, ay: ay, bx: bx, by: by, cx: cx, cy: cy})
	}

	for idx, tc := range generatedTestCases {
		// Calculate expected area using the robust logic
		exp := calculateExpectedArea(tc.ax, tc.ay, tc.bx, tc.by, tc.cx, tc.cy)

		// Run the user's binary
		cmd := exec.Command(binary)
		// Print input with 6 decimal places as per problem statement, and trim trailing zeros from floating points.
		// Sprintf with %g format will drop insignificant trailing zeros, or use %.6f if exact precision is needed. Let's use %.6f to be explicit.
		cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("%.6f %.6f %.6f %.6f %.6f %.6f\n", tc.ax, tc.ay, tc.bx, tc.by, tc.cx, tc.cy))
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Test %d (Input: %.6f %.6f %.6f %.6f %.6f %.6f): runtime error: %v\nstderr: %s\n", idx+1, tc.ax, tc.ay, tc.bx, tc.by, tc.cx, tc.cy, err, errBuf.String())
			os.Exit(1)
		}
		var got float64
		outStr := strings.TrimSpace(outBuf.String())
		_, err = fmt.Sscan(outStr, &got)
		if err != nil {
			fmt.Printf("Test %d (Input: %.6f %.6f %.6f %.6f %.6f %.6f): Failed to parse output '%s': %v\n", idx+1, tc.ax, tc.ay, tc.bx, tc.by, tc.cx, tc.cy, outStr, err)
			os.Exit(1)
		}

		// Compare with tolerance, considering relative and absolute error
		if math.Abs(got-exp) > eps*math.Max(1.0, math.Abs(exp)) {
			fmt.Printf("Test %d (Input: %.6f %.6f %.6f %.6f %.6f %.6f) failed: expected %.8f got %.8f (raw: '%s')\n", idx+1, tc.ax, tc.ay, tc.bx, tc.by, tc.cx, tc.cy, exp, got, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(generatedTestCases))
}
