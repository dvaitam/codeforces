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

func buildOracle() (string, error) {
	exe := "oracleE"
	cmd := exec.Command("go", "build", "-o", exe, "./0-999/0-99/80-89/87/87E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func genTriangle(rng *rand.Rand) []point {
	x0 := rng.Intn(10)
	y0 := rng.Intn(10)
	dx1 := rng.Intn(5) + 1
	dx2 := rng.Intn(dx1)
	dy := rng.Intn(5) + 1
	p1 := point{x0, y0}
	p2 := point{x0 + dx1, y0}
	p3 := point{x0 + dx2, y0 + dy}
	return []point{p1, p2, p3}
}

func generateCase(rng *rand.Rand) string {
	var sb strings.Builder
	for i := 0; i < 3; i++ {
		pts := genTriangle(rng)
		fmt.Fprintf(&sb, "%d\n", len(pts))
		for _, p := range pts {
			fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
		}
	}
	m := rng.Intn(5) + 1
	fmt.Fprintf(&sb, "%d\n", m)
	for i := 0; i < m; i++ {
		x := rng.Intn(15)
		y := rng.Intn(15)
		fmt.Fprintf(&sb, "%d %d\n", x, y)
	}
	return sb.String()
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := runProg("./"+oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
