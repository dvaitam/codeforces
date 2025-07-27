package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	return out.String(), err
}

func expectedE(a []int) int {
	n := len(a)
	prefix := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + a[i]
	}
	seen := make([]bool, n+1)
	for l := 0; l < n; l++ {
		for r := l + 2; r <= n; r++ {
			sum := prefix[r] - prefix[l]
			if sum > n {
				break
			}
			seen[sum] = true
		}
	}
	ans := 0
	for i := 0; i < n; i++ {
		if seen[a[i]] {
			ans++
		}
	}
	return ans
}

func runCase(bin string, arr []int) error {
	n := len(arr)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	out, err := runBinary(bin, sb.String())
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out)
	}
	gotStr := strings.TrimSpace(out)
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("invalid output %q", gotStr)
	}
	exp := expectedE(arr)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(sc.Text())
	for i := 0; i < t; i++ {
		sc.Scan()
		n, _ := strconv.Atoi(sc.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			sc.Scan()
			v, _ := strconv.Atoi(sc.Text())
			arr[j] = v
		}
		if err := runCase(bin, arr); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
