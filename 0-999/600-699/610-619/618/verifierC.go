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

func run(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand, n int) (string, [][2]int) {
	points := make([][2]int, n)
	// ensure first three points are non-collinear
	points[0] = [2]int{0, 0}
	points[1] = [2]int{1, 0}
	points[2] = [2]int{0, 1}
	used := map[[2]int]bool{{0, 0}: true, {1, 0}: true, {0, 1}: true}
	for i := 3; i < n; i++ {
		for {
			x := rng.Intn(2001) - 1000
			y := rng.Intn(2001) - 1000
			pt := [2]int{x, y}
			if !used[pt] {
				points[i] = pt
				used[pt] = true
				break
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", points[i][0], points[i][1]))
	}
	return sb.String(), points
}

func cross(ax, ay, bx, by int64) int64 {
	return ax*by - ay*bx
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 120; t++ {
		n := rng.Intn(20) + 3
		input, points := genCase(rng, n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		fields := strings.Fields(got)
		if len(fields) < 3 {
			fmt.Fprintf(os.Stderr, "case %d failed: not enough output values\ninput:\n%s", t+1, input)
			os.Exit(1)
		}
		a, err1 := strconv.Atoi(fields[0])
		b, err2 := strconv.Atoi(fields[1])
		c, err3 := strconv.Atoi(fields[2])
		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: output is not three integers\ninput:\n%soutput:\n%s", t+1, input, got)
			os.Exit(1)
		}
		if a < 1 || a > n || b < 1 || b > n || c < 1 || c > n {
			fmt.Fprintf(os.Stderr, "case %d failed: indices out of range\ninput:\n%soutput:\n%s", t+1, input, got)
			os.Exit(1)
		}
		if a == b || a == c || b == c {
			fmt.Fprintf(os.Stderr, "case %d failed: indices are not distinct\ninput:\n%soutput:\n%s", t+1, input, got)
			os.Exit(1)
		}
		ax, ay := int64(points[a-1][0]), int64(points[a-1][1])
		bx, by := int64(points[b-1][0]), int64(points[b-1][1])
		cx, cy := int64(points[c-1][0]), int64(points[c-1][1])
		area := cross(bx-ax, by-ay, cx-ax, cy-ay)
		if area == 0 {
			fmt.Fprintf(os.Stderr, "case %d failed: chosen points are collinear\ninput:\n%soutput:\n%s", t+1, input, got)
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			if i == a-1 || i == b-1 || i == c-1 {
				continue
			}
			px, py := int64(points[i][0]), int64(points[i][1])
			f1 := cross(bx-ax, by-ay, px-ax, py-ay)
			f2 := cross(cx-bx, cy-by, px-bx, py-by)
			f3 := cross(ax-cx, ay-cy, px-cx, py-cy)
			if (f1 >= 0 && f2 >= 0 && f3 >= 0) || (f1 <= 0 && f2 <= 0 && f3 <= 0) {
				fmt.Fprintf(os.Stderr, "case %d failed: point %d lies inside triangle %d %d %d\ninput:\n%s", t+1, i+1, a, b, c, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
