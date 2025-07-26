package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func buildRef() (string, error) {
	ref := "refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "524D.go")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("build ref failed: %v\n%s", err, string(out))
	}
	return ref, nil
}

func runProg(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func formatTime(sec int) string {
	h := sec / 3600
	m := (sec % 3600) / 60
	s := sec % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func genTests() []string {
	rng := rand.New(rand.NewSource(4))
	tests := make([]string, 100)
	for t := 0; t < 100; t++ {
		n := rng.Intn(8) + 2
		M := rng.Intn(n) + 1
		T := rng.Intn(10) + 1
		times := make([]int, n)
		for i := 0; i < n; i++ {
			times[i] = rng.Intn(24 * 3600)
		}
		sort.Ints(times)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, M, T)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%s\n", formatTime(times[i]))
		}
		tests[t] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	target := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(ref)
	tests := genTests()
	for i, tc := range tests {
		expOut, err1 := runProg("./"+ref, tc)
		if err1 != nil {
			fmt.Printf("reference solution runtime error on test %d: %v\n", i+1, err1)
			return
		}
		gotOut, err2 := runProg(target, tc)
		if err2 != nil {
			fmt.Printf("target runtime error on test %d: %v\n", i+1, err2)
			return
		}
		if strings.TrimSpace(expOut) != strings.TrimSpace(gotOut) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, tc, expOut, gotOut)
			return
		}
	}
	fmt.Println("All tests passed")
}
