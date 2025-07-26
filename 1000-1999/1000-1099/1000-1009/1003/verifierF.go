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

func expected(n int, words []string) int {
	a := make([]int, n+1)
	lens := make([]int, n+1)
	m := make(map[string]int)
	for i := 1; i <= n; i++ {
		s := words[i-1]
		if m[s] == 0 {
			m[s] = i
			a[i] = i
			lens[i] = len(s)
		} else {
			a[i] = m[s]
		}
	}
	sum := n - 1
	for i := 1; i <= n; i++ {
		sum += lens[a[i]]
	}
	ans := 0
	for i := 1; i <= n; i++ {
		for j := i; j <= n; j++ {
			cnt := 1
			x := -1
			for k := i; k <= j; k++ {
				x += lens[a[k]]
			}
			poi := j + 1
			length := j - i + 1
			for poi+length-1 <= n {
				ok := true
				for k := 0; k < length; k++ {
					if a[i+k] != a[poi+k] {
						ok = false
						break
					}
				}
				if ok {
					cnt++
					poi += length
				} else {
					poi++
				}
			}
			if cnt > 1 {
				if v := x * cnt; v > ans {
					ans = v
				}
			}
		}
	}
	return sum - ans
}

func runCase(bin string, words []string) error {
	n := len(words)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, w := range words {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(w)
	}
	sb.WriteByte('\n')
	input := sb.String()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("invalid output %q", gotStr)
	}
	want := expected(n, words)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Println("could not open testcasesF.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		words := make([]string, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			words[j] = scan.Text()
		}
		if err := runCase(bin, words); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
