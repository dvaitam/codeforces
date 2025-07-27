package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
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
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func solveCaseC(arr []int) int {
	n := len(arr)
	positions := make(map[int][]int)
	uniq := make([]int, 0)
	for i, x := range arr {
		if _, ok := positions[x]; !ok {
			uniq = append(uniq, x)
		}
		positions[x] = append(positions[x], i+1)
	}
	ans := n
	for _, x := range uniq {
		pos := positions[x]
		cnt := 0
		if pos[0] != 1 {
			cnt++
		}
		for i := 1; i < len(pos); i++ {
			if pos[i] != pos[i-1]+1 {
				cnt++
			}
		}
		if pos[len(pos)-1] != n {
			cnt++
		}
		if cnt < ans {
			ans = cnt
		}
	}
	return ans
}

func generateTests() ([][]int, string) {
	const t = 100
	r := rand.New(rand.NewSource(3))
	arrays := make([][]int, t)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := r.Intn(20) + 1
		fmt.Fprintf(&sb, "%d\n", n)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = r.Intn(n) + 1
			fmt.Fprintf(&sb, "%d ", arr[j])
		}
		fmt.Fprintln(&sb)
		arrays[i] = arr
	}
	return arrays, sb.String()
}

func verify(arrays [][]int, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	for idx, arr := range arrays {
		if !scanner.Scan() {
			return fmt.Errorf("case %d: missing output", idx+1)
		}
		var ans int
		fmt.Sscan(scanner.Text(), &ans)
		expected := solveCaseC(arr)
		if ans != expected {
			return fmt.Errorf("case %d: expected %d got %d", idx+1, expected, ans)
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output: %s", scanner.Text())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go <binary>")
		os.Exit(1)
	}
	arrays, input := generateTests()
	out, err := runBinary(os.Args[1], input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	if err := verify(arrays, out); err != nil {
		fmt.Fprintln(os.Stderr, "verification failed:", err)
		os.Exit(1)
	}
	fmt.Println("All tests passed for problem C")
}
