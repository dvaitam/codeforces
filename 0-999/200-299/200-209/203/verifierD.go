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

func reflectPos(pos, limit float64) float64 {
	period := 2 * limit
	pos = math.Mod(pos, period)
	if pos < 0 {
		pos += period
	}
	if pos > limit {
		return period - pos
	}
	return pos
}

func computeExpected(line string) (float64, float64, error) {
	parts := strings.Fields(line)
	if len(parts) != 6 {
		return 0, 0, fmt.Errorf("input should contain six numbers")
	}
	vals := make([]float64, 6)
	for i, p := range parts {
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid number: %v", err)
		}
		vals[i] = v
	}
	a, b, m := vals[0], vals[1], vals[2]
	vx, vy, vz := vals[3], vals[4], vals[5]
	t := m / -vy
	xUnfold := a/2 + vx*t
	zUnfold := vz * t
	x0 := reflectPos(xUnfold, a)
	z0 := reflectPos(zUnfold, b)
	return x0, z0, nil
}

func parseFloats(s string) (float64, float64, error) {
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

	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
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
		expX, expZ, err := computeExpected(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle compute error: %v\n", err)
			os.Exit(1)
		}
		input := line + "\n"

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		gotX, gotZ, err := parseFloats(strings.TrimSpace(out.String()))
		if err != nil {
			fmt.Printf("test %d: cannot parse output: %v\n", idx, err)
			os.Exit(1)
		}
		if math.Abs(gotX-expX) > 1e-6 || math.Abs(gotZ-expZ) > 1e-6 {
			fmt.Printf("test %d failed\nexpected: %.6f %.6f\n got: %.6f %.6f\n", idx, expX, expZ, gotX, gotZ)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
