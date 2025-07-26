package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type tile struct {
	num  int
	suit byte
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseTile(s string) tile {
	return tile{num: int(s[0] - '0'), suit: s[1]}
}

func solveCase(a, b, c string) int {
	t := [3]tile{parseTile(a), parseTile(b), parseTile(c)}
	if t[0] == t[1] && t[1] == t[2] {
		return 0
	}
	if t[0].suit == t[1].suit && t[1].suit == t[2].suit {
		nums := []int{t[0].num, t[1].num, t[2].num}
		sort.Ints(nums)
		if nums[0]+1 == nums[1] && nums[1]+1 == nums[2] {
			return 0
		}
	}
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 3; j++ {
			ti, tj := t[i], t[j]
			if ti == tj {
				return 1
			}
			if ti.suit == tj.suit {
				d := ti.num - tj.num
				if d < 0 {
					d = -d
				}
				if d <= 2 {
					return 1
				}
			}
		}
	}
	return 2
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var a, b, c string
		fmt.Sscan(line, &a, &b, &c)
		expect := solveCase(a, b, c)
		input := fmt.Sprintf("%s %s %s\n", a, b, c)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx, err, input)
			os.Exit(1)
		}
		if got != fmt.Sprintf("%d", expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", idx, expect, got, input)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
