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

type Sofa struct{ x1, y1, x2, y2 int }

type Test struct {
	d, n, m                int
	sofas                  []Sofa
	cntl, cntr, cntt, cntb int
}

func buildRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "818C.go")
	bin := filepath.Join(os.TempDir(), "ref818C.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return bin, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []Test {
	rand.Seed(3)
	tests := make([]Test, 100)
	for idx := 0; idx < 100; idx++ {
		n := rand.Intn(4) + 2
		m := rand.Intn(4) + 2
		maxPairs := n * m / 2
		d := rand.Intn(maxPairs) + 1
		// generate positions
		type Pair [4]int
		var pairs []Pair
		for x := 1; x <= n; x++ {
			for y := 1; y <= m; y++ {
				if x < n {
					pairs = append(pairs, Pair{x, y, x + 1, y})
				}
				if y < m {
					pairs = append(pairs, Pair{x, y, x, y + 1})
				}
			}
		}
		rand.Shuffle(len(pairs), func(i, j int) { pairs[i], pairs[j] = pairs[j], pairs[i] })
		used := make([][]bool, n+1)
		for i := range used {
			used[i] = make([]bool, m+1)
		}
		sofas := make([]Sofa, 0, d)
		for _, p := range pairs {
			if len(sofas) == d {
				break
			}
			x1, y1, x2, y2 := p[0], p[1], p[2], p[3]
			if used[x1][y1] || used[x2][y2] {
				continue
			}
			used[x1][y1] = true
			used[x2][y2] = true
			sofas = append(sofas, Sofa{x1, y1, x2, y2})
		}
		d = len(sofas)
		cntl := rand.Intn(d + 1)
		cntr := rand.Intn(d + 1)
		cntt := rand.Intn(d + 1)
		cntb := rand.Intn(d + 1)
		tests[idx] = Test{d, n, m, sofas, cntl, cntr, cntt, cntb}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	tests := genTests()
	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t.d))
		sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.m))
		for _, s := range t.sofas {
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", s.x1, s.y1, s.x2, s.y2))
		}
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", t.cntl, t.cntr, t.cntt, t.cntb))
		input := sb.String()
		want, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:%sexpected:%s\nactual:%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
