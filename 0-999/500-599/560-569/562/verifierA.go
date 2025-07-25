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
	ref := "./refA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "562A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(rng.Intn(6)))
	}
	sb.WriteByte('\n')
	for i := 2; i <= n; i++ {
		v := rng.Intn(i-1) + 1
		l := rng.Intn(5) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", i, v, l))
	}
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func parse(out string) (int, float64, error) {
	var f int
	var c float64
	_, err := fmt.Sscan(out, &f, &c)
	return f, c, err
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		expStr, err := runBinary(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr, err := runBinary(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		ef, ec, err := parse(expStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad reference output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gf, gc, err := parse(gotStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if ef != gf || abs(ec-gc) > 1e-6 {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %d %.10f\ngot: %d %.10f\n", i+1, in, ef, ec, gf, gc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
