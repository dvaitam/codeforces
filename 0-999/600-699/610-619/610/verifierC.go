package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveC(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var k int
	if _, err := fmt.Fscan(in, &k); err != nil {
		return ""
	}
	n := 1 << uint(k)
	mat := [][]byte{[]byte{'+'}}
	size := 1
	for size < n {
		newSize := size * 2
		newMat := make([][]byte, newSize)
		for i := range newMat {
			newMat[i] = make([]byte, newSize)
		}
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				ch := mat[i][j]
				newMat[i][j] = ch
				newMat[i][j+size] = ch
				newMat[i+size][j] = ch
				if ch == '+' {
					newMat[i+size][j+size] = '*'
				} else {
					newMat[i+size][j+size] = '+'
				}
			}
		}
		mat = newMat
		size = newSize
	}
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		buf.Write(mat[i])
		buf.WriteByte('\n')
	}
	return strings.TrimSpace(buf.String())
}

func genTestC(rng *rand.Rand) string {
	k := rng.Intn(5) // 0..4
	return fmt.Sprintf("%d\n", k)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestC(rng)
		expect := solveC(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
