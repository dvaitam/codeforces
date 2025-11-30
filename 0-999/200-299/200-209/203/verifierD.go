package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type point struct {
	x, y, z float64
}

type testCase struct {
	a, b, m    int
	vx, vy, vz float64
}

// Embedded testcases (previously from testcasesD.txt) to keep verifier self contained.
const rawTestcasesD = `
4 10 9 -3 -3 4
8 10 2 4 -5 3
5 9 4 -2 -2 4
9 8 7 5 -4 1
3 9 7 -5 -5 1
8 5 1 -5 -1 -1
10 9 8 0 -4 5
10 9 5 4 -1 2
3 7 4 1 -1 3
3 5 5 -4 -4 -5
10 1 3 5 -4 -4
1 10 1 -5 -1 -2
5 9 1 1 -4 1
2 6 7 -4 -2 1
8 2 6 -4 -4 -3
6 8 7 -1 -5 -3
1 4 8 -3 -1 1
9 10 9 3 -4 -1
7 10 4 5 -3 -1
3 4 5 -4 -4 -1
3 10 9 -4 -4 -3
1 2 7 -5 -5 5
7 10 10 0 -5 -5
8 9 7 0 -5 1
1 8 4 -2 -5 -5
1 5 2 1 -1 2
3 9 10 -5 -3 1
9 10 2 2 -1 2
7 3 10 -1 -1 -1
8 10 1 -3 -1 0
5 1 2 2 -1 1
8 2 10 -4 -2 -1
2 5 10 -3 -2 -5
3 3 6 -3 -2 4
6 6 1 -1 -3 -1
2 5 10 1 -2 2
3 10 3 -1 -1 -3
10 7 10 1 -3 -4
1 2 10 -2 -1 1
1 9 4 -1 -1 4
1 3 6 -1 -1 5
10 9 10 2 -5 0
6 8 7 3 -2 -4
10 2 9 1 -5 -1
2 9 7 -4 -4 -1
10 6 10 2 -3 3
4 5 3 1 -3 -4
9 10 6 5 -2 4
1 2 4 -2 -4 5
6 1 6 -3 -1 3
9 2 4 -4 -1 -1
3 1 7 0 -5 -2
2 2 5 3 -2 -1
6 2 4 -3 -5 4
1 10 7 2 -1 3
5 2 10 -2 -3 -1
10 8 5 1 -5 2
10 8 6 -4 -3 -4
6 10 1 -1 -1 4
1 9 9 -2 -2 2
1 9 1 1 -1 0
7 1 3 1 -4 -1
6 8 5 -1 -1 3
5 4 8 -2 -3 -2
10 3 4 -2 -5 1
10 1 9 5 -1 2
3 2 5 -3 -1 -1
6 7 10 -4 -2 3
5 8 3 -2 -3 -4
5 3 6 -3 -1 0
3 6 3 -2 -3 0
2 6 3 -3 -2 4
5 3 10 0 -2 -1
10 1 10 1 -1 -3
2 7 7 2 -2 1
8 7 5 4 -2 -5
6 2 1 -1 -1 -1
3 8 7 5 -2 4
9 7 6 -4 -1 -4
5 6 1 0 -4 -5
2 7 6 2 -2 5
1 2 8 -4 -4 -4
2 6 2 4 -3 -5
6 10 1 -2 -1 5
8 10 10 4 -4 0
1 2 2 2 -2 3
8 6 7 -2 -5 -3
4 10 8 5 -3 -4
7 5 5 -2 -1 -3
2 4 7 -2 -1 1
8 1 8 -4 -3 -1
10 1 9 2 -1 -3
4 7 6 0 -2 -5
10 1 8 5 -1 -1
4 9 2 -1 -5 5
1 10 10 1 -3 1
5 7 4 0 -1 0
1 5 5 2 -1 1
8 8 6 2 -2 0
9 8 3 5 -1 4
`

