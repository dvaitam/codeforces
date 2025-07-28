package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runCmd(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = io.Discard
	err := cmd.Run()
	return out.String(), err
}

func randCase() (int, int64, []int64, []int64) {
	n := rand.Intn(5) + 1
	m := int64(rand.Intn(10))
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(rand.Intn(10))
		b[i] = int64(rand.Intn(10))
		if a[i]+b[i] < m {
			a[i] += m - (a[i] + b[i])
		}
	}
	return n, m, a, b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	refBin := "./refG.bin"
	if err := exec.Command("go", "build", "-o", refBin, "1607G.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rand.Seed(7)
	type tc struct {
		n    int
		m    int64
		a, b []int64
	}
	cases := []tc{{1, 0, []int64{0}, []int64{0}}}
	for len(cases) < 100 {
		n, m, a, b := randCase()
		cases = append(cases, tc{n, m, a, b})
	}

	for i, c := range cases {
		input := fmt.Sprintf("1\n%d %d\n", c.n, c.m)
		for j := 0; j < c.n; j++ {
			input += fmt.Sprintf("%d %d\n", c.a[j], c.b[j])
		}
		exp, err := runCmd(refBin, input)
		if err != nil {
			fmt.Println("reference solution failed:", err)
			os.Exit(1)
		}
		got, err := runCmd(candidate, input)
		if err != nil {
			fmt.Printf("test %d: candidate runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		exp = strings.TrimSpace(exp)
		got = strings.TrimSpace(got)
		if exp != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
