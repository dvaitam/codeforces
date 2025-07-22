package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func expected(nums []int) string {
	n := len(nums)
	switch n {
	case 1:
		if nums[0] > 0 {
			return "BitLGM"
		}
		return "BitAryo"
	case 2:
		sort.Ints(nums)
		phi := (1 + math.Sqrt(5)) / 2
		k := nums[1] - nums[0]
		t := int(math.Floor(float64(k) * phi))
		if t == nums[0] {
			return "BitAryo"
		}
		return "BitLGM"
	case 3:
		x := nums[0] ^ nums[1] ^ nums[2]
		if x != 0 {
			return "BitLGM"
		}
		return "BitAryo"
	default:
		return "BitLGM"
	}
}

func runCase(bin string, line string, idx int) error {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil || n < 1 || n > 3 {
		return fmt.Errorf("test %d: invalid n", idx)
	}
	if len(fields) != 1+n {
		return fmt.Errorf("test %d: expected %d numbers got %d", idx, 1+n, len(fields))
	}
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		v, _ := strconv.Atoi(fields[1+i])
		nums[i] = v
	}
	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprintf(&input, "%d", nums[i])
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	ans := strings.TrimSpace(out.String())
	exp := expected(nums)
	if ans != exp {
		return fmt.Errorf("expected %s got %s", exp, ans)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open testcasesD.txt: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		if err := runCase(bin, line, idx); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
