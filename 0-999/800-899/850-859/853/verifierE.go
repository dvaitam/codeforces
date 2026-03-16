package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func cross(ax, ay, bx, by int64) int64 {
	return ax*by - ay*bx
}

func solveCase(line string) string {
	fields := strings.Fields(line)
	idx := 0
	if len(fields) < 7 {
		return ""
	}
	k, _ := strconv.Atoi(fields[idx])
	idx++
	n, _ := strconv.Atoi(fields[idx])
	idx++
	q, _ := strconv.Atoi(fields[idx])
	idx++
	vx := make([]int64, k)
	vy := make([]int64, k)
	for i := 0; i < k; i++ {
		vx[i], _ = strconv.ParseInt(fields[idx], 10, 64)
		idx++
		vy[i], _ = strconv.ParseInt(fields[idx], 10, 64)
		idx++
	}
	S := make([]float64, k)
	for j := 0; j < k; j++ {
		sum := int64(0)
		for i := 0; i < k; i++ {
			if i == j {
				continue
			}
			sum += abs64(cross(vx[j], vy[j], vx[i], vy[i]))
		}
		S[j] = float64(sum)
		if S[j] == 0 {
			S[j] = 1
		}
	}
	type Point struct {
		coord []float64
		w     int64
	}
	factories := make([]Point, n)
	for i := 0; i < n; i++ {
		fx, _ := strconv.ParseInt(fields[idx], 10, 64)
		idx++
		fy, _ := strconv.ParseInt(fields[idx], 10, 64)
		idx++
		a, _ := strconv.ParseInt(fields[idx], 10, 64)
		idx++
		coord := make([]float64, k)
		for j := 0; j < k; j++ {
			coord[j] = float64(cross(vx[j], vy[j], fx, fy)) / S[j]
		}
		factories[i] = Point{coord: coord, w: a}
	}
	type Query struct {
		center []float64
		t      float64
	}
	queries := make([]Query, q)
	for i := 0; i < q; i++ {
		px, _ := strconv.ParseInt(fields[idx], 10, 64)
		idx++
		py, _ := strconv.ParseInt(fields[idx], 10, 64)
		idx++
		t, _ := strconv.ParseInt(fields[idx], 10, 64)
		idx++
		c := make([]float64, k)
		for j := 0; j < k; j++ {
			c[j] = float64(cross(vx[j], vy[j], px, py)) / S[j]
		}
		queries[i] = Query{center: c, t: float64(t)}
	}
	var sb strings.Builder
	for _, qu := range queries {
		var sum int64
		for _, f := range factories {
			dist := 0.0
			for j := 0; j < k; j++ {
				d := f.coord[j] - qu.center[j]
				if d < 0 {
					d = -d
				}
				if d > dist {
					dist = d
				}
			}
			if dist <= qu.t {
				sum += f.w
			}
		}
		sb.WriteString(fmt.Sprintf("%d\n", sum))
	}
	return strings.TrimSpace(sb.String())
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesRaw = `2 1 1 1 0 0 1 -4 1 2 -1 5 0
2 2 3 1 0 0 1 4 -1 5 -5 4 1 1 1 4 0 -2 1 -3 5 5
2 1 1 1 0 0 1 3 4 2 -1 0 3
2 3 1 1 0 0 1 1 -3 5 4 -3 3 -3 5 1 0 -3 0
2 3 2 1 0 0 1 4 -2 1 3 -4 4 3 4 4 -1 0 4 -4 -2 0
2 3 2 1 0 0 1 -1 3 3 -2 4 3 -1 -4 4 -3 -5 5 -4 2 5
2 2 2 1 0 0 1 -1 2 1 0 0 2 0 0 2 1 -2 0
2 2 3 1 0 0 1 4 3 3 1 4 1 3 -4 2 -5 1 1 -4 -3 0
2 3 3 1 0 0 1 -1 -2 3 -2 4 1 -2 5 4 0 5 5 5 -1 0 -4 2 2
2 3 2 1 0 0 1 -2 -3 4 -5 -3 1 -3 4 1 -4 0 0 5 -2 1
2 2 2 1 0 0 1 -1 4 3 5 3 3 5 2 4 3 -5 1
2 1 3 1 0 0 1 -2 -2 1 -3 -4 0 4 3 0 4 4 3
2 2 3 1 0 0 1 2 1 1 4 5 5 5 0 5 3 -2 1 4 -1 4
2 2 3 1 0 0 1 1 4 5 -5 2 5 4 -2 3 -1 3 5 -1 -4 2
2 3 3 1 0 0 1 -1 3 3 2 2 1 4 -1 3 -4 1 2 -5 -4 5 0 4 0
2 3 1 1 0 0 1 4 5 1 -5 -5 3 2 -4 5 5 -5 3
2 3 1 1 0 0 1 -1 3 3 -1 -2 4 4 2 4 -1 0 0
2 3 3 1 0 0 1 4 -1 1 0 -4 2 1 2 3 0 -1 4 5 -2 5 -3 3 3
2 2 2 1 0 0 1 2 -4 1 5 -3 5 -5 -1 2 0 0 0
2 2 3 1 0 0 1 0 -3 3 2 -5 5 3 -1 2 -1 4 5 0 -4 2
2 1 2 1 0 0 1 4 -3 3 4 -1 5 -1 2 0
2 2 1 1 0 0 1 -2 0 2 -4 5 1 1 2 3
2 3 3 1 0 0 1 -5 0 1 5 -5 4 -1 -2 2 2 -4 5 1 3 0 -1 -4 1
2 1 1 1 0 0 1 -1 1 5 5 4 2
2 1 2 1 0 0 1 2 4 3 -1 0 1 -5 3 1
2 2 3 1 0 0 1 -5 5 4 0 -5 5 -2 -4 4 -2 -5 5 3 0 1
2 3 1 1 0 0 1 4 2 5 2 1 2 2 -1 2 -5 0 5
2 2 3 1 0 0 1 0 5 3 -2 1 1 -5 5 0 -3 4 2 4 -3 0
2 1 3 1 0 0 1 -3 4 2 -4 -3 1 3 -3 3 1 -5 2
2 2 2 1 0 0 1 4 -5 2 4 -5 3 0 0 0 4 5 5
2 2 2 1 0 0 1 -4 1 2 0 4 3 -3 -2 4 -3 4 1
2 3 1 1 0 0 1 0 4 5 3 5 5 -4 4 2 0 -1 3
2 2 2 1 0 0 1 -3 1 4 0 -4 3 -3 -1 1 0 2 0
2 2 2 1 0 0 1 -1 5 3 -3 1 2 -2 -3 4 2 4 3
2 3 3 1 0 0 1 -2 -5 3 -3 -5 4 1 4 1 -5 1 2 4 -2 3 -5 4 1
2 2 3 1 0 0 1 2 2 5 1 1 4 -5 -1 1 -3 0 0 3 4 3
2 1 3 1 0 0 1 1 -2 2 -1 -5 3 0 2 2 3 -2 2
2 3 1 1 0 0 1 2 -3 3 -2 0 1 -1 4 4 -2 -5 1
2 3 2 1 0 0 1 -3 -1 5 3 -2 2 1 4 4 0 5 3 4 -3 0
2 3 3 1 0 0 1 1 5 4 0 -1 2 -1 2 3 1 4 1 2 -2 5 -3 2 3
2 1 3 1 0 0 1 1 -2 5 -1 1 0 -1 -1 1 -3 -2 2
2 1 1 1 0 0 1 0 -2 5 4 -3 0
2 1 1 1 0 0 1 -1 4 3 -1 -3 0
2 2 3 1 0 0 1 0 -4 4 -1 0 5 3 3 5 -1 -4 4 3 -5 0
2 2 3 1 0 0 1 1 4 4 0 5 5 -2 -4 3 1 5 0 4 -1 4
2 2 3 1 0 0 1 -3 -5 1 -3 1 3 -5 0 2 3 4 4 1 3 5
2 1 1 1 0 0 1 4 3 2 -2 0 5
2 3 2 1 0 0 1 -3 -5 1 -5 5 3 -2 1 5 5 -4 3 -2 0 0
2 1 3 1 0 0 1 -3 -1 2 5 -3 1 1 1 1 -2 -5 4
2 2 2 1 0 0 1 -2 3 1 4 0 3 -1 3 1 2 4 1
2 3 1 1 0 0 1 4 5 1 -1 1 1 -5 -4 3 0 -4 4
2 2 1 1 0 0 1 -2 -1 5 -4 -4 4 2 -3 2
2 1 1 1 0 0 1 -2 -3 1 4 -5 3
2 2 3 1 0 0 1 3 -2 4 -2 -2 5 -4 0 1 -1 -4 3 2 -2 5
2 2 2 1 0 0 1 2 -2 4 5 -3 5 -2 4 3 0 1 2
2 2 1 1 0 0 1 -3 4 5 2 -4 2 4 0 5
2 3 2 1 0 0 1 -5 -1 4 -2 1 2 -1 4 3 -5 0 0 1 -5 4
2 3 3 1 0 0 1 -5 -5 1 1 4 3 2 0 1 3 1 0 1 -1 5 -5 0 0
2 1 2 1 0 0 1 5 5 1 4 -5 5 3 5 0
2 3 3 1 0 0 1 -5 -2 5 4 1 3 -4 -4 2 -4 -3 4 -2 -2 3 4 -1 3
2 1 1 1 0 0 1 3 3 4 -4 -2 1
2 2 3 1 0 0 1 5 2 1 -2 -2 1 -3 1 5 2 1 2 4 -2 4
2 1 3 1 0 0 1 5 1 2 2 4 1 -5 -1 3 4 2 5
2 2 3 1 0 0 1 1 0 5 -2 -4 3 3 3 1 0 -1 1 3 1 2
2 2 2 1 0 0 1 0 0 2 -4 -1 3 -2 -4 1 1 3 1
2 2 1 1 0 0 1 1 2 2 -1 3 3 1 -4 4
2 2 2 1 0 0 1 5 -4 2 5 -4 5 -1 -1 0 1 1 1
2 1 3 1 0 0 1 -3 -5 4 1 0 3 1 -5 2 2 5 4
2 3 1 1 0 0 1 -4 -5 2 4 3 2 -3 0 4 -3 3 2
2 2 1 1 0 0 1 5 -4 3 -3 4 2 -2 -3 0
2 3 1 1 0 0 1 2 -1 3 3 -5 4 2 4 2 -5 -3 5
2 1 3 1 0 0 1 1 0 1 0 -5 5 2 3 1 -1 1 0
2 2 2 1 0 0 1 2 -3 2 5 -4 5 3 -4 4 -2 5 5
2 3 3 1 0 0 1 1 0 3 -2 4 1 0 5 3 -4 4 4 -3 -4 3 -5 2 2
2 2 2 1 0 0 1 -3 -5 5 -3 -3 4 -2 -1 4 3 4 0
2 1 1 1 0 0 1 -5 4 3 -1 2 5
2 3 2 1 0 0 1 -1 3 3 0 -3 1 5 2 5 -4 -4 1 -3 0 3
2 1 3 1 0 0 1 -1 0 3 -1 3 5 -5 3 0 5 1 5
2 1 3 1 0 0 1 1 2 1 -4 2 5 -3 1 5 -3 1 4
2 1 2 1 0 0 1 1 -1 3 -3 -4 5 4 2 4
2 2 3 1 0 0 1 -2 3 1 -2 2 5 1 -3 4 3 -5 5 -2 -1 4
2 1 2 1 0 0 1 -5 3 3 -1 1 1 -5 -2 4
2 2 1 1 0 0 1 2 -1 1 -5 1 4 3 -5 5
2 1 3 1 0 0 1 -4 0 4 -2 -1 3 3 0 4 -3 -4 2
2 1 3 1 0 0 1 -5 1 4 -1 0 1 3 -2 4 4 4 0
2 2 1 1 0 0 1 -2 -1 2 -3 2 2 -5 -4 1
2 2 2 1 0 0 1 1 5 1 -3 4 3 -1 5 0 4 -2 2
2 2 2 1 0 0 1 -4 4 3 -4 5 4 0 -2 4 4 -5 0
2 2 3 1 0 0 1 2 -2 3 -1 -4 1 -4 5 4 1 -2 3 -2 2 4
2 3 3 1 0 0 1 1 3 4 -3 -3 1 3 2 1 1 5 5 0 -3 0 4 4 2
2 2 3 1 0 0 1 -2 2 4 -5 -3 3 -1 3 2 5 1 4 -5 5 3
2 1 2 1 0 0 1 -1 2 5 0 2 4 -4 5 2
2 1 2 1 0 0 1 2 0 4 5 -3 5 2 0 0
2 2 3 1 0 0 1 4 -2 5 -3 5 5 -1 5 0 -3 -5 4 4 5 5
2 3 2 1 0 0 1 3 2 3 -3 -5 4 0 -1 4 -5 -5 0 -2 1 5
2 1 1 1 0 0 1 1 2 5 -2 -4 2
2 1 1 1 0 0 1 2 4 1 -4 -5 2
2 1 2 1 0 0 1 -3 -2 2 5 -2 4 5 -5 0
2 2 2 1 0 0 1 0 -3 1 -5 5 5 -5 -5 5 -2 1 2
2 1 2 1 0 0 1 -1 -4 5 -3 -4 2 0 5 4`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		expected := solveCase(line)
		got, err := run(bin, line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
