package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCandidate(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func maxDistance(s string) int {
	dist := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		up := int('z' - c)
		down := int(c - 'a')
		if up > down {
			dist += up
		} else {
			dist += down
		}
	}
	return dist
}

func checkAnswer(s string, k int, ans string) bool {
	if ans == "-1" {
		return maxDistance(s) < k
	}
	if len(ans) != len(s) {
		return false
	}
	dist := 0
	for i := 0; i < len(s); i++ {
		if ans[i] < 'a' || ans[i] > 'z' {
			return false
		}
		diff := int(s[i]) - int(ans[i])
		if diff < 0 {
			diff = -diff
		}
		dist += diff
	}
	if dist != k {
		return false
	}
	if maxDistance(s) < k {
		return false
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 1; i <= t; i++ {
		if !scan.Scan() {
			fmt.Printf("missing n for case %d\n", i)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		kVal, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		s := scan.Text()
		input := fmt.Sprintf("%d %d\n%s\n", n, kVal, s)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if !checkAnswer(s, kVal, out) {
			fmt.Printf("case %d failed: invalid output %q\n", i, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
