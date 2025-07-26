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
	out := filepath.Join(os.TempDir(), "1097H_ref")
	cmd := exec.Command("go", "build", "-o", out, "1097H.go")
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

func genTest() string {
	d := rand.Intn(5) + 2
	m := rand.Intn(20) + 2
	gen := make([]string, d)
	gen[0] = "0"
	for i := 1; i < d; i++ {
		gen[i] = fmt.Sprintf("%d", rand.Intn(m))
	}
	n := rand.Intn(5) + 1
	B := make([]string, n)
	for i := range B {
		B[i] = fmt.Sprintf("%d", rand.Intn(m))
	}
	l := rand.Int63n(1000) + 1
	r := l + int64(n) + rand.Int63n(100)
	lines := []string{fmt.Sprintf("%d %d", d, m), strings.Join(gen, " "), fmt.Sprintf("%d", n), strings.Join(B, " "), fmt.Sprintf("%d %d", l, r)}
	return strings.Join(lines, "\n") + "\n"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
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
