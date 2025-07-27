package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveC(s string) string {
	a := []byte(s)
	n := len(a)
	changes := 0
	for i := 0; i < n; i++ {
		if (i > 0 && a[i] == a[i-1]) || (i > 1 && a[i] == a[i-2]) {
			changes++
			a[i] = '?' // mark
		}
	}
	return fmt.Sprint(changes)
}

func genCases() []string {
	rand.Seed(3)
	cases := make([]string, 100)
	letters := []byte("abc")
	for i := 0; i < 100; i++ {
		l := rand.Intn(10) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = letters[rand.Intn(len(letters))]
		}
		s := string(b)
		tc := fmt.Sprintf("1\n%s\n", s)
		cases[i] = tc
	}
	return cases
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		s := strings.TrimSpace(strings.Split(tc, "\n")[1])
		want := solveC(s)
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected: %s Got: %s\n", i+1, tc, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
