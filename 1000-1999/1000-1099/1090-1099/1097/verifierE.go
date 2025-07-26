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
	out := filepath.Join(os.TempDir(), "1097E_ref")
	cmd := exec.Command("go", "build", "-o", out, "1097E.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

func runBin(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func perm(n int) []int {
	a := make([]int, n)
	for i := range a {
		a[i] = i + 1
	}
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
	return a
}

func genTest() string {
	q := rand.Intn(3) + 1
	lines := []string{fmt.Sprintf("%d", q)}
	for i := 0; i < q; i++ {
		n := rand.Intn(10) + 1
		a := perm(n)
		line := fmt.Sprintf("%d\n", n)
		parts := make([]string, n)
		for j, v := range a {
			parts[j] = fmt.Sprintf("%d", v)
		}
		line += strings.Join(parts, " ")
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n") + "\n"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
