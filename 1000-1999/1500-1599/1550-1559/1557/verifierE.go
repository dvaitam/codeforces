package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveE(t int) string {
	var sb strings.Builder
	for i := 0; i < t; i++ {
		for j := 1; j <= 64 && j <= 130; j++ {
			x := (j-1)%8 + 1
			y := (j-1)%8 + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		}
	}
	return strings.TrimSpace(sb.String())
}

func runCase(bin string, t int) (string, error) {
	input := fmt.Sprintf("%d\n", t)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]int, 0, 102)
	cases = append(cases, 1)
	cases = append(cases, 2)
	for i := 0; i < 100; i++ {
		cases = append(cases, rng.Intn(60)+1)
	}
	for i, tcase := range cases {
		expect := solveE(tcase)
		out, err := runCase(bin, tcase)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\n\nGot:\n%s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
