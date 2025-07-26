package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func solveB(n, B int, arr []int) string {
	var costs []int
	odd, even := 0, 0
	for i := 0; i < n-1; i++ {
		if arr[i]%2 == 0 {
			even++
		} else {
			odd++
		}
		if odd == even {
			diff := arr[i] - arr[i+1]
			if diff < 0 {
				diff = -diff
			}
			costs = append(costs, diff)
		}
	}
	sort.Ints(costs)
	total := 0
	cnt := 0
	for _, c := range costs {
		if total+c <= B {
			total += c
			cnt++
		} else {
			break
		}
	}
	return fmt.Sprintf("%d", cnt)
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		budget, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			arr[i], _ = strconv.Atoi(scan.Text())
		}
		var inSB strings.Builder
		inSB.WriteString(fmt.Sprintf("%d %d\n", n, budget))
		for i, v := range arr {
			if i > 0 {
				inSB.WriteByte(' ')
			}
			inSB.WriteString(strconv.Itoa(v))
		}
		inSB.WriteByte('\n')
		input := inSB.String()
		expected := solveB(n, budget, arr)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("case %d failed: expected %q got %q\ninput:\n%s", caseIdx+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
