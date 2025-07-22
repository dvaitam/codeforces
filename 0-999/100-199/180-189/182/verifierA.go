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

const refDirA = "./0-999/100-199/180-189/182"

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildRef() (string, error) {
	refPath := filepath.Join(refDirA, "refA.bin")
	cmd := exec.Command("go", "build", "-o", refPath, "182A.go")
	cmd.Dir = refDirA
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return refPath, nil
}

func genTest() string {
	a := rand.Intn(5) + 1
	b := rand.Intn(5) + 1
	ax := rand.Intn(21) - 10
	ay := rand.Intn(21) - 10
	bx := rand.Intn(21) - 10
	by := rand.Intn(21) - 10
	n := rand.Intn(3) // 0..2 trenches
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", a, b)
	fmt.Fprintf(&sb, "%d %d %d %d\n", ax, ay, bx, by)
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if rand.Intn(2) == 0 {
			x := rand.Intn(21) - 10
			y := rand.Intn(21) - 10
			l := rand.Intn(b) + 1
			y2 := y + l
			fmt.Fprintf(&sb, "%d %d %d %d\n", x, y, x, y2)
		} else {
			x := rand.Intn(21) - 10
			y := rand.Intn(21) - 10
			l := rand.Intn(b) + 1
			x2 := x + l
			fmt.Fprintf(&sb, "%d %d %d %d\n", x, y, x2, y)
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
