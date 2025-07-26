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

func buildRef() (string, error) {
	ref := "./refG.bin"
	cmd := exec.Command("go", "build", "-o", ref, "825G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func genTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		a, b := p, i
		if a > b {
			a, b = b, a
		}
		edges[i-2] = [2]int{a, b}
	}
	return edges
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 3
	q := rng.Intn(20) + 1
	edges := genTree(rng, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	// first query type 1
	sb.WriteString(fmt.Sprintf("1 %d\n", rng.Intn(n)))
	for i := 1; i < q; i++ {
		t := rng.Intn(2) + 1
		z := rng.Intn(n)
		sb.WriteString(fmt.Sprintf("%d %d\n", t, z))
	}
	return sb.String()
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
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
		in := genCase(rng)
		want, err := run(ref, in)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", in)
			os.Exit(1)
		}
		got, err := run(exe, in)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", in)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\n", i+1)
			fmt.Println("input:\n", in)
			fmt.Println("expected:\n", want)
			fmt.Println("got:\n", got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
