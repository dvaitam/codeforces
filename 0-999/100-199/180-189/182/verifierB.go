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

const refDirB = "./0-999/100-199/180-189/182"

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildRef() (string, error) {
	refPath := filepath.Join(refDirB, "refB.bin")
	cmd := exec.Command("go", "build", "-o", refPath, "182B.go")
	cmd.Dir = refDirB
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return refPath, nil
}

func genTest() string {
	d := rand.Intn(50) + 1
	n := rand.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", d)
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rand.Intn(d)+1)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	candidate := os.Args[1]
	rand.Seed(1)
	ref, err := buildRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	total := 100
	passed := 0
	for i := 0; i < total; i++ {
		input := genTest()
		exp, err := runBinary(ref, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			return
		}
		out, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			continue
		}
		if out == exp {
			passed++
		} else {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\n got:%s\n", i+1, input, exp, out)
		}
	}
	fmt.Printf("Passed %d/%d tests\n", passed, total)
}
