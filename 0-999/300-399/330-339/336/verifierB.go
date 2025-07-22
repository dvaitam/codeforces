package main

import (
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

func clc(q, r float64) float64 {
	if q < 1.0 {
		return 0.0
	}
	q1 := math.Min(1.0, q)
	base := math.Sqrt(r*r + r*r)
	rs := base*q1 + q1*q1*r
	q -= q1
	return rs + base*q*2.0 + q*q*r
}

func expectedAnswerB(m, R int) float64 {
	n := float64(m)
	r := float64(R)
	rs := 0.0
	for i := 1; i <= m; i++ {
		rs += r + r
		q1 := float64(i - 1)
		rs += r*q1 + clc(q1, r)
		q2 := n - float64(i)
		rs += r*q2 + clc(q2, r)
	}
	rs /= n * n
	return rs
}

func generateCaseB(rng *rand.Rand) (int, int) {
	m := rng.Intn(8) + 1
	R := rng.Intn(10) + 1
	return m, R
}

func runCaseB(bin string, m, R int) error {
	input := fmt.Sprintf("%d %d\n", m, R)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseFloat(gotStr, 64)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expected := expectedAnswerB(m, R)
	if math.Abs(got-expected) > 1e-6 {
		return fmt.Errorf("expected %.6f got %s", expected, gotStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		m, R := generateCaseB(rng)
		if err := runCaseB(bin, m, R); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d\n", i+1, err, m, R)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
