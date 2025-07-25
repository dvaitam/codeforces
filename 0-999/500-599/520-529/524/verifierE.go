package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildRef() (string, error) {
	ref := "refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "524E.go")
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

func genTests() []string {
	rng := rand.New(rand.NewSource(5))
	tests := make([]string, 100)
	for t := 0; t < 100; t++ {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		maxk := n * m
		if maxk > 15 {
			maxk = 15
		}
		k := rng.Intn(maxk) + 1
		q := rng.Intn(5) + 1
		used := map[[2]int]bool{}
		rooks := make([][2]int, k)
		for i := 0; i < k; i++ {
			for {
				x := rng.Intn(n) + 1
				y := rng.Intn(m) + 1
				p := [2]int{x, y}
				if !used[p] {
					used[p] = true
					rooks[i] = p
					break
				}
			}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d\n", n, m, k, q)
		for _, p := range rooks {
			fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
		}
		for i := 0; i < q; i++ {
			x1 := rng.Intn(n) + 1
			x2 := rng.Intn(n) + 1
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			y1 := rng.Intn(m) + 1
			y2 := rng.Intn(m) + 1
			if y1 > y2 {
				y1, y2 = y2, y1
			}
			fmt.Fprintf(&sb, "%d %d %d %d\n", x1, y1, x2, y2)
		}
		tests[t] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
	fmt.Println("OK")
}
