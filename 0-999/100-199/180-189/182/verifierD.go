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

const refDirD = "./0-999/100-199/180-189/182"

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildRef() (string, error) {
	refPath := filepath.Join(refDirD, "refD.bin")
	cmd := exec.Command("go", "build", "-o", refPath, "182D.go")
	cmd.Dir = refDirD
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return refPath, nil
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rand.Intn(3))
	}
	return string(b)
}

func genTest() string {
	l1 := rand.Intn(8) + 1
	l2 := rand.Intn(8) + 1
	s1 := randString(l1)
	s2 := randString(l2)
	return fmt.Sprintf("%s %s\n", s1, s2)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
			fmt.Printf("test %d failed\ninput:%sexpected:%s\n got:%s\n", i+1, input, exp, out)
		}
	}
	fmt.Printf("Passed %d/%d tests\n", passed, total)
}
