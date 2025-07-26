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

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

var word = "kotlin"

func generatePieces() ([]string, []int) {
	times := rand.Intn(4) + 1
	full := strings.Repeat(word, times)
	n := rand.Intn(len(full)) + 1
	cuts := make([]int, 0, n-1)
	for len(cuts) < n-1 {
		pos := rand.Intn(len(full)-1) + 1
		found := false
		for _, c := range cuts {
			if c == pos {
				found = true
				break
			}
		}
		if !found {
			cuts = append(cuts, pos)
		}
	}
	sortInts(cuts)
	pieces := make([]string, 0, n)
	last := 0
	for _, c := range cuts {
		pieces = append(pieces, full[last:c])
		last = c
	}
	pieces = append(pieces, full[last:])
	idx := rand.Perm(n)
	shuffled := make([]string, n)
	for i, p := range pieces {
		shuffled[idx[i]] = p
	}
	return shuffled, idx
}

func sortInts(a []int) { // simple insertion sort for small slices
	for i := 1; i < len(a); i++ {
		v := a[i]
		j := i - 1
		for j >= 0 && a[j] > v {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = v
	}
}

func checkOutput(out string, pieces []string) bool {
	fields := strings.Fields(out)
	if len(fields) != len(pieces) {
		return false
	}
	seen := make([]bool, len(pieces))
	var sb strings.Builder
	for _, f := range fields {
		var idx int
		if _, err := fmt.Sscan(f, &idx); err != nil {
			return false
		}
		if idx < 1 || idx > len(pieces) || seen[idx-1] {
			return false
		}
		seen[idx-1] = true
		sb.WriteString(pieces[idx-1])
	}
	s := sb.String()
	if len(s)%len(word) != 0 || len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] != word[i%6] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]

	for t := 0; t < 100; t++ {
		pieces, idx := generatePieces()
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", len(pieces)))
		for i := 0; i < len(pieces); i++ {
			if i > 0 {
				input.WriteByte('\n')
			}
			input.WriteString(pieces[i])
		}
		input.WriteByte('\n')
		out, err := runProgram(bin, []byte(input.String()))
		if err != nil || !checkOutput(out, pieces) {
			fmt.Println("Test", t+1, "failed")
			fmt.Println("Input:\n", input.String())
			fmt.Println("Output:", out)
			if err != nil {
				fmt.Println("Error:", err)
			}
			fmt.Println("Expected order (one possible):", idx)
			return
		}
	}
	fmt.Println("All tests passed")
}
