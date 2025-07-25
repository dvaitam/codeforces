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

type op struct {
	pos int
	typ byte
}

func expectedCase(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	a := make([]int64, n)
	var suma int64
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		suma += a[i]
	}
	var k int
	fmt.Fscan(reader, &k)
	b := make([]int64, k)
	var sumb int64
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &b[i])
		sumb += b[i]
	}
	if suma != sumb {
		return "NO\n"
	}
	flag := make([]int, n)
	nxt := make([]int, n)
	var cur, cnt int
	var temp int64
	for i := 0; i < n; i++ {
		temp += a[i]
		if temp == b[cur] {
			flag[cnt] = i
			cnt++
			cur++
			temp = 0
		} else if temp > b[cur] {
			return "NO\n"
		}
	}
	if cur != k {
		return "NO\n"
	}
	var ans []op
	initNxt := func(left, right int) {
		for i := left; i < right; i++ {
			nxt[i] = i + 1
		}
		if right < len(nxt) {
			nxt[right] = -1
		}
	}
	check := func(left, right int) bool {
		return nxt[left] == -1
	}
	var findIndex func(left, index int) int
	findIndex = func(left, index int) int {
		cnt := 0
		for i := left; i != -1; i = nxt[i] {
			cnt++
			if index == i {
				return cnt
			}
		}
		return -1
	}
	solve := func(left, right, pre int) int64 {
		var mx int64 = -1
		var index int
		first := true
		for i := left; i != -1; i = nxt[i] {
			j := nxt[i]
			if j != -1 && a[i] != a[j] {
				sum := a[i] + a[j]
				if first {
					first = false
					mx = sum
					index = i
				} else if sum > mx {
					mx = sum
					index = i
				}
			}
		}
		if mx != -1 {
			pos1 := findIndex(left, index)
			pos2 := findIndex(left, nxt[index])
			if a[index] > a[nxt[index]] {
				ans = append(ans, op{pos: pre + pos1, typ: 'R'})
			} else {
				ans = append(ans, op{pos: pre + pos2, typ: 'L'})
			}
			a[index] += a[nxt[index]]
			nxt[index] = nxt[nxt[index]]
		}
		return mx
	}
	noFail := false
	for i := 0; i < cnt; i++ {
		var left int
		right := flag[i]
		if i == 0 {
			left = 0
		} else {
			left = flag[i-1] + 1
		}
		initNxt(left, right)
		for !check(left, right) {
			p := solve(left, right, i)
			if p == -1 {
				noFail = true
				break
			}
		}
		if noFail {
			break
		}
	}
	if noFail {
		return "NO\n"
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for _, t := range ans {
		sb.WriteString(fmt.Sprintf("%d %c\n", t.pos, t.typ))
	}
	return sb.String()
}

func runCase(exe string, input string, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		a := make([]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			a[j], _ = strconv.Atoi(scan.Text())
		}
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		b := make([]int, k)
		for j := 0; j < k; j++ {
			scan.Scan()
			b[j], _ = strconv.Atoi(scan.Text())
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(strconv.Itoa(a[j]))
		}
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("%d\n", k))
		for j := 0; j < k; j++ {
			if j > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(strconv.Itoa(b[j]))
		}
		sb.WriteString("\n")
		exp := expectedCase(sb.String())
		if err := runCase(exe, sb.String(), exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
