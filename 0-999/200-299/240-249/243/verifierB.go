package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func compileRef(src string) (string, error) {
	tmp, err := os.CreateTemp("", "refB-")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())
	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return tmp.Name(), nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTest() string {
	n := rand.Intn(8) + 2
	maxEdges := n * (n - 1) / 2
	m := rand.Intn(maxEdges + 1)
	h := rand.Intn(3) + 1
	t := rand.Intn(3) + 1
	edges := make(map[[2]int]struct{})
	for len(edges) < m {
		u := rand.Intn(n) + 1
		v := rand.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		edges[[2]int{u, v}] = struct{}{}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d %d\n", n, len(edges), h, t)
	for e := range edges {
		fmt.Fprintf(&b, "%d %d\n", e[0], e[1])
	}
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate, _ := filepath.Abs(os.Args[1])
	rand.Seed(time.Now().UnixNano())
	refBin, err := compileRef("243B.go")
	if err != nil {
		fmt.Println("failed to compile reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)
	for i := 0; i < 100; i++ {
		input := genTest()
		expOut, err := runBinary(refBin, input)
		if err != nil {
			fmt.Println("reference failed:", err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(expOut)
		actOut, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		actual := strings.TrimSpace(actOut)
		if actual != expected {
			fmt.Printf("test %d failed:\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, input, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
