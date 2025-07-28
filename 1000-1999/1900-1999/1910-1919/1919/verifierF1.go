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
	ref := "refF1.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1919F1.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runExe(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	q := rng.Intn(3) + 1
	a := make([]int, n)
	b := make([]int, n)
	c := make([]int, n-1)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(5)
		b[i] = rng.Intn(5)
	}
	for i := 0; i < n-1; i++ {
		c[i] = rng.Intn(5)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d ", a[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d ", b[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < n-1; i++ {
		fmt.Fprintf(&sb, "%d ", c[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		p := rng.Intn(n) + 1
		x := rng.Intn(5)
		y := rng.Intn(5)
		z := rng.Intn(5)
		fmt.Fprintf(&sb, "%d %d %d %d\n", p, x, y, z)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
