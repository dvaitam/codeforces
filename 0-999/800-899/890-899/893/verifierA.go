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

func runBinary(bin, input string) (string, error) {
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

func solveCase(wins []int) string {
	p1, p2, s := 1, 2, 3
	for _, w := range wins {
		if w != p1 && w != p2 {
			return "NO"
		}
		var loser int
		if w == p1 {
			loser = p2
		} else {
			loser = p1
		}
		p1, p2, s = w, s, loser
	}
	return "YES"
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(100) + 1
	wins := make([]int, n)
	spectators := make([]int, n)
	p1, p2, s := 1, 2, 3
	for i := 0; i < n; i++ {
		if r.Intn(2) == 0 {
			wins[i] = p1
			spectators[i] = s
			p1, p2, s = p1, s, p2
		} else {
			wins[i] = p2
			spectators[i] = s
			p1, p2, s = p2, s, p1
		}
	}
	if r.Intn(2) == 0 {
		idx := r.Intn(n)
		wins[idx] = spectators[idx]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprint(wins[i]))
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
	expect := solveCase(wins)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(r)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
