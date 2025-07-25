package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runProgram(bin string, input string) (string, error) {
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

func digits(x int) int {
	d := 0
	for x > 0 {
		d++
		x /= 10
	}
	if d == 0 {
		d = 1
	}
	return d
}

func randomEdges(rng *rand.Rand, n int) [][2]string {
	D := digits(n)
	edges := make([][2]string, n-1)
	for i := 0; i < n-1; i++ {
		a := rng.Intn(D) + 1
		b := rng.Intn(D) + 1
		edges[i] = [2]string{strings.Repeat("?", a), strings.Repeat("?", b)}
	}
	return edges
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	refBin := "ref611H.bin"
	if err := exec.Command("go", "build", "-o", refBin, "611H.go").Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(8))
	for t := 0; t < 100; t++ {
		n := rng.Intn(8) + 2
		edges := randomEdges(rng, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%s %s\n", e[0], e[1])
		}
		input := sb.String()
		want, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
			os.Exit(1)
		}
		out, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\n=== got ===\n%s\n", t+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
