package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const grav = 9.8
const tolerance = 1e-4

func prepareProgram(path string) (string, func(), error) {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".go" && ext != ".cpp" && ext != ".cc" && ext != ".cxx" {
		return path, nil, nil
	}
	dir, err := os.MkdirTemp("", "verifier47E-*")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "cand")
	var cmd *exec.Cmd
	if ext == ".go" {
		cmd = exec.Command("go", "build", "-o", bin, path)
	} else {
		cmd = exec.Command("g++", "-O2", "-std=c++17", "-o", bin, path)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("compile error: %v\n%s", err, out.String())
	}
	if runtime.GOOS == "windows" {
		if _, err2 := os.Stat(bin); err2 != nil {
			if _, e3 := os.Stat(bin + ".exe"); e3 == nil {
				bin += ".exe"
			}
		}
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// solveExpected: brute-force reference.
// For each ball, find the closest wall (smallest x) that blocks it.
// Walls with the same x are deduplicated by keeping the tallest.
func solveExpected(V int, angles []float64, walls [][2]float64) [][2]float64 {
	Vf := float64(V)

	// Sort walls by x, then deduplicate keeping max y per x.
	w := make([][2]float64, len(walls))
	copy(w, walls)
	sort.Slice(w, func(i, j int) bool { return w[i][0] < w[j][0] })
	dedup := w[:0]
	for _, wall := range w {
		if len(dedup) > 0 && math.Abs(dedup[len(dedup)-1][0]-wall[0]) < 1e-9 {
			if wall[1] > dedup[len(dedup)-1][1] {
				dedup[len(dedup)-1][1] = wall[1]
			}
		} else {
			dedup = append(dedup, [2]float64{wall[0], wall[1]})
		}
	}

	result := make([][2]float64, len(angles))
	for i, alpha := range angles {
		cosA := math.Cos(alpha)
		sinA := math.Sin(alpha)
		tLand := 2 * Vf * sinA / grav
		xLand := Vf * cosA * tLand

		result[i] = [2]float64{xLand, 0}
		for _, wall := range dedup {
			xi, yi := wall[0], wall[1]
			if xi > xLand+1e-9 {
				break
			}
			t := xi / (Vf * cosA)
			h := Vf*sinA*t - grav*t*t/2
			if h < 0 {
				h = 0
			}
			if h <= yi+1e-9 {
				result[i] = [2]float64{xi, h}
				break
			}
		}
	}
	return result
}

// buildInput formats a test case as the correct solution expects:
// n and V on one line, each angle on its own line, m on one line,
// each wall (xi yi) on its own line.
func buildInput(V int, angles []float64, walls [][2]float64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", len(angles), V))
	for _, a := range angles {
		sb.WriteString(fmt.Sprintf("%.4f\n", a))
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(walls)))
	for _, wall := range walls {
		sb.WriteString(fmt.Sprintf("%.4f %.4f\n", wall[0], wall[1]))
	}
	return sb.String()
}

func parseOutput(s string, n int) ([][2]float64, error) {
	fields := strings.Fields(s)
	if len(fields) < 2*n {
		return nil, fmt.Errorf("expected %d floats, got %d", 2*n, len(fields))
	}
	res := make([][2]float64, n)
	for i := 0; i < n; i++ {
		x, err := strconv.ParseFloat(fields[2*i], 64)
		if err != nil {
			return nil, fmt.Errorf("bad x[%d]: %v", i, err)
		}
		y, err := strconv.ParseFloat(fields[2*i+1], 64)
		if err != nil {
			return nil, fmt.Errorf("bad y[%d]: %v", i, err)
		}
		res[i] = [2]float64{x, y}
	}
	return res, nil
}

func compareResults(exp, got [][2]float64) error {
	for i := range exp {
		dx := math.Abs(exp[i][0] - got[i][0])
		dy := math.Abs(exp[i][1] - got[i][1])
		relX := dx / math.Max(1, math.Abs(exp[i][0]))
		relY := dy / math.Max(1, math.Abs(exp[i][1]))
		if dx > tolerance && relX > tolerance {
			return fmt.Errorf("ball %d: x expected %.6f got %.6f", i+1, exp[i][0], got[i][0])
		}
		if dy > tolerance && relY > tolerance {
			return fmt.Errorf("ball %d: y expected %.6f got %.6f", i+1, exp[i][1], got[i][1])
		}
	}
	return nil
}

func genCase(rng *rand.Rand) (int, []float64, [][2]float64) {
	n := rng.Intn(5) + 1
	V := rng.Intn(90) + 10
	piOver4 := math.Pi / 4
	angles := make([]float64, n)
	for i := range angles {
		// angle strictly in (0.0001, pi/4 - 0.0001), represented with 4 decimal places
		lo := 1
		hi := int((piOver4-0.0001)*10000) - 1
		aInt := lo + rng.Intn(hi-lo+1)
		angles[i] = float64(aInt) / 10000.0
	}
	m := rng.Intn(6) + 1
	walls := make([][2]float64, m)
	for i := range walls {
		// xi in [1, 100] (integer), yi in [0, 100] (integer)
		walls[i] = [2]float64{
			float64(rng.Intn(100) + 1),
			float64(rng.Intn(101)),
		}
	}
	return V, angles, walls
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin, cleanup, err := prepareProgram(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}

	rng := rand.New(rand.NewSource(42))
	const numTests = 200
	for tc := 1; tc <= numTests; tc++ {
		V, angles, walls := genCase(rng)
		input := buildInput(V, angles, walls)
		exp := solveExpected(V, angles, walls)

		got, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", tc, err)
			os.Exit(1)
		}
		gotParsed, err := parseOutput(got, len(angles))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: parse error: %v\noutput: %q\n", tc, err, got)
			os.Exit(1)
		}
		if err := compareResults(exp, gotParsed); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", tc, err, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", numTests)
}
