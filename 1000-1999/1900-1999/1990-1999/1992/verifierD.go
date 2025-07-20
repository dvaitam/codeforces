package main

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type state struct {
	pos   int
	water bool
}

func solve(n, m, k int, s string) string {
	river := make([]byte, n+2)
	river[0] = 'L'
	for i := 0; i < n; i++ {
		river[i+1] = s[i]
	}
	river[n+1] = 'L'
	const INF = int(1e9)
	surf := make([]int, n+2)
	wat := make([]int, n+2)
	for i := 0; i < n+2; i++ {
		surf[i] = INF
		wat[i] = INF
	}
	dq := list.New()
	surf[0] = 0
	dq.PushFront(state{0, false})
	for dq.Len() > 0 {
		cur := dq.Remove(dq.Front()).(state)
		var d int
		if cur.water {
			d = wat[cur.pos]
		} else {
			d = surf[cur.pos]
		}
		if cur.water {
			nxt := cur.pos + 1
			if nxt <= n+1 && river[nxt] != 'C' {
				nw := river[nxt] == 'W'
				nd := d + 1
				if nw {
					if nd < wat[nxt] {
						wat[nxt] = nd
						dq.PushBack(state{nxt, true})
					}
				} else {
					if nd < surf[nxt] {
						surf[nxt] = nd
						dq.PushBack(state{nxt, false})
					}
				}
			}
		} else {
			for j := cur.pos + 1; j <= n+1 && j <= cur.pos+m; j++ {
				if river[j] == 'C' {
					continue
				}
				nw := river[j] == 'W'
				if nw {
					if d < wat[j] {
						wat[j] = d
						dq.PushFront(state{j, true})
					}
				} else {
					if d < surf[j] {
						surf[j] = d
						dq.PushFront(state{j, false})
					}
				}
			}
		}
	}
	ans := surf[n+1]
	if wat[n+1] < ans {
		ans = wat[n+1]
	}
	if ans <= k {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Println("missing string")
			os.Exit(1)
		}
		s := scan.Text()
		expected[i] = solve(n, m, k, s)
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
