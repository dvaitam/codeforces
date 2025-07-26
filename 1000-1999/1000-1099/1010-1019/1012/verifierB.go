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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierB.go /path/to/binary\n")
		os.Exit(1)
	}
	candidate := os.Args[1]
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	ref := filepath.Join(dir, "refB")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "1012B.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(1)
	for t := 0; t < 100; t++ {
		n := rand.Intn(10) + 1
		m := rand.Intn(10) + 1
		maxQ := n * m
		if maxQ > 50 {
			maxQ = 50
		}
		q := rand.Intn(maxQ + 1)
		used := make(map[[2]int]bool)
		var b bytes.Buffer
		fmt.Fprintf(&b, "%d %d %d\n", n, m, q)
		for i := 0; i < q; i++ {
			var r, c int
			for {
				r = rand.Intn(n) + 1
				c = rand.Intn(m) + 1
				if !used[[2]int{r, c}] {
					used[[2]int{r, c}] = true
					break
				}
			}
			fmt.Fprintf(&b, "%d %d\n", r, c)
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
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%s\nactual:%s\n", t+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
