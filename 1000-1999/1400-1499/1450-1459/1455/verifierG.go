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
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	ref := filepath.Join(dir, "refG.bin")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "1455G.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func genProgram(rng *rand.Rand) (int, int, []string) {
	n := rng.Intn(8) + 1
	s := rng.Intn(10) + 1
	lines := make([]string, 0, n)
	open := 0
	for len(lines) < n {
		if open > 0 && rng.Intn(4) == 0 && len(lines) < n { // close block
			lines = append(lines, "end")
			open--
			continue
		}
		typ := rng.Intn(2)
		if typ == 0 && open < 2 {
			y := rng.Intn(10)
			lines = append(lines, fmt.Sprintf("if %d", y))
			open++
		} else {
			y := rng.Intn(10)
			v := rng.Intn(10) + 1
			lines = append(lines, fmt.Sprintf("set %d %d", y, v))
		}
	}
	for open > 0 {
		lines = append(lines, "end")
		open--
	}
	return len(lines), s, lines
}

func buildInput(n, s int, lines []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, s))
	for _, ln := range lines {
		sb.WriteString(ln)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, s, lines := genProgram(rng)
		input := buildInput(n, s, lines)
		exp, err := runBinary(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(exe, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
