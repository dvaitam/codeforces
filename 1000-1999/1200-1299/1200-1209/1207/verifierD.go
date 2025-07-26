package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const mod int64 = 998244353

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func expected(n int, pairs [][2]int) int64 {
	fac := make([]int64, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	tmp := make([][2]int, n)
	copy(tmp, pairs)
	sort.Slice(tmp, func(i, j int) bool { return tmp[i][0] < tmp[j][0] })
	cntA := int64(1)
	for i := 0; i < n; {
		j := i
		for j < n && tmp[j][0] == tmp[i][0] {
			j++
		}
		cntA = cntA * fac[j-i] % mod
		i = j
	}
	copy(tmp, pairs)
	sort.Slice(tmp, func(i, j int) bool { return tmp[i][1] < tmp[j][1] })
	cntB := int64(1)
	for i := 0; i < n; {
		j := i
		for j < n && tmp[j][1] == tmp[i][1] {
			j++
		}
		cntB = cntB * fac[j-i] % mod
		i = j
	}
	copy(tmp, pairs)
	sort.Slice(tmp, func(i, j int) bool {
		if tmp[i][0] == tmp[j][0] {
			return tmp[i][1] < tmp[j][1]
		}
		return tmp[i][0] < tmp[j][0]
	})
	valid := true
	for i := 1; i < n; i++ {
		if tmp[i][1] < tmp[i-1][1] {
			valid = false
			break
		}
	}
	cntAB := int64(0)
	if valid {
		cntAB = 1
		for i := 0; i < n; {
			j := i
			for j < n && tmp[j] == tmp[i] {
				j++
			}
			cntAB = cntAB * fac[j-i] % mod
			i = j
		}
	}
	ans := (fac[n] - cntA - cntB + cntAB) % mod
	if ans < 0 {
		ans += mod
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	cand, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("could not open testcasesD.txt:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		f := strings.Fields(line)
		if len(f) < 1 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(f[0])
		if len(f) != 1+2*n {
			fmt.Printf("test %d bad field count\n", idx)
			os.Exit(1)
		}
		pairs := make([][2]int, n)
		pos := 1
		for i := 0; i < n; i++ {
			a, _ := strconv.Atoi(f[pos])
			b, _ := strconv.Atoi(f[pos+1])
			pairs[i] = [2]int{a, b}
			pos += 2
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			input.WriteString(fmt.Sprintf("%d %d\n", pairs[i][0], pairs[i][1]))
		}
		expect := expected(n, pairs)
		cmd := exec.Command(cand)
		cmd.Stdin = strings.NewReader(input.String())
		out, err := cmd.Output()
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		got, _ := strconv.ParseInt(gotStr, 10, 64)
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
