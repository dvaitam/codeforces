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

type Point struct{ x, y int }

func isRight(p1, p2, p3 Point) bool {
	dx := p2.x - p1.x
	dy := p2.y - p1.y
	a := dx*dx + dy*dy
	dx = p3.x - p2.x
	dy = p3.y - p2.y
	b := dx*dx + dy*dy
	dx = p1.x - p3.x
	dy = p1.y - p3.y
	c := dx*dx + dy*dy
	if a == 0 || b == 0 || c == 0 {
		return false
	}
	return a+b == c || a+c == b || b+c == a
}

func expected(p [3]Point) string {
	if isRight(p[0], p[1], p[2]) {
		return "RIGHT"
	}
	dirs := []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for i := 0; i < 3; i++ {
		orig := p[i]
		for _, d := range dirs {
			p[i] = Point{orig.x + d.x, orig.y + d.y}
			if isRight(p[0], p[1], p[2]) {
				return "ALMOST"
			}
		}
		p[i] = orig
	}
	return "NEITHER"
}

func generateCase(rng *rand.Rand) (string, string) {
	for {
		var p [3]Point
		for i := 0; i < 3; i++ {
			p[i].x = rng.Intn(21) - 10
			p[i].y = rng.Intn(21) - 10
		}
		area := (p[1].x-p[0].x)*(p[2].y-p[0].y) - (p[1].y-p[0].y)*(p[2].x-p[0].x)
		if area != 0 {
			input := fmt.Sprintf("%d %d %d %d %d %d\n", p[0].x, p[0].y, p[1].x, p[1].y, p[2].x, p[2].y)
			exp := expected(p)
			return input, exp
		}
	}
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	res := strings.TrimSpace(out.String())
	if res != exp {
		return fmt.Errorf("expected %s got %s", exp, res)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
