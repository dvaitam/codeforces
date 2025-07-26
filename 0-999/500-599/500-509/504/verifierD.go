package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func canRepresent(nums []int, target int) (bool, []int) {
	m := len(nums)
	bestMask := -1
	for mask := 1; mask < (1 << m); mask++ {
		xor := 0
		for i := 0; i < m; i++ {
			if mask&(1<<i) != 0 {
				xor ^= nums[i]
			}
		}
		if xor == target {
			bestMask = mask
			break
		}
	}
	if bestMask == -1 {
		return false, nil
	}
	idx := []int{}
	for i := 0; i < m; i++ {
		if bestMask&(1<<i) != 0 {
			idx = append(idx, i)
		}
	}
	return true, idx
}

func genTest(rng *rand.Rand) (string, []int) {
	m := rng.Intn(8) + 1
	nums := make([]int, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 0; i < m; i++ {
		nums[i] = rng.Intn(256)
		sb.WriteString(fmt.Sprintf("%d\n", nums[i]))
	}
	return sb.String(), nums
}

func check(out string, nums []int) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	for i, x := range nums {
		if !scanner.Scan() {
			return fmt.Errorf("not enough lines")
		}
		fields := strings.Fields(scanner.Text())
		k, err := strconv.Atoi(fields[0])
		if err != nil {
			return fmt.Errorf("invalid k")
		}
		if k == 0 {
			ok, _ := canRepresent(nums[:i], x)
			if ok {
				return fmt.Errorf("representation exists but got 0")
			}
			continue
		}
		if len(fields)-1 != k {
			return fmt.Errorf("invalid number of indices")
		}
		mask := 0
		xor := 0
		for j := 0; j < k; j++ {
			idx, err := strconv.Atoi(fields[1+j])
			if err != nil || idx < 0 || idx >= i || (mask&(1<<idx)) != 0 {
				return fmt.Errorf("bad index")
			}
			mask |= 1 << idx
			xor ^= nums[idx]
		}
		if xor != x {
			return fmt.Errorf("wrong xor")
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, nums := genTest(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(out, nums); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
