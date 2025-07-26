package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func expectedAdds(n, k int, marks []int) int {
	sum := 0
	for _, v := range marks {
		sum += v
	}
	count := n
	adds := 0
	for {
		avg := float64(sum) / float64(count)
		rounded := int(math.Floor(avg + 0.5))
		if rounded >= k {
			break
		}
		sum += k
		count++
		adds++
	}
	return adds
}

func runCase(bin string, n, k int, marks []int) error {
	input := fmt.Sprintf("%d %d\n", n, k)
	for i, v := range marks {
		if i > 0 {
			input += " "
		}
		input += strconv.Itoa(v)
	}
	input += "\n"

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	outStr := strings.TrimSpace(string(out))
	expected := strconv.Itoa(expectedAdds(n, k, marks))
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not open testcasesA.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Printf("missing n on case %d\n", i+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Printf("missing k on case %d\n", i+1)
			os.Exit(1)
		}
		k, _ := strconv.Atoi(scan.Text())
		marks := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Printf("missing mark on case %d index %d\n", i+1, j+1)
				os.Exit(1)
			}
			marks[j], _ = strconv.Atoi(scan.Text())
		}
		if err := runCase(bin, n, k, marks); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