func loadTestcases() ([]testCase, error) {
	fields := strings.Fields(rawTestcasesD)
	if len(fields)%6 != 0 {
		return nil, fmt.Errorf("unexpected token count %d (want multiple of 6)", len(fields))
	}
	cases := make([]testCase, 0, len(fields)/6)
	for i := 0; i < len(fields); i += 6 {
		a, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("parse a at token %d (%q): %w", i+1, fields[i], err)
		}
		b, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, fmt.Errorf("parse b at token %d (%q): %w", i+2, fields[i+1], err)
		}
		m, err := strconv.Atoi(fields[i+2])
		if err != nil {
			return nil, fmt.Errorf("parse m at token %d (%q): %w", i+3, fields[i+2], err)
		}
		vx, err := strconv.ParseFloat(fields[i+3], 64)
		if err != nil {
			return nil, fmt.Errorf("parse vx at token %d (%q): %w", i+4, fields[i+3], err)
		}
		vy, err := strconv.ParseFloat(fields[i+4], 64)
		if err != nil {
			return nil, fmt.Errorf("parse vy at token %d (%q): %w", i+5, fields[i+4], err)
		}
		vz, err := strconv.ParseFloat(fields[i+5], 64)
		if err != nil {
			return nil, fmt.Errorf("parse vz at token %d (%q): %w", i+6, fields[i+5], err)
		}
		cases = append(cases, testCase{a: a, b: b, m: m, vx: vx, vy: vy, vz: vz})
	}
	return cases, nil
}

// abs and solve203DCase are lifted directly from 203D.go so the verifier is self contained.
func abs(a float64) float64 {
	return math.Abs(a)
}

func solve203DCase(tc testCase) (float64, float64) {
	aF := float64(tc.a)
	bF := float64(tc.b)
	eps := 1e-9
	INF := math.Inf(1)

	var hitTime [5]float64
	var hitPoints [5]point
	var hitVectors [5]point

	curPoint := point{x: aF * 0.5, y: float64(tc.m), z: 0}
	curVector := point{x: tc.vx, y: tc.vy, z: tc.vz}

	for {
		if abs(curVector.y) < eps {
			hitTime[0] = INF
		} else {
			t := -curPoint.y / curVector.y
			hitTime[0] = t
			hitPoints[0] = point{
				x: curPoint.x + t*curVector.x,
				y: 0,
				z: curPoint.z + t*curVector.z,
			}
		}
		if abs(curVector.z) < eps {
			hitTime[1], hitTime[2] = INF, INF
		} else {
			t1 := (bF - curPoint.z) / curVector.z
			t2 := -curPoint.z / curVector.z
			hitTime[1], hitTime[2] = t1, t2
			v1 := curVector
			v2 := curVector
			v1.z = -v1.z
			v2.z = -v2.z
			hitVectors[1], hitVectors[2] = v1, v2
			hitPoints[1] = point{
				x: curPoint.x + t1*curVector.x,
				y: curPoint.y + t1*curVector.y,
				z: bF,
			}
			hitPoints[2] = point{
				x: curPoint.x + t2*curVector.x,
				y: curPoint.y + t2*curVector.y,
				z: 0,
			}
		}
		if abs(curVector.x) < eps {
			hitTime[3], hitTime[4] = INF, INF
		} else {
			t3 := (aF - curPoint.x) / curVector.x
			t4 := -curPoint.x / curVector.x
			hitTime[3], hitTime[4] = t3, t4
			v3 := curVector
			v4 := curVector
			v3.x = -v3.x
			v4.x = -v4.x
			hitVectors[3], hitVectors[4] = v3, v4
			hitPoints[3] = point{
				x: aF,
				y: curPoint.y + t3*curVector.y,
				z: curPoint.z + t3*curVector.z,
			}
			hitPoints[4] = point{
				x: 0,
				y: curPoint.y + t4*curVector.y,
				z: curPoint.z + t4*curVector.z,
			}
		}
		minind := 0
		for i := 1; i < 5; i++ {
			if hitTime[i] < hitTime[minind] && hitTime[i] > eps {
				minind = i
			}
		}
		if minind == 0 {
			curPoint = hitPoints[0]
			break
		}
		curPoint = hitPoints[minind]
		curVector = hitVectors[minind]
	}
	return curPoint.x, curPoint.z
}

func parseOutput(s string) (float64, float64, error) {
	parts := strings.Fields(s)
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("output should contain two numbers")
	}
	x0, err1 := strconv.ParseFloat(parts[0], 64)
	z0, err2 := strconv.ParseFloat(parts[1], 64)
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("invalid floats")
	}
	return x0, z0, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		expX, expZ := solve203DCase(tc)
		input := fmt.Sprintf("%d %d %d %.10g %.10g %.10g\n", tc.a, tc.b, tc.m, tc.vx, tc.vy, tc.vz)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		gotX, gotZ, err := parseOutput(strings.TrimSpace(out.String()))
		if err != nil {
			fmt.Printf("test %d: cannot parse output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if math.Abs(gotX-expX) > 1e-6 || math.Abs(gotZ-expZ) > 1e-6 {
			fmt.Printf("test %d failed\nexpected: %.6f %.6f\n got: %.6f %.6f\n", idx+1, expX, expZ, gotX, gotZ)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
