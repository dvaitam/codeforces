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
	ref := "refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "524C.go")
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
	rng := rand.New(rand.NewSource(3))
	tests := make([]string, 100)
	for t := 0; t < 100; t++ {
		n := rng.Intn(4) + 1
		k := rng.Intn(20) + 1
		den := make([]int, n)
		used := map[int]bool{}
		for i := 0; i < n; i++ {
			for {
				v := rng.Intn(50) + 1
				if !used[v] {
					used[v] = true
					den[i] = v
					break
				}
			}
		}
		sort.Ints(den)
		q := rng.Intn(5) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for _, v := range den {
			fmt.Fprintf(&sb, "%d ", v)
		}
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "%d\n", q)
		for i := 0; i < q; i++ {
			x := rng.Intn(200) + 1
			fmt.Fprintf(&sb, "%d\n", x)
		}
		tests[t] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
