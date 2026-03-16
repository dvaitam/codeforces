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
const testcasesDRaw = `
5 4 0 4 0 3 3
2 1 0 0
1 3 1
2 2 0 0
6 5 1 5 0 3 4 0
5 4 1 4 4 0 2
3 4 0 0 3
1 3 2
6 2 1 0 2 1 0 2
4 3 2 3 3 0
3 5 2 4 3
1 5 4
5 3 0 3 3 0 3
3 4 0 0 2
6 1 1 1 0 1 1 1
2 1 1 0
5 5 5 2 1 1 2
6 1 0 0 1 1 1 1
3 4 3 4 3
2 1 0 0
4 2 1 0 2 1
6 3 1 0 3 2 1 1
1 2 2
1 5 5
5 5 3 3 2 2 2
1 4 2
3 5 4 4 2
3 2 2 1 0
1 5 1
6 5 3 2 3 0 4 3
6 5 1 0 2 4 1 5
2 3 3 0
6 2 2 0 1 0 1 0
5 4 1 3 4 0 1
2 4 0 3
6 4 2 2 3 2 4 2
1 3 0
4 1 1 0 1 1
4 4 0 4 3 2
5 2 2 2 1 1 0
4 4 4 1 0 0
5 3 2 0 0 1 0
6 5 3 0 2 0 2 3
2 3 3 1
5 2 1 1 2 0 0
2 2 1 2
6 1 0 1 0 1 0 0
5 2 1 0 1 1 2
6 3 2 3 1 3 1 0
6 5 1 5 4 0 2 0
2 4 4 1
4 5 1 5 2 5
2 4 4 0
5 1 1 0 0 0 1
5 3 3 2 3 2 0
2 3 0 1
6 2 0 1 2 0 0 1
6 1 1 1 0 0 1 0
4 1 1 1 1 1
5 4 3 1 0 3 2
2 4 0 2
1 1 0
3 1 0 1 0
3 4 4 3 3
1 1 1
3 1 0 0 1
3 1 1 0 0
5 3 0 0 3 2 1
5 2 0 0 0 2 2
2 5 2 0
4 4 0 1 1 2
3 2 0 1 0
2 4 3 2
5 2 1 2 0 0 2
1 4 1
2 1 0 1
5 1 0 0 0 0 0
6 4 0 3 3 4 2 1
2 3 1 3
1 2 1
5 3 0 3 3 1 3
4 2 1 0 2 1
4 5 0 2 2 4
3 5 0 4 4
4 2 1 0 2 1
4 1 0 0 1 0
4 4 3 0 3 4
4 2 0 1 1 1
2 5 2 4
1 3 2
2 3 3 0
2 3 1 2
3 4 0 2 4
4 1 1 1 1 0
3 2 0 0 2
3 2 0 1 2
4 4 0 0 3 4
6 5 1 1 1 1 1 1
1 5 1
4 3 2 2 0 3
`


const mod = 1000000007

func solve(n, h int, a []int) int64 {
	d := make([]int, n+2)
	for i := 1; i <= n; i++ {
		if a[i-1] > h {
			return 0
		}
		d[i] = h - a[i-1]
	}
	d[0], d[n+1] = 0, 0
	ans := int64(1)
	for i := 0; i <= n; i++ {
		cur := d[i]
		next := d[i+1]
		delta := next - cur
		switch {
		case delta == 1:
			// start
		case delta == 0:
			ans = ans * int64(cur+1) % mod
		case delta == -1:
			ans = ans * int64(cur) % mod
		default:
			return 0
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesDRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		h, _ := strconv.Atoi(parts[1])
		if len(parts)-2 < n {
			fmt.Printf("test %d invalid len\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[2+i])
			arr[i] = v
		}
		exp := solve(n, h, arr)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d\n", n, h)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&buf, "%d ", arr[i])
		}
		buf.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != fmt.Sprint(exp) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
