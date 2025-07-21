package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	errorTol = 1e-4
	infinite = 1e9
)

func calTime(x, v []float64, i, j int) float64 {
	if math.Abs(v[i]-v[j]) > errorTol && math.Abs(x[j]-x[i]) > errorTol && (v[i]-v[j])*(x[j]-x[i]) > 0 {
		return (x[j] - x[i]) / (v[i] - v[j])
	}
	return infinite
}

func collide(v []float64, m []int, i, j int) {
	v1, v2 := v[i], v[j]
	m1, m2 := float64(m[i]), float64(m[j])
	v[i] = ((m1-m2)*v1 + 2*m2*v2) / (m1 + m2)
	v[j] = ((m2-m1)*v2 + 2*m1*v1) / (m1 + m2)
}

func passTime(x, v []float64, t float64) {
	for i := range x {
		x[i] += v[i] * t
	}
}

func solve(in string) string {
	scanner := bufio.NewScanner(strings.NewReader(in))
	scanner.Split(bufio.ScanWords)
	var n int
	var T float64
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &n)
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &T)
	x := make([]float64, n)
	v := make([]float64, n)
	m := make([]int, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &x[i])
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &v[i])
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &m[i])
	}
	for T > 0 {
		colTime := infinite
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				t := calTime(x, v, i, j)
				if t < colTime && t < T {
					colTime = t
				}
			}
		}
		if colTime == infinite {
			passTime(x, v, T)
			break
		}
		passTime(x, v, colTime)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if math.Abs(x[i]-x[j]) < errorTol {
					collide(v, m, i, j)
				}
			}
		}
		T -= colTime
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%.6f\n", x[i]))
	}
	sb.WriteByte('\n')
	return strings.TrimSpace(sb.String())
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase() string {
	n := rand.Intn(5) + 1
	T := rand.Float64() * 100
	coords := map[int]bool{}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %.6f\n", n, T))
	for i := 0; i < n; i++ {
		var x int
		for {
			x = rand.Intn(201) - 100
			if !coords[x] {
				coords[x] = true
				break
			}
		}
		v := rand.Intn(200) + 1
		if rand.Intn(2) == 0 {
			v = -v
		}
		m := rand.Intn(100) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", x, v, m))
	}
	return sb.String()
}

func compareFloats(a, b string) bool {
	sa := strings.Fields(a)
	sb := strings.Fields(b)
	if len(sa) != len(sb) {
		return false
	}
	for i := 0; i < len(sa); i++ {
		fa, _ := strconv.ParseFloat(sa[i], 64)
		fb, _ := strconv.ParseFloat(sb[i], 64)
		if math.Abs(fa-fb) > 1e-3 {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := generateCase()
		exp := solve(tc)
		got, err := run(bin, tc)
		if err != nil {
			fmt.Printf("case %d: error executing binary: %v\n", i+1, err)
			os.Exit(1)
		}
		if !compareFloats(got, exp) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
