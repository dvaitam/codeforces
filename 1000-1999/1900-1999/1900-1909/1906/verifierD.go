package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const ref1906D = "1906D.go"

func genConvex(rng *rand.Rand, n int) [][2]int {
	for {
		angles := make([]float64, n)
		for i := 0; i < n; i++ {
			angles[i] = rng.Float64() * 2 * math.Pi
		}
		sort.Float64s(angles)
		pts := make([][2]int, n)
		R := 10000000.0
		for i := 0; i < n; i++ {
			pts[i][0] = int(math.Round(R * math.Cos(angles[i])))
			pts[i][1] = int(math.Round(R * math.Sin(angles[i])))
		}
		ok := true
		for i := 0; i < n; i++ {
			p1 := pts[i]
			p2 := pts[(i+1)%n]
			p3 := pts[(i+2)%n]
			cross := int64(p2[0]-p1[0])*int64(p3[1]-p1[1]) - int64(p2[1]-p1[1])*int64(p3[0]-p1[0])
			if cross <= 0 {
				ok = false
				break
			}
		}
		if ok {
			return pts
		}
	}
}

func isStrictlyInside(pts [][2]int, q [2]int) bool {
	n := len(pts)
	for i := 0; i < n; i++ {
		p1 := pts[i]
		p2 := pts[(i+1)%n]
		cross := int64(p2[0]-p1[0])*int64(q[1]-p1[1]) - int64(p2[1]-p1[1])*int64(q[0]-p1[0])
		if cross <= 0 {
			return false
		}
	}
	return true
}

func generateCase(rng *rand.Rand) ([]byte, int) {
	n := rng.Intn(10) + 3
	pts := genConvex(rng, n)
	
	q := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, p := range pts {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		var a, b [2]int
		for {
			a[0] = rng.Intn(40000000) - 20000000
			a[1] = rng.Intn(40000000) - 20000000
			if !isStrictlyInside(pts, a) {
				break
			}
		}
		for {
			b[0] = rng.Intn(40000000) - 20000000
			b[1] = rng.Intn(40000000) - 20000000
			if !isStrictlyInside(pts, b) && (a[0] != b[0] || a[1] != b[1]) {
				break
			}
		}
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", a[0], a[1], b[0], b[1]))
	}
	return []byte(sb.String()), q
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference(ref1906D)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for testIdx := 0; testIdx < 100; testIdx++ {
		input, q := generateCase(rng)
		
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error: %v\nInput:\n%s\n", err, string(input))
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error: %v\nInput:\n%s\n", err, string(input))
			os.Exit(1)
		}

		refVals, err := parseFloatLines(refOut)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to parse reference output:", err)
			os.Exit(1)
		}
		candVals, err := parseFloatLines(candOut)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to parse candidate output:", err)
			os.Exit(1)
		}
		if len(refVals) != q {
			fmt.Fprintf(os.Stderr, "reference output line count mismatch: expected %d, got %d\n", q, len(refVals))
			os.Exit(1)
		}
		if len(candVals) != q {
			fmt.Fprintf(os.Stderr, "candidate output line count mismatch: expected %d, got %d\n", q, len(candVals))
			os.Exit(1)
		}
		for i := 0; i < q; i++ {
			if !closeEnough(refVals[i], candVals[i]) {
				fmt.Fprintf(os.Stderr, "case %d line %d mismatch: expected %.9f, got %.9f\nInput:\n%s\n", testIdx+1, i+1, refVals[i], candVals[i], string(input))
				os.Exit(1)
			}
		}
	}
	fmt.Println("Accepted")
}

func parseFloatLines(out string) ([]float64, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	vals := make([]float64, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		v, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float %q", line)
		}
		vals = append(vals, v)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return vals, nil
}

func closeEnough(a, b float64) bool {
	if a < 0 && b < 0 {
		return true // Both are -1
	}
	if a < 0 || b < 0 {
		return false // One is -1, other is not
	}
	diff := math.Abs(a - b)
	limit := 1e-6 * math.Max(1.0, math.Abs(a))
	return diff <= limit+1e-9
}

func buildReference(src string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-1906D-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref.bin")
	
	srcPath := src
	if !filepath.IsAbs(srcPath) && !strings.HasPrefix(srcPath, ".") {
		srcPath = "./" + srcPath
	}

	cmd := exec.Command("go", "build", "-o", bin, srcPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}