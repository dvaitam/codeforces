package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var notes = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "B", "H"}
var noteMap = map[string]int{
	"C": 0, "C#": 1, "D": 2, "D#": 3, "E": 4, "F": 5,
	"F#": 6, "G": 7, "G#": 8, "A": 9, "B": 10, "H": 11,
}
var perms = [6][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}

func classify(a, b, c string) string {
	nums := []int{noteMap[a], noteMap[b], noteMap[c]}
	for _, p := range perms {
		x := nums[p[0]]
		y := nums[p[1]]
		z := nums[p[2]]
		d1 := (y - x + 12) % 12
		d2 := (z - y + 12) % 12
		if d1 == 4 && d2 == 3 {
			return "major"
		}
		if d1 == 3 && d2 == 4 {
			return "minor"
		}
	}
	return "strange"
}

func runCase(bin, input, expected string, idx int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error on case %d: %v", idx, err)
	}
	out := strings.TrimSpace(string(outBytes))
	if out != expected {
		return fmt.Errorf("wrong answer on case %d: input=%q expected=%q got=%q", idx, strings.TrimSpace(input), expected, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	caseIdx := 0
	for _, a := range notes {
		for _, b := range notes {
			for _, c := range notes {
				input := fmt.Sprintf("%s %s %s\n", a, b, c)
				expected := classify(a, b, c)
				caseIdx++
				if err := runCase(bin, input, expected, caseIdx); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("ok %d\n", caseIdx)
}
