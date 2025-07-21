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
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierC.go /path/to/binary\n")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	ref := filepath.Join(dir, "refC")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "83C.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(3)
	for t := 0; t < 100; t++ {
		n := rand.Intn(4) + 2
		m := rand.Intn(4) + 2
		k := rand.Intn(4) + 1
		grid := make([][]byte, n)
		for i := range grid {
			grid[i] = make([]byte, m)
			for j := range grid[i] {
				grid[i][j] = byte('a' + rand.Intn(4))
			}
		}
		sr, sc := rand.Intn(n), rand.Intn(m)
		tr, tc := rand.Intn(n), rand.Intn(m)
		for tr == sr && tc == sc {
			tr, tc = rand.Intn(n), rand.Intn(m)
		}
		grid[sr][sc] = 'S'
		grid[tr][tc] = 'T'

		var b bytes.Buffer
		fmt.Fprintf(&b, "%d %d %d\n", n, m, k)
		for i := 0; i < n; i++ {
			b.Write(grid[i])
			b.WriteByte('\n')
		}
		input := b.String()

		candOut, cErr := runBinary(candidate, input)
		refOut, rErr := runBinary(ref, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "test %d: candidate error: %v\n", t+1, cErr)
			os.Exit(1)
		}
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "test %d: reference error: %v\n", t+1, rErr)
			os.Exit(1)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%sactual:%s\n", t+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}
