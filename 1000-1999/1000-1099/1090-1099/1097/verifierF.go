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
	out := filepath.Join(os.TempDir(), "1097F_ref")
	cmd := exec.Command("go", "build", "-o", out, "1097F.go")
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

func genTest() string {
	n := rand.Intn(5) + 1
	q := rand.Intn(30) + 1
	lines := []string{fmt.Sprintf("%d %d", n, q)}
	queryCount := 0
	for i := 0; i < q; i++ {
		typ := rand.Intn(4) + 1
		if i == q-1 && queryCount == 0 {
			typ = 4
		}
		switch typ {
		case 1:
			x := rand.Intn(n) + 1
			v := rand.Intn(50) + 1
			lines = append(lines, fmt.Sprintf("1 %d %d", x, v))
		case 2:
			x := rand.Intn(n) + 1
			y := rand.Intn(n) + 1
			z := rand.Intn(n) + 1
			lines = append(lines, fmt.Sprintf("2 %d %d %d", x, y, z))
		case 3:
			x := rand.Intn(n) + 1
			y := rand.Intn(n) + 1
			z := rand.Intn(n) + 1
			lines = append(lines, fmt.Sprintf("3 %d %d %d", x, y, z))
		case 4:
			x := rand.Intn(n) + 1
			v := rand.Intn(50) + 1
			lines = append(lines, fmt.Sprintf("4 %d %d", x, v))
			queryCount++
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
