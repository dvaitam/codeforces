package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type point struct {
	x, w int
	id   int
}

func buildExecutable(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "bin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

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

func oracle(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(reader, &t)
	var n, m int
	fmt.Fscan(reader, &n, &m)
	pts := make([]point, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &pts[i].x, &pts[i].w)
		pts[i].id = i + 1
	}
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].w != pts[j].w {
			return pts[i].w < pts[j].w
		}
		return pts[i].x < pts[j].x
	})
	selected := pts[:2*n]
	sum := 0
	for _, p := range selected {
		sum += p.w
	}
	sort.Slice(selected, func(i, j int) bool {
		return selected[i].x < selected[j].x
	})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", sum))
	l, r := 0, len(selected)-1
	for l < r {
		sb.WriteString(fmt.Sprintf("%d %d\n", selected[l].id, selected[r].id))
		l++
		r--
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := 2*n + rng.Intn(3)
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, m)
	usedX := map[int]bool{}
	pts := make([]point, m)
	for i := 0; i < m; i++ {
		x := rng.Intn(1000) - 500
		for usedX[x] {
			x = rng.Intn(1000) - 500
		}
		usedX[x] = true
		w := rng.Intn(20) - 10
		fmt.Fprintf(&sb, "%d %d\n", x, w)
		pts[i] = point{x: x, w: w, id: i + 1}
	}
	input := sb.String()
	out := oracle(input)
	return input, out
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binPath := os.Args[1]
	bin, cleanup, err := buildExecutable(binPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()
	rng := rand.New(rand.NewSource(3))
	for i := 0; i < 100; i++ {
		input, expected := genCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
