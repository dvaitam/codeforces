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
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "562D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	q := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 0; i < q; i++ {
		t := rng.Intn(3) + 1
		x := rng.Intn(n) + 1
		y := rng.Intn(n) + 1
		if t == 2 {
			if x > y {
				x, y = y, x
			}
			if x == y {
				if y < n {
					y++
				} else {
					x--
				}
			}
		}
		sb.WriteString(fmt.Sprint(t))
		sb.WriteString(fmt.Sprintf(" %d %d\n", x, y))
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		expected, err := runBinary(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
