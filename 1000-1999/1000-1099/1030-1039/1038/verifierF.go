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

func expected(n int, s string) string {
	m := len(s)
	a := make([]int, m+2)
	for i := 1; i <= m; i++ {
		a[i] = int(s[i-1] - '0')
	}
	Next := make([]int, m+2)
	j := 0
	for i := 2; i <= m; i++ {
		for j > 0 && a[j+1] != a[i] {
			j = Next[j]
		}
		if a[j+1] == a[i] {
			j++
		}
		Next[i] = j
	}
	used := make([]int, n+2)
	F := make([][][]int64, n+2)
	for i := 0; i <= n; i++ {
		F[i] = make([][]int64, m)
		for k := 0; k < m; k++ {
			F[i][k] = make([]int64, 2)
		}
	}
	tmp := make([][]int64, m)
	for k := 0; k < m; k++ {
		tmp[k] = make([]int64, 2)
	}
	var ans int64
	clear := func() {
		for i := 0; i <= n; i++ {
			for k := 0; k < m; k++ {
				F[i][k][0] = 0
				F[i][k][1] = 0
			}
		}
	}
	var dp func(st, ed int)
	dp = func(st, ed int) {
		for i := st; i <= ed; i++ {
			for k := 0; k < m; k++ {
				for ok := 0; ok < 2; ok++ {
					cnt := F[i-1][k][ok]
					if cnt == 0 {
						continue
					}
					for cur := 0; cur < 2; cur++ {
						if used[i] != -1 && used[i] != cur {
							continue
						}
						j := k
						for j > 0 && a[j+1] != cur {
							j = Next[j]
						}
						if a[j+1] == cur {
							j++
						}
						if j == m {
							F[i][0][1] += cnt
						} else {
							F[i][j][ok] += cnt
						}
					}
				}
			}
		}
	}
	for i := 1; i <= n; i++ {
		used[i] = -1
	}
	clear()
	F[0][0][0] = 1
	dp(1, n)
	for k := 0; k < m; k++ {
		ans += F[n][k][1]
	}
	for length := 1; length < m; length++ {
		for i := 1; i <= n; i++ {
			used[i] = -1
		}
		for i := 1; i <= length; i++ {
			used[n-length+i] = a[i]
		}
		for i := length + 1; i <= m; i++ {
			used[i-length] = a[i]
		}
		clear()
		start := n - length + 1
		F[start][0][0] = 1
		dp(start+1, n)
		for k := 0; k < m; k++ {
			for ok := 0; ok < 2; ok++ {
				tmp[k][ok] = F[n][k][ok]
			}
		}
		clear()
		for k := 0; k < m; k++ {
			for ok := 0; ok < 2; ok++ {
				F[0][k][ok] = tmp[k][ok]
			}
		}
		dp(1, n)
		for k := 0; k < m; k++ {
			ans += F[n][k][0]
		}
	}
	return fmt.Sprintf("%d", ans)
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Println("could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		s := scan.Text()
		input := fmt.Sprintf("%d\n%s\n", n, s)
		exp := expected(n, s)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
