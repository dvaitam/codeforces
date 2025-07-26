package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type rect struct{ x1, y1, x2, y2 int }

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) (string, []rect) {
	n := r.Intn(8) + 2
	x1 := r.Intn(10) - 5
	y1 := r.Intn(10) - 5
	x2 := x1 + r.Intn(5) + 1
	y2 := y1 + r.Intn(5) + 1
	rects := make([]rect, n)
	for i := 0; i < n-1; i++ {
		a := x1 - r.Intn(3)
		b := y1 - r.Intn(3)
		c := x2 + r.Intn(3)
		d := y2 + r.Intn(3)
		rects[i] = rect{a, b, c, d}
	}
	if r.Intn(2) == 0 {
		rects[n-1] = rect{x1 - r.Intn(3), y1 - r.Intn(3), x2 + r.Intn(3), y2 + r.Intn(3)}
	} else {
		a := x2 + 1 + r.Intn(3)
		b := y2 + 1 + r.Intn(3)
		c := a + r.Intn(3) + 1
		d := b + r.Intn(3) + 1
		rects[n-1] = rect{a, b, c, d}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, rc := range rects {
		fmt.Fprintf(&sb, "%d %d %d %d\n", rc.x1, rc.y1, rc.x2, rc.y2)
	}
	return sb.String(), rects
}

func contains(r rect, x, y int) bool {
	return r.x1 <= x && x <= r.x2 && r.y1 <= y && y <= r.y2
}

func check(rects []rect, out string) error {
	f := strings.Fields(out)
	if len(f) < 2 {
		return fmt.Errorf("invalid output")
	}
	x, err := strconv.Atoi(f[0])
	if err != nil {
		return err
	}
	y, err := strconv.Atoi(f[1])
	if err != nil {
		return err
	}
	cnt := 0
	for _, rc := range rects {
		if contains(rc, x, y) {
			cnt++
		}
	}
	if cnt < len(rects)-1 {
		return fmt.Errorf("point not in at least n-1 rectangles")
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	cand := os.Args[1]
	for i := 1; i <= 100; i++ {
		in, rects := genCase(rand.New(rand.NewSource(int64(i))))
		out, err := run(cand, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if err := check(rects, out); err != nil {
			fmt.Printf("wrong answer on test %d: %v\ninput:\n%soutput:%s\n", i, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
