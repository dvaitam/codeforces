package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func compileRef() (string, error) {
	out := filepath.Join(os.TempDir(), "1097G_ref")
	cmd := exec.Command("go", "build", "-o", out, "1097G.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

func runBin(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genTree(n int) [][2]int {
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	return edges
}

func genTest() string {
	n := rand.Intn(6) + 2
	k := rand.Intn(n) + 1
	edges := genTree(n)
	lines := []string{fmt.Sprintf("%d %d", n, k)}
	for _, e := range edges {
		lines = append(lines, fmt.Sprintf("%d %d", e[0], e[1]))
	}
	return strings.Join(lines, "\n") + "\n"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to compile reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(1)
	for t := 0; t < 100; t++ {
		test := genTest()
		exp, err := runBin(ref, test)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference run failed:", err)
			os.Exit(1)
		}
		got, err := runBin(cand, test)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate run failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("mismatch on test %d\ninput:\n%sexpected:%s got:%s\n", t+1, test, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
