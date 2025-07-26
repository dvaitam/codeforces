package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func buildRef() (string, error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	tmp, err := os.CreateTemp("", "refE-*.bin")
	if err != nil {
		return "", err
	}
	tmp.Close()
	path := tmp.Name()
	cmd := exec.Command("go", "build", "-o", path, filepath.Join(dir, "1016E.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return path, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	y := float64(rng.Intn(10) + 1)
	a := float64(rng.Intn(5))
	b := a + float64(rng.Intn(5)+1)
	n := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%.1f %.1f %.1f\n%d\n", y, a, b, n)
	for i := 0; i < n; i++ {
		l := float64(rng.Intn(10))
		r := l + float64(rng.Intn(5)+1)
		fmt.Fprintf(&sb, "%.1f %.1f\n", l, r)
	}
	m := rng.Intn(5) + 1
	fmt.Fprintf(&sb, "%d\n", m)
	for i := 0; i < m; i++ {
		px := float64(rng.Intn(10))
		py := y + float64(rng.Intn(5)+1)
		fmt.Fprintf(&sb, "%.1f %.1f\n", px, py)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
		exp, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", t+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:\n%s\nGot:\n%s\ninput:\n%s", t+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
