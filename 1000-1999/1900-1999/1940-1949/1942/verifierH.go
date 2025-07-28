package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildRef() (string, error) {
	tmp := filepath.Join(os.TempDir(), "refH_1942")
	cmd := exec.Command("go", "build", "-o", tmp, "1942H.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return tmp, nil
}

func runProg(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTest(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	q := rng.Intn(4) + 1
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "1\n%d %d\n", n, q)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		if i > 2 {
			buf.WriteByte(' ')
		}
		fmt.Fprint(&buf, p)
	}
	if n > 1 {
		buf.WriteByte('\n')
	} else {
		buf.WriteByte('\n')
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprint(&buf, rng.Intn(5))
	}
	buf.WriteByte('\n')
	for i := 0; i < q; i++ {
		t := rng.Intn(2) + 1
		x := rng.Intn(n) + 1
		v := rng.Intn(5) + 1
		fmt.Fprintf(&buf, "%d %d %d\n", t, x, v)
	}
	return buf.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(50))
	for i := 0; i < 100; i++ {
		test := genTest(rng)
		expected, err := runProg(ref, test)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference error:", err)
			os.Exit(1)
		}
		got, err := runProg(target, test)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d execution error: %v\n", i+1, err)
			os.Exit(1)
		}
		if expected != got {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, test, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
