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

func fact(n int) int {
	f := 1
	for i := 2; i <= n; i++ {
		f *= i
	}
	return f
}

func ord(p []int) int {
	n := len(p)
	used := make([]bool, n)
	res := 0
	f := fact(n - 1)
	for i := 0; i < n; i++ {
		cnt := 0
		for j := 0; j < p[i]; j++ {
			if !used[j] {
				cnt++
			}
		}
		res += cnt * f
		used[p[i]] = true
		if i < n-1 {
			f /= n - 1 - i
		}
	}
	return res
}

func permFromIndex(n, idx int) []int {
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = i
	}
	res := make([]int, n)
	f := fact(n - 1)
	for i := 0; i < n; i++ {
		pos := 0
		if f > 0 {
			pos = idx / f
			idx %= f
		}
		res[i] = nums[pos]
		nums = append(nums[:pos], nums[pos+1:]...)
		if i < n-1 {
			f /= n - 1 - i
		}
	}
	return res
}

func genTest(rng *rand.Rand) (string, []int) {
	n := rng.Intn(6) + 1
	p := rng.Perm(n)
	q := rng.Perm(n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range q {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	idx := (ord(p) + ord(q)) % fact(n)
	expect := permFromIndex(n, idx)
	return sb.String(), expect
}

func check(out string, expect []int) error {
	fields := strings.Fields(out)
	if len(fields) != len(expect) {
		return fmt.Errorf("expected %d numbers", len(expect))
	}
	seen := make(map[int]bool)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid number")
		}
		if v < 0 || v >= len(expect) || seen[v] {
			return fmt.Errorf("invalid permutation")
		}
		if v != expect[i] {
			return fmt.Errorf("wrong result")
		}
		seen[v] = true
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genTest(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(out, exp); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
