package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type edge struct{ u, v int }

func buildRef() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("read reference: %v", err)
	}
	ref := "./refE.bin"
	if strings.Contains(string(data), "#include") {
		cppPath := "refE.cpp"
		if err := os.WriteFile(cppPath, data, 0644); err != nil {
			return "", fmt.Errorf("write cpp: %v", err)
		}
		cmd := exec.Command("g++", "-O2", "-o", ref, cppPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build reference cpp: %v: %s", err, string(out))
		}
	} else {
		cmd := exec.Command("go", "build", "-o", ref, src)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build reference: %v: %s", err, string(out))
		}
	}
	return ref, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func genTestE() (int, int, []edge) {
	n := rand.Intn(10) + 1
	k := rand.Intn(min(20, n-1) + 1)
	edges := make([]edge, 0, n-1)
	for i := 1; i < n; i++ {
		p := rand.Intn(i)
		edges = append(edges, edge{p, i})
	}
	return n, k, edges
}

func runBinary(path string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	path := os.Args[1]

	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	for i := 0; i < 100; i++ {
		n, k, edges := genTestE()
		var b strings.Builder
		fmt.Fprintf(&b, "%d %d\n", n, k)
		for _, e := range edges {
			fmt.Fprintf(&b, "%d %d\n", e.u+1, e.v+1)
		}
		input := b.String()
		expStr, err := runBinary(ref, input)
		if err != nil {
			fmt.Printf("test %d: reference error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		gotStr, err := runBinary(path, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(gotStr) != strings.TrimSpace(expStr) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, expStr, gotStr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
