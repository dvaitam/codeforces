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

type point struct{ x, y int }

func run(bin string, in []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1184C1.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genCase(rng *rand.Rand) []byte {
	n := rng.Intn(9) + 2 // 2..10
	x0 := rng.Intn(40)
	y0 := rng.Intn(40)
	L := rng.Intn(10) + n
	if x0+L > 50 {
		x0 = 50 - L
	}
	if y0+L > 50 {
		y0 = 50 - L
	}
	pts := make([]point, 0, 4*n+1)
	seen := map[point]bool{}
	add := func(p point) {
		if !seen[p] {
			seen[p] = true
			pts = append(pts, p)
		}
	}
	// left
	for len(pts) < n {
		add(point{x0, y0 + rng.Intn(L+1)})
	}
	// right
	for len(pts) < 2*n {
		add(point{x0 + L, y0 + rng.Intn(L+1)})
	}
	// bottom
	for len(pts) < 3*n {
		add(point{x0 + rng.Intn(L+1), y0})
	}
	// top
	for len(pts) < 4*n {
		add(point{x0 + rng.Intn(L+1), y0 + L})
	}
	// interior
	ix := x0 + 1 + rng.Intn(L-1)
	iy := y0 + 1 + rng.Intn(L-1)
	add(point{ix, iy})
	rng.Shuffle(len(pts), func(i, j int) { pts[i], pts[j] = pts[j], pts[i] })

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, p := range pts {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := genCase(rng)
		want, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", t+1, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", t+1, string(input), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
