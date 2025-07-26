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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	t := 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		k := rand.Intn(3) + 1
		n := 2*k + 1
		nums := rand.Perm(4*k + 1)
		board1 := make([]string, n)
		board2 := make([]string, n)
		pos := 0
		emptyPlaced := false
		for j := 0; j < n; j++ {
			if !emptyPlaced && rand.Intn(4*k+1-pos) == 0 {
				board1[j] = "E"
				emptyPlaced = true
			} else {
				board1[j] = fmt.Sprintf("%d", nums[pos]+1)
				pos++
			}
		}
		for j := 0; j < n; j++ {
			if !emptyPlaced && rand.Intn(4*k+1-pos) == 0 {
				board2[j] = "E"
				emptyPlaced = true
			} else {
				board2[j] = fmt.Sprintf("%d", nums[pos]+1)
				pos++
			}
		}
		if !emptyPlaced {
			board2[n-1] = "E"
		}
		fmt.Fprintln(&input, k)
		fmt.Fprintln(&input, strings.Join(board1, " "))
		fmt.Fprintln(&input, strings.Join(board2, " "))
		expected[i] = "SURGERY FAILED"
	}
	cmd := exec.Command(bin)
	cmd.Stdin = &input
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("binary error:", err)
		fmt.Print(string(out))
		return
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) != t {
		fmt.Printf("expected %d lines, got %d\n", t, len(lines))
		fmt.Print(string(out))
		return
	}
	for i := 0; i < t; i++ {
		if strings.TrimSpace(lines[i]) != expected[i] {
			fmt.Printf("Test %d failed: expected %s got %s\n", i+1, expected[i], strings.TrimSpace(lines[i]))
			return
		}
	}
	fmt.Println("All tests passed!")
}
