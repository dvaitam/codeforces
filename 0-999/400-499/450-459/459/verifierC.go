package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func possible(n, k, d int) bool {
	limit := big.NewInt(int64(n))
	cur := big.NewInt(1)
	for i := 0; i < d; i++ {
		cur.Mul(cur, big.NewInt(int64(k)))
		if cur.Cmp(limit) >= 0 {
			return true
		}
	}
	return cur.Cmp(limit) >= 0
}

func validArrangement(out string, n, k, d int) bool {
	lines := strings.Fields(out)
	if len(lines) == 1 && lines[0] == "-1" {
		return false
	}
	nums := []int{}
	for _, tok := range lines {
		var v int
		if _, err := fmt.Sscanf(tok, "%d", &v); err != nil {
			return false
		}
		nums = append(nums, v)
	}
	if len(nums) != n*d {
		return false
	}
	mat := make([][]int, d)
	idx := 0
	for i := 0; i < d; i++ {
		mat[i] = make([]int, n)
		for j := 0; j < n; j++ {
			v := nums[idx]
			if v < 1 || v > k {
				return false
			}
			mat[i][j] = v
			idx++
		}
	}
	for a := 0; a < n; a++ {
		for b := a + 1; b < n; b++ {
			same := true
			for i := 0; i < d; i++ {
				if mat[i][a] != mat[i][b] {
					same = false
					break
				}
			}
			if same {
				return false
			}
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(10) + 1
		k := rand.Intn(5) + 1
		d := rand.Intn(5) + 1
		input := fmt.Sprintf("%d %d %d\n", n, k, d)
		possibleFlag := possible(n, k, d)
		outStr, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		if !possibleFlag {
			if strings.TrimSpace(outStr) != "-1" {
				fmt.Printf("wrong answer on test %d (should be -1)\ninput:%soutput:%s\n", t+1, input, outStr)
				os.Exit(1)
			}
			continue
		}
		if !validArrangement(outStr, n, k, d) {
			fmt.Printf("invalid arrangement on test %d\ninput:%soutput:%s\n", t+1, input, outStr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
