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
type Rect struct{ x1, y1, x2, y2 int }

func buildRef() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("read reference: %v", err)
	}
	ref := "./refG.bin"
	if strings.Contains(string(data), "#include") {
		cppPath := "refG.cpp"
		if err := os.WriteFile(cppPath, data, 0644); err != nil {
			return "", fmt.Errorf("write cpp: %v", err)
		}
		cmd := exec.Command("g++", "-O2", "-o", ref, cppPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build reference cpp: %v: %s", err, string(out))
		}
	} else {
		cmd := exec.Command("go", "build", "-o", ref, src)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build reference: %v: %s", err, string(out))
		}
	}
	return ref, nil
}

func genCase(rng *rand.Rand) string {
	x1 := rng.Intn(20)
	x2 := x1 + rng.Intn(5) + 1
	y2 := rng.Intn(20)
	y1 := y2 + rng.Intn(5) + 1
	// polygon as rectangle
	px1 := rng.Intn(20)
	py2 := rng.Intn(20)
	px2 := px1 + rng.Intn(5) + 1
	py1 := py2 + rng.Intn(5) + 1
	poly := []Point{{px1, py1}, {px2, py1}, {px2, py2}, {px1, py2}}
	n := len(poly)
	input := fmt.Sprintf("%d %d %d %d\n%d\n", x1, y1, x2, y2, n)
	for _, p := range poly {
		input += fmt.Sprintf("%d %d\n", p.x, p.y)
	}
	return input
}

func runBin(path, input string) (string, error) {
	cmd := exec.Command(path)
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
		fmt.Println("usage: verifierG /path/to/binary")
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
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp, err := runBin(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: reference error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
