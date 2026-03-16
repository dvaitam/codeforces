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

func buildReference() (string, error) {
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	refBin := "ref_solution"
	cmd := exec.Command("go", "build", "-o", refBin, refSrc)
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, errBuf.String())
	}
	return refBin, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCaseG(rng *rand.Rand) string {
	typ := rng.Intn(3)
	switch typ {
	case 0:
		target := rng.Intn(32768)
		return fmt.Sprintf("%d\nint f(int n){\nreturn n;\n}\n", target)
	case 1:
		c := rng.Intn(10) + 1
		target := rng.Intn(32768)
		code := fmt.Sprintf("int f(int n){\nif (n>%d) return n-%d;\nreturn n;\n}\n", c, c)
		return fmt.Sprintf("%d\n%s", target, code)
	default:
		c := rng.Intn(10) + 1
		target := rng.Intn(32768)
		code := fmt.Sprintf("int f(int n){\nif (n<%d) return n;\nreturn n-%d;\n}\n", c, c)
		return fmt.Sprintf("%d\n%s", target, code)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierG /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "build reference: %v\n", err)
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseG(rng)

		expect, err := runBinary("./"+refBin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: reference error: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}

		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate error: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}

		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
