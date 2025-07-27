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

func runBinary(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genTest() (int, int, []byte) {
	n := rand.Intn(9) + 1
	s := rand.Intn(20) + n // ensure s >= n
	return n, s, []byte(fmt.Sprintf("%d %d\n", n, s))
}

func canWin(n, s int) bool { return s >= 2*n }

func hasSubarray(arr []int, target int) bool {
	for i := 0; i < len(arr); i++ {
		sum := 0
		for j := i; j < len(arr); j++ {
			sum += arr[j]
			if sum == target {
				return true
			}
		}
	}
	return false
}

func checkOutput(out string, n, s int, expectYes bool) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no output")
	}
	first := strings.ToUpper(strings.TrimSpace(lines[0]))
	if first == "NO" {
		if expectYes {
			return fmt.Errorf("expected YES got NO")
		}
		return nil
	}
	if first != "YES" {
		return fmt.Errorf("first line should be YES or NO")
	}
	if !expectYes {
		return fmt.Errorf("expected NO got YES")
	}
	if len(lines) < 3 {
		return fmt.Errorf("not enough lines in output")
	}
	nums := strings.Fields(lines[1])
	if len(nums) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(nums))
	}
	arr := make([]int, n)
	sum := 0
	for i, t := range nums {
		v, err := strconv.Atoi(t)
		if err != nil || v <= 0 {
			return fmt.Errorf("invalid array element")
		}
		arr[i] = v
		sum += v
	}
	if sum != s {
		return fmt.Errorf("array sum %d != %d", sum, s)
	}
	k, err := strconv.Atoi(strings.TrimSpace(lines[2]))
	if err != nil || k < 0 || k > s {
		return fmt.Errorf("invalid k")
	}
	if hasSubarray(arr, k) || hasSubarray(arr, s-k) {
		return fmt.Errorf("Vasya can win")
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	cand := os.Args[1]

	for i := 1; i <= 100; i++ {
		n, s, in := genTest()
		expectedYes := canWin(n, s)
		got, err := runBinary(cand, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if err := checkOutput(got, n, s, expectedYes); err != nil {
			fmt.Printf("wrong answer on test %d: %v\ninput:\n%soutput:\n%s\n", i, err, string(in), got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
