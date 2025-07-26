package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	input string
}

func runExe(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "984B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

type Grid struct {
	n, m  int
	cells [][]byte
}

func (g Grid) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", g.n, g.m))
	for i := 0; i < g.n; i++ {
		sb.Write(g.cells[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func generateValidGrid(rng *rand.Rand) Grid {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	bombs := make([][]bool, n)
	for i := 0; i < n; i++ {
		bombs[i] = make([]bool, m)
		for j := 0; j < m; j++ {
			bombs[i][j] = rng.Intn(4) == 0
		}
	}
	cells := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if bombs[i][j] {
				row[j] = '*'
			} else {
				count := 0
				for di := -1; di <= 1; di++ {
					for dj := -1; dj <= 1; dj++ {
						if di == 0 && dj == 0 {
							continue
						}
						ni, nj := i+di, j+dj
						if ni >= 0 && ni < n && nj >= 0 && nj < m && bombs[ni][nj] {
							count++
						}
					}
				}
				if count == 0 {
					row[j] = '.'
				} else {
					row[j] = byte('0' + count)
				}
			}
		}
		cells[i] = row
	}
	return Grid{n: n, m: m, cells: cells}
}

func genValidCase(rng *rand.Rand) Test {
	g := generateValidGrid(rng)
	return Test{g.String()}
}

func genInvalidCase(rng *rand.Rand) Test {
	g := generateValidGrid(rng)
	i := rng.Intn(g.n)
	j := rng.Intn(g.m)
	ch := g.cells[i][j]
	options := []byte{'.', '1', '2', '3', '4', '5', '6', '7', '8', '*'}
	newCh := options[rng.Intn(len(options))]
	for newCh == ch {
		newCh = options[rng.Intn(len(options))]
	}
	g.cells[i][j] = newCh
	return Test{g.String()}
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(0))
	tests := make([]Test, 0, 102)
	for i := 0; i < 50; i++ {
		tests = append(tests, genValidCase(rng))
	}
	for i := 0; i < 50; i++ {
		tests = append(tests, genInvalidCase(rng))
	}
	tests = append(tests, Test{"1 1\n.\n"})
	tests = append(tests, Test{"1 1\n*\n"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		exp, err := runExe(ref, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
