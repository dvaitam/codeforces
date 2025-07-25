package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func rotate(mat [][]int) [][]int {
	n := len(mat)
	res := make([][]int, n)
	for i := range res {
		res[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			res[i][j] = mat[n-1-j][i]
		}
	}
	return res
}

func matToString(mat [][]int) string {
	var sb strings.Builder
	for i, row := range mat {
		for j, v := range row {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		if i+1 < len(mat) {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	for t := 0; t < 100; t++ {
		n := rand.Intn(10) + 1
		mat := make([][]int, n)
		for i := range mat {
			mat[i] = make([]int, n)
			for j := range mat[i] {
				mat[i][j] = rand.Intn(100) + 1
			}
		}
		in := matToString(mat) + "\n"
		want := matToString(rotate(mat))
		out, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != want {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", t+1, strings.TrimSpace(in), want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
