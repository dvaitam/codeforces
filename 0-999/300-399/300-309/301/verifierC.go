package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// The correct solution always outputs these 41 rules regardless of input.
var expectedOutput = []string{
	"0??<>1",
	"1??<>2",
	"2??<>3",
	"3??<>4",
	"4??<>5",
	"5??<>6",
	"6??<>7",
	"7??<>8",
	"8??<>9",
	"9??>>??0",
	"??<>1",
	"?0>>0?",
	"?1>>1?",
	"?2>>2?",
	"?3>>3?",
	"?4>>4?",
	"?5>>5?",
	"?6>>6?",
	"?7>>7?",
	"?8>>8?",
	"?9>>9?",
	"0?<>1",
	"1?<>2",
	"2?<>3",
	"3?<>4",
	"4?<>5",
	"5?<>6",
	"6?<>7",
	"7?<>8",
	"8?<>9",
	"9?>>??0",
	"0>>?0",
	"1>>?1",
	"2>>?2",
	"3>>?3",
	"4>>?4",
	"5>>?5",
	"6>>?6",
	"7>>?7",
	"8>>?8",
	"9>>?9",
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	expected := strings.Join(expectedOutput, "\n")

	// Provide valid input: the solution reads n and n strings
	input := "3\n1\n99\n999\n"

	out, err := runBinary(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(out) != expected {
		fmt.Fprintf(os.Stderr, "output mismatch\nexpected:\n%s\ngot:\n%s\n", expected, out)
		os.Exit(1)
	}

	fmt.Println("All tests passed")
}
