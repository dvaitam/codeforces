package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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
	"?>>??",
	">>?",
}

func runBinary(path string) (string, error) {
	cmd := exec.Command(path)
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
	passed := 0
	for i := 0; i < 100; i++ { // run 100 times
		out, err := runBinary(bin)
		if err != nil {
			fmt.Printf("run %d: runtime error: %v\n", i+1, err)
			continue
		}
		if strings.TrimSpace(out) != expected {
			fmt.Printf("run %d failed: output mismatch\n", i+1)
		} else {
			passed++
		}
	}
	fmt.Printf("passed %d/100 runs\n", passed)
	if passed != 100 {
		os.Exit(1)
	}
}
