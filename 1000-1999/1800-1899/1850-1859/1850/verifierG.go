package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func solveCase(points [][2]int) int64 {
	mX := make(map[int]int)
	mY := make(map[int]int)
	mD1 := make(map[int]int)
	mD2 := make(map[int]int)
	for _, p := range points {
		x, y := p[0], p[1]
		mX[x]++
		mY[y]++
		mD1[x-y]++
		mD2[x+y]++
	}
	var ans int64
	for _, v := range mX {
		ans += int64(v * (v - 1))
	}
	for _, v := range mY {
		ans += int64(v * (v - 1))
	}
	for _, v := range mD1 {
		ans += int64(v * (v - 1))
	}
	for _, v := range mD2 {
		ans += int64(v * (v - 1))
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesG.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var inputs []string
	var exps []int64
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		pts := make([][2]int, n)
		idx := 1
		for i := 0; i < n; i++ {
			x, _ := strconv.Atoi(parts[idx])
			y, _ := strconv.Atoi(parts[idx+1])
			pts[i] = [2]int{x, y}
			idx += 2
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d\n", n)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%d %d\n", pts[i][0], pts[i][1])
		}
		inputs = append(inputs, sb.String())
		exps = append(exps, solveCase(pts))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}

	for idx, input := range inputs {
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		var got int64
		fmt.Fscan(strings.NewReader(out), &got)
		if got != exps[idx] {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %d got %d\n", idx+1, input, exps[idx], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
