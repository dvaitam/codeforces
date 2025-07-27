package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func expected(n, k int) string {
	if n >= k*k && n%2 == k%2 {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(1)
	t := 100
	cases := make([][2]int, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(1000000) + 1
		k := rand.Intn(1000) + 1
		cases[i] = [2]int{n, k}
	}
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", t))
	for _, c := range cases {
		input.WriteString(fmt.Sprintf("%d %d\n", c[0], c[1]))
	}
	in := input.String()

	var expectedOut strings.Builder
	for _, c := range cases {
		expectedOut.WriteString(expected(c[0], c[1]))
		expectedOut.WriteByte('\n')
	}
	want := strings.TrimSpace(expectedOut.String())

	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("Runtime error: %v\n%s", err, out.String())
		os.Exit(1)
	}

	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	wantLines := strings.Split(want, "\n")
	if len(gotLines) != len(wantLines) {
		fmt.Println("Wrong answer: line count mismatch")
		fmt.Println("Expected:")
		fmt.Println(want)
		fmt.Println("Got:")
		fmt.Println(out.String())
		os.Exit(1)
	}
	for i := range wantLines {
		if strings.TrimSpace(gotLines[i]) != strings.TrimSpace(wantLines[i]) {
			fmt.Printf("Wrong answer on case %d: expected %s got %s\n", i+1, wantLines[i], gotLines[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
