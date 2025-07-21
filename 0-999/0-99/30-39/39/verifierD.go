package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveD(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var x1, y1, z1 int
	var x2, y2, z2 int
	if _, err := fmt.Fscan(r, &x1, &y1, &z1, &x2, &y2, &z2); err != nil {
		return ""
	}
	if x1 == x2 || y1 == y2 || z1 == z2 {
		return "YES\n"
	}
	return "NO\n"
}

func generateCaseD(rng *rand.Rand) string {
	a := make([]int, 6)
	for i := range a {
		a[i] = rng.Intn(2)
	}
	return fmt.Sprintf("%d %d %d %d %d %d\n", a[0], a[1], a[2], a[3], a[4], a[5])
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseD(rng)
	}
	for i, tc := range cases {
		expect := solveD(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
