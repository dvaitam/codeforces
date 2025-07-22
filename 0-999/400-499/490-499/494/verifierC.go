package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func baseDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

func prepareBinary(path, tag string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		bin := filepath.Join(os.TempDir(), tag)
		cmd := exec.Command("go", "build", "-o", bin, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", path, err, out)
		}
		return bin, nil
	}
	return path, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

type interval struct{ l, r int }

func disjointOrNested(a interval, b interval) bool {
	if a.r < b.l || b.r < a.l {
		return true
	}
	if a.l <= b.l && b.r <= a.r {
		return true
	}
	if b.l <= a.l && a.r <= b.r {
		return true
	}
	return false
}

func genTests() []string {
	rand.Seed(1)
	var tests []string
	for t := 0; t < 100; t++ {
		n := rand.Intn(5) + 1
		q := rand.Intn(3) + 1
		arr := make([]int, n)
		for i := range arr {
			arr[i] = rand.Intn(10)
		}
		lines := []string{fmt.Sprintf("%d %d", n, q)}
		var parts []string
		for _, v := range arr {
			parts = append(parts, fmt.Sprintf("%d", v))
		}
		lines = append(lines, strings.Join(parts, " "))
		intervals := []interval{}
		for i := 0; i < q; i++ {
			for {
				l := rand.Intn(n) + 1
				r := l + rand.Intn(n-l+1)
				cand := interval{l, r}
				ok := true
				for _, iv := range intervals {
					if !disjointOrNested(cand, iv) {
						ok = false
						break
					}
				}
				if ok {
					intervals = append(intervals, cand)
					p := float64(rand.Intn(1000)) / 1000.0
					lines = append(lines, fmt.Sprintf("%d %d %.3f", l, r, p))
					break
				}
			}
		}
		tests = append(tests, strings.Join(lines, "\n")+"\n")
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go <binary>")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candC")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refSrc := filepath.Join(baseDir(), "494C.go")
	refPath, err := prepareBinary(refSrc, "refC")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	tests := genTests()
	for i, input := range tests {
		exp, err := runBinary(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
