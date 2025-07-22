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
	cmd := exec.Command("go", "build", "-o", ref, "249D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		a := rng.Intn(5)
		b := rng.Intn(5) + 1
		c := a + rng.Intn(5) + 1
		d := b
		input := fmt.Sprintf("%d\n%d/%d %d/%d\n", n, a, b, c, d)
		for j := 0; j < n; j++ {
			x := rng.Intn(50) + 1
			y := rng.Intn(50) + 1
			input += fmt.Sprintf("%d %d\n", x, y)
		}
		expect, err := runBinary(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(expect) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "test %d failed:\nexpected: %s\n got: %s\ninput:\n%s", i+1, strings.TrimSpace(expect), strings.TrimSpace(got), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
