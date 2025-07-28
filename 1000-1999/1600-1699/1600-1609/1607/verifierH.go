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

func randCase() (int, []int, []int, []int) {
	n := rand.Intn(5) + 1
	a := make([]int, n)
	b := make([]int, n)
	m := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rand.Intn(10)
		b[i] = rand.Intn(10)
		limit := a[i] + b[i]
		m[i] = rand.Intn(limit + 1)
	}
	return n, a, b, m
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	refBin := "./refH.bin"
	if err := exec.Command("go", "build", "-o", refBin, "1607H.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rand.Seed(8)
	type tc struct {
		n       int
		a, b, m []int
	}
	cases := []tc{{1, []int{0}, []int{0}, []int{0}}}
	for len(cases) < 100 {
		n, a, b, m := randCase()
		cases = append(cases, tc{n, a, b, m})
	}

	for i, c := range cases {
		input := fmt.Sprintf("1\n%d\n", c.n)
		for j := 0; j < c.n; j++ {
			input += fmt.Sprintf("%d %d %d\n", c.a[j], c.b[j], c.m[j])
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
