package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func hasSolution(x1, y1, x2, y2 int) bool {
	if x1 == x2 || y1 == y2 || abs(x1-x2) == abs(y1-y2) {
		return true
	}
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func validOutput(out string, x1, y1, x2, y2 int) bool {
	out = strings.TrimSpace(out)
	if out == "-1" {
		return !hasSolution(x1, y1, x2, y2)
	}
	var x3, y3, x4, y4 int
	n, err := fmt.Sscanf(out, "%d %d %d %d", &x3, &y3, &x4, &y4)
	if err != nil || n != 4 {
		return false
	}
	if x3 == x4 && y3 == y4 {
		return false
	}
	// verify square with sides parallel to axes
	// gather points
	pts := [][2]int{{x1, y1}, {x2, y2}, {x3, y3}, {x4, y4}}
	// compute unique x and y values
	xs := map[int]bool{}
	ys := map[int]bool{}
	for _, p := range pts {
		xs[p[0]] = true
		ys[p[1]] = true
	}
	if len(xs) != 2 || len(ys) != 2 {
		return false
	}
	var ax, bx int
	i := 0
	for v := range xs {
		if i == 0 {
			ax = v
		} else {
			bx = v
		}
		i++
	}
	var ay, by int
	i = 0
	for v := range ys {
		if i == 0 {
			ay = v
		} else {
			by = v
		}
		i++
	}
	w1 := abs(ax - bx)
	w2 := abs(ay - by)
	if w1 != w2 || w1 == 0 {
		return false
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		x1 := rand.Intn(201) - 100
		y1 := rand.Intn(201) - 100
		x2 := rand.Intn(201) - 100
		y2 := rand.Intn(201) - 100
		if x1 == x2 && y1 == y2 {
			x2++
		}
		input := fmt.Sprintf("%d %d %d %d\n", x1, y1, x2, y2)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		if !validOutput(out, x1, y1, x2, y2) {
			fmt.Printf("wrong answer on test %d\ninput:\n%soutput:\n%s\n", t+1, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
