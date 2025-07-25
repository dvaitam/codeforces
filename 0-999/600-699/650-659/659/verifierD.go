package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Point struct{ x, y int64 }

func solve(input string) string {
	lines := strings.Fields(input)
	idx := 0
	n := 0
	fmt.Sscan(lines[idx], &n)
	idx++
	pts := make([]Point, n+2)
	for i := 0; i <= n; i++ {
		fmt.Sscan(lines[idx], &pts[i].x)
		idx++
		fmt.Sscan(lines[idx], &pts[i].y)
		idx++
	}
	pts[n+1] = pts[1]
	var area int64
	for i := 0; i < n; i++ {
		area += pts[i].x*pts[i+1].y - pts[i+1].x*pts[i].y
	}
	orient := int64(1)
	if area < 0 {
		orient = -1
	}
	cnt := 0
	for i := 0; i < n; i++ {
		v1x := pts[i+1].x - pts[i].x
		v1y := pts[i+1].y - pts[i].y
		v2x := pts[i+2].x - pts[i+1].x
		v2y := pts[i+2].y - pts[i+1].y
		cross := v1x*v2y - v1y*v2x
		if cross*orient < 0 {
			cnt++
		}
	}
	return fmt.Sprintln(cnt)
}

func rectTest(w, h int) string {
	return fmt.Sprintf("4\n0 0\n0 %d\n%d %d\n%d 0\n0 0\n", h, w, h, w)
}

func lShapeTest(w1, w2, h, cut int) string {
	return fmt.Sprintf("6\n0 0\n0 %d\n%d %d\n%d %d\n%d 0\n0 0\n", h, w1, h, w1, cut, w1+w2)
}

func generateTests() []string {
	rand.Seed(45)
	tests := make([]string, 100)
	for i := 0; i < 50; i++ {
		w := rand.Intn(9) + 2
		h := rand.Intn(9) + 2
		tests[i] = rectTest(w, h)
	}
	for i := 50; i < 100; i++ {
		w1 := rand.Intn(5) + 2
		w2 := rand.Intn(5) + 2
		h := rand.Intn(5) + 3
		cut := rand.Intn(h-1) + 1
		tests[i] = lShapeTest(w1, w2, h, cut)
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expect := strings.TrimSpace(solve(t))
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if expect != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
