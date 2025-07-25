package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a / gcd(a, b) * b
}

func brute(n, m int, arr []int) (int, []int) {
	bestLen := -1
	bestL := 1
	bestIdx := []int{}
	for mask := 0; mask < (1 << uint(n)); mask++ {
		l := 1
		curIdx := []int{}
		valid := true
		for i := 0; i < n; i++ {
			if mask&(1<<uint(i)) != 0 {
				l = lcm(l, arr[i])
				if l > m {
					valid = false
					break
				}
				curIdx = append(curIdx, i)
			}
		}
		if valid {
			if len(curIdx) > bestLen {
				bestLen = len(curIdx)
				bestL = l
				bestIdx = append([]int(nil), curIdx...)
			}
		}
	}
	for i := range bestIdx {
		bestIdx[i]++
	}
	return bestL, bestIdx
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(6) + 1
		m := rand.Intn(20) + 1
		arr := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			arr[i] = rand.Intn(m) + 1
			sb.WriteString(fmt.Sprintf("%d ", arr[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expL, expIdx := brute(n, m, arr)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nInput:\n%s\nOutput:\n%s\n", t+1, err, input, out)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) < 2 {
			fmt.Printf("Test %d failed to parse output\nInput:\n%s\nOutput:\n%s\n", t+1, input, out)
			os.Exit(1)
		}
		gotL, err1 := strconv.Atoi(fields[0])
		gotK, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			fmt.Printf("Test %d failed: invalid numbers\n", t+1)
			os.Exit(1)
		}
		idxOut := []int{}
		for _, f := range fields[2:] {
			v, err := strconv.Atoi(f)
			if err != nil {
				fmt.Printf("Test %d failed: invalid index\n", t+1)
				os.Exit(1)
			}
			idxOut = append(idxOut, v)
		}
		if gotL != expL || gotK != len(expIdx) {
			fmt.Printf("Test %d failed\nInput:\n%s\nExpected: %d %d %v\nGot: %s\n", t+1, input, expL, len(expIdx), expIdx, out)
			os.Exit(1)
		}
		// order may matter but we expect ascending
		for i := 0; i < gotK && i < len(expIdx); i++ {
			if idxOut[i] != expIdx[i] {
				fmt.Printf("Test %d failed indexes\nExpected %v\nGot %v\n", t+1, expIdx, idxOut)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed!")
}
