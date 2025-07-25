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

type point struct{ x, y float64 }

func expectedC(n int, w, v, u float64, pts []point) float64 {
	ratio := v / u
	direct := true
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		xi, yi := pts[i].x, pts[i].y
		xj, yj := pts[j].x, pts[j].y
		if yi == yj {
			continue
		}
		A := (xj - xi) / (yj - yi)
		denom := A - ratio
		if math.Abs(denom) < 1e-12 {
			continue
		}
		y := (A*yi - xi) / denom
		if y > math.Min(yi, yj) && y < math.Max(yi, yj) {
			direct = false
			break
		}
	}
	tPed := w / u
	if direct {
		return tPed
	}
	maxT := -1e300
	for i := 0; i < n; i++ {
		t := pts[i].x / v
		if t > maxT {
			maxT = t
		}
	}
	return maxT + tPed
}

func runCase(exe string, tc struct {
	n       int
	w, v, u float64
	pts     []point
}) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %.0f %.0f %.0f\n", tc.n, tc.w, tc.v, tc.u))
	for i := 0; i < tc.n; i++ {
		sb.WriteString(fmt.Sprintf("%.2f %.2f\n", tc.pts[i].x, tc.pts[i].y))
	}
	input := sb.String()
	exp := fmt.Sprintf("%.10f", expectedC(tc.n, tc.w, tc.v, tc.u, tc.pts))
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(sc.Text())
	for idx := 0; idx < t; idx++ {
		sc.Scan()
		n, _ := strconv.Atoi(sc.Text())
		sc.Scan()
		w, _ := strconv.ParseFloat(sc.Text(), 64)
		sc.Scan()
		vVal, _ := strconv.ParseFloat(sc.Text(), 64)
		sc.Scan()
		uVal, _ := strconv.ParseFloat(sc.Text(), 64)
		pts := make([]point, n)
		for i := 0; i < n; i++ {
			sc.Scan()
			x, _ := strconv.ParseFloat(sc.Text(), 64)
			sc.Scan()
			y, _ := strconv.ParseFloat(sc.Text(), 64)
			pts[i] = point{x, y}
		}
		tc := struct {
			n       int
			w, v, u float64
			pts     []point
		}{n, w, vVal, uVal, pts}
		if err := runCase(exe, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
