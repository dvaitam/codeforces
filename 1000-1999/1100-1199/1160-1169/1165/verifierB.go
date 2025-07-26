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

const numTestsB = 100

func solveB(a []int) int {
	sort.Ints(a)
	days := 0
	for _, v := range a {
		if v >= days+1 {
			days++
		}
	}
	return days
}

func run(binary, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(2)
	for t := 1; t <= numTestsB; t++ {
		n := rand.Intn(50) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rand.Intn(100) + 1
		}
		input := fmt.Sprintf("%d\n", n)
		for i, v := range a {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		expect := solveB(append([]int(nil), a...))
		out, err := run(binary, input)
		if err != nil {
			fmt.Printf("test %d failed to run: %v\noutput:%s\n", t, err, out)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) == 0 {
			fmt.Printf("test %d: no output\n", t)
			os.Exit(1)
		}
		var got int
		fmt.Sscanf(fields[0], "%d", &got)
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\ninput:%s\n", t, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
