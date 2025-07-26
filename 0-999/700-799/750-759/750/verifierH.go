package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	input string
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildRef() (string, error) {
	ref := "./refH.bin"
	cmd := exec.Command("go", "build", "-o", ref, "750H.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genGrid(h, w int) string {
	var sb strings.Builder
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if rand.Intn(4) == 0 {
				sb.WriteByte('#')
			} else {
				sb.WriteByte('.')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func genTests() []Test {
	rand.Seed(7)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		h := rand.Intn(4) + 2
		w := rand.Intn(4) + 2
		q := rand.Intn(3) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", h, w, q)
		sb.WriteString(genGrid(h, w))
		for j := 0; j < q; j++ {
			k := rand.Intn(3)
			fmt.Fprintf(&sb, "%d\n", k)
			for t := 0; t < k; t++ {
				r := rand.Intn(h) + 1
				c := rand.Intn(w) + 1
				fmt.Fprintf(&sb, "%d %d\n", r, c)
			}
		}
		tests = append(tests, Test{sb.String()})
	}
	tests = append(tests, Test{"2 2 1\n..\n..\n0\n"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierH.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		exp, err := runExe(ref, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
