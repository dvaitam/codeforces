package main

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveOne(r, v float64, s, f int) float64 {
	d := float64(f - s)
	e := d / r
	b := e - 2*math.Pi
	if b < 0 {
		b = 0
	}
	for i := 0; i < 50; i++ {
		mid := 0.5 * (b + e)
		t := math.Sin(mid * 0.5)
		if t < 0 {
			t = -t
		}
		t = r*mid + 2*r*t
		if t < d {
			b = mid
		} else {
			e = mid
		}
	}
	return b * r / v
}

func solveB(n int, r, v int, segs [][2]int) []string {
	res := make([]string, n)
	rf := float64(r)
	vf := float64(v)
	for i := 0; i < n; i++ {
		t := solveOne(rf, vf, segs[i][0], segs[i][1])
		res[i] = fmt.Sprintf("%.10f", t)
	}
	return res
}

func run(binary, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(2)
	for t := 0; t < 100; t++ {
		n := rand.Intn(10) + 1
		r := rand.Intn(9) + 1
		v := rand.Intn(9) + 1
		segs := make([][2]int, n)
		for i := 0; i < n; i++ {
			s := rand.Intn(50)
			f := s + rand.Intn(50) + 1
			segs[i] = [2]int{s, f}
		}
		input := fmt.Sprintf("%d %d %d\n", n, r, v)
		for i, sg := range segs {
			input += fmt.Sprintf("%d %d", sg[0], sg[1])
			if i+1 < n {
				input += "\n"
			} else {
				input += "\n"
			}
		}
		expectedLines := solveB(n, r, v, segs)
		expected := strings.Join(expectedLines, "\n")
		got, err := run(bin, input)
		if err != nil {
			fmt.Println("test", t, "runtime error:", err)
			fmt.Println("output:", got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Println("test", t, "failed")
			fmt.Println("input:\n" + input)
			fmt.Println("expected:\n" + expected)
			fmt.Println("got:\n" + got)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
