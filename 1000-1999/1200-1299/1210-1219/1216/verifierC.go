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

type Rect struct {
	x1, y1, x2, y2 int
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
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

func area(r Rect) int {
	if r.x2 <= r.x1 || r.y2 <= r.y1 {
		return 0
	}
	return (r.x2 - r.x1) * (r.y2 - r.y1)
}

func intersect(a, b Rect) Rect {
	x1 := max(a.x1, b.x1)
	y1 := max(a.y1, b.y1)
	x2 := min(a.x2, b.x2)
	y2 := min(a.y2, b.y2)
	if x1 >= x2 || y1 >= y2 {
		return Rect{x1, y1, x1, y1}
	}
	return Rect{x1, y1, x2, y2}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func expected(w, b1, b2 Rect) string {
	wArea := area(w)
	cov1 := area(intersect(w, b1))
	cov2 := area(intersect(w, b2))
	covBoth := area(intersect(intersect(w, b1), b2))
	covered := cov1 + cov2 - covBoth
	if covered < wArea {
		return "YES"
	}
	return "NO"
}

func genCase(rng *rand.Rand) (string, Rect, Rect, Rect) {
	x1 := rng.Intn(10)
	y1 := rng.Intn(10)
	x2 := x1 + rng.Intn(10) + 1
	y2 := y1 + rng.Intn(10) + 1
	w := Rect{x1, y1, x2, y2}

	x3 := rng.Intn(20)
	y3 := rng.Intn(20)
	x4 := x3 + rng.Intn(10) + 1
	y4 := y3 + rng.Intn(10) + 1
	b1 := Rect{x3, y3, x4, y4}

	x5 := rng.Intn(20)
	y5 := rng.Intn(20)
	x6 := x5 + rng.Intn(10) + 1
	y6 := y5 + rng.Intn(10) + 1
	b2 := Rect{x5, y5, x6, y6}

	input := fmt.Sprintf("%d %d %d %d\n%d %d %d %d\n%d %d %d %d\n",
		w.x1, w.y1, w.x2, w.y2,
		b1.x1, b1.y1, b1.x2, b1.y2,
		b2.x1, b2.y1, b2.x2, b2.y2)
	return input, w, b1, b2
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, w, b1, b2 := genCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		exp := expected(w, b1, b2)
		out = strings.ToUpper(strings.TrimSpace(out))
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
