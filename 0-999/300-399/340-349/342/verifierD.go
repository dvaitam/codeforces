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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 1
	grid := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			row[j] = '.'
		}
		grid[i] = row
	}
	or := rng.Intn(3)
	oc := rng.Intn(n)
	grid[or][oc] = 'O'
	for i := 0; i < 3; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 'O' {
				continue
			}
			if rng.Intn(3) == 0 {
				grid[i][j] = 'X'
			}
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < 3; i++ {
		b.WriteString(string(grid[i]))
		b.WriteByte('\n')
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	ref := filepath.Join(dir, "refD")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "342D.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		in := generateCase(rng)
		candOut, cErr := runBinary(cand, in)
		if cErr != nil {
			fmt.Printf("test %d: candidate error: %v\n", t+1, cErr)
			os.Exit(1)
		}
		refOut, rErr := runBinary(ref, in)
		if rErr != nil {
			fmt.Printf("test %d: reference error: %v\n", t+1, rErr)
			os.Exit(1)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%sactual:%s\n", t+1, in, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
