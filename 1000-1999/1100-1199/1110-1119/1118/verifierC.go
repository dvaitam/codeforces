package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func checkMatrix(n int, nums []int, mat [][]int) bool {
	cnt := make(map[int]int)
	for _, v := range nums {
		cnt[v]++
	}
	for i := 0; i < n; i++ {
		if len(mat[i]) != n {
			return false
		}
		for j := 0; j < n; j++ {
			cnt[mat[i][j]]--
		}
	}
	for _, v := range cnt {
		if v != 0 {
			return false
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if mat[i][j] != mat[n-1-i][j] || mat[i][j] != mat[i][n-1-j] || mat[i][j] != mat[n-1-i][n-1-j] {
				return false
			}
		}
	}
	return true
}

func possible(n int, nums []int) bool {
	cnt := make(map[int]int)
	for _, v := range nums {
		cnt[v]++
	}
	k := n / 2
	quads := 0
	pairs := 0
	singles := 0
	for _, c := range cnt {
		quads += c / 4
		c %= 4
		pairs += c / 2
		c %= 2
		singles += c
	}
	if n%2 == 0 {
		return singles == 0 && pairs == 0 && quads >= k*k
	}
	return singles == 1 && pairs >= 2*k && quads >= k*k
}

func parseMatrix(lines []string, n int) ([][]int, bool) {
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		fields := strings.Fields(lines[i])
		if len(fields) != n {
			return nil, false
		}
		row := make([]int, n)
		for j := 0; j < n; j++ {
			var val int
			if _, err := fmt.Sscan(fields[j], &val); err != nil {
				return nil, false
			}
			row[j] = val
		}
		mat[i] = row
	}
	return mat, true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(3)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(5) + 1
		total := n * n
		nums := make([]int, total)
		for i := 0; i < total; i++ {
			nums[i] = rand.Intn(3)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range nums {
			if i%20 == 0 && i > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(fmt.Sprintf("%d ", v))
		}
		sb.WriteString("\n")
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t, err)
			fmt.Println(out)
			return
		}
		out = strings.TrimSpace(out)
		lines := strings.Split(out, "\n")
		if strings.EqualFold(lines[0], "NO") {
			if possible(n, nums) {
				fmt.Printf("test %d failed: expected YES but got NO\ninput:\n%s\n", t, sb.String())
				return
			}
			continue
		}
		if !strings.EqualFold(lines[0], "YES") {
			fmt.Printf("test %d invalid output\n", t)
			return
		}
		if !possible(n, nums) {
			fmt.Printf("test %d failed: expected NO but got YES\ninput:\n%s\n", t, sb.String())
			return
		}
		if len(lines)-1 != n {
			fmt.Printf("test %d failed: expected %d rows got %d\n", t, n, len(lines)-1)
			return
		}
		mat, ok := parseMatrix(lines[1:], n)
		if !ok {
			fmt.Printf("test %d failed: cannot parse matrix\n", t)
			return
		}
		if !checkMatrix(n, nums, mat) {
			fmt.Printf("test %d failed: matrix invalid\n", t)
			return
		}
	}
	fmt.Println("all tests passed")
}
