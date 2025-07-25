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

func solveF(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(in, &n)
	rows := make([]int, n)
	for i := 0; i < n; i++ {
		var r, c int
		fmt.Fscan(in, &r, &c)
		rows[r-1] = c
	}
	ans := 0
	for i := 0; i < n; i++ {
		minVal := rows[i]
		maxVal := rows[i]
		for j := i; j < n; j++ {
			if rows[j] < minVal {
				minVal = rows[j]
			}
			if rows[j] > maxVal {
				maxVal = rows[j]
			}
			if maxVal-minVal == j-i {
				ans++
			}
		}
	}
	return fmt.Sprint(ans)
}

func genTestF(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	cols := rng.Perm(n)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&buf, "%d %d\n", i+1, cols[i]+1)
	}
	return buf.String()
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestF(rng)
		expect := solveF(in)
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
