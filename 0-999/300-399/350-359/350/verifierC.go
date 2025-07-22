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

type bomb struct {
	x, y, d int
}

type test struct {
	input    string
	expected string
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func solveC(bs []bomb) string {
	// copy algorithm from 350C.go
	// sort by distance
	sort.Slice(bs, func(i, j int) bool { return bs[i].d < bs[j].d })
	var k int
	for _, b := range bs {
		if b.x != 0 {
			k++
		}
		if b.y != 0 {
			k++
		}
		k++ // pick
		if b.y != 0 {
			k++
		}
		if b.x != 0 {
			k++
		}
		k++ // destroy
	}
	var out strings.Builder
	fmt.Fprintf(&out, "%d\n", k)
	for _, b := range bs {
		if b.x > 0 {
			fmt.Fprintf(&out, "1 %d R\n", b.x)
		} else if b.x < 0 {
			fmt.Fprintf(&out, "1 %d L\n", -b.x)
		}
		if b.y > 0 {
			fmt.Fprintf(&out, "1 %d U\n", b.y)
		} else if b.y < 0 {
			fmt.Fprintf(&out, "1 %d D\n", -b.y)
		}
		fmt.Fprintln(&out, 2)
		if b.y > 0 {
			fmt.Fprintf(&out, "1 %d D\n", b.y)
		} else if b.y < 0 {
			fmt.Fprintf(&out, "1 %d U\n", -b.y)
		}
		if b.x > 0 {
			fmt.Fprintf(&out, "1 %d L\n", b.x)
		} else if b.x < 0 {
			fmt.Fprintf(&out, "1 %d R\n", -b.x)
		}
		fmt.Fprintln(&out, 3)
	}
	return out.String()
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(44))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		seen := make(map[[2]int]bool)
		bs := make([]bomb, 0, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for len(bs) < n {
			x := rng.Intn(21) - 10
			y := rng.Intn(21) - 10
			if x == 0 && y == 0 {
				continue
			}
			if seen[[2]int{x, y}] {
				continue
			}
			seen[[2]int{x, y}] = true
			bs = append(bs, bomb{x: x, y: y, d: abs(x) + abs(y)})
			fmt.Fprintf(&sb, "%d %d\n", x, y)
		}
		expected := solveC(append([]bomb(nil), bs...))
		tests = append(tests, test{sb.String(), expected})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
