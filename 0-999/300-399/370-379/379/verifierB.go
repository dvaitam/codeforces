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

func runCase(exe string, n int, a []int) error {
	input := fmt.Sprintf("%d\n", n)
	for i, v := range a {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v)
	}
	input += "\n"
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	instr := strings.TrimSpace(out.String())
	if len(instr) > 1000000 {
		return fmt.Errorf("instruction length > 1e6")
	}
	pos := 0
	coins := make([]int, n)
	lastP := false
	for i := 0; i < len(instr); i++ {
		c := instr[i]
		switch c {
		case 'L':
			if pos == 0 {
				return fmt.Errorf("move left from position 1")
			}
			pos--
			lastP = false
		case 'R':
			if pos == n-1 {
				return fmt.Errorf("move right from position n")
			}
			pos++
			lastP = false
		case 'P':
			if lastP {
				return fmt.Errorf("two consecutive P")
			}
			coins[pos]++
			lastP = true
		default:
			return fmt.Errorf("invalid char %c", c)
		}
	}
	for i := 0; i < n; i++ {
		if coins[i] != a[i] {
			return fmt.Errorf("wallet %d expected %d got %d", i+1, a[i], coins[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := [][]int{{1, 0, 2}, {0, 3}, {5, 0, 0, 1}}
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 2
		arr := make([]int, n)
		pos := rng.Intn(n)
		arr[pos] = rng.Intn(5) + 1
		for j := 0; j < n; j++ {
			if j != pos {
				arr[j] = rng.Intn(5)
			}
		}
		cases = append(cases, arr)
	}
	for idx, arr := range cases {
		if err := runCase(exe, len(arr), arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
