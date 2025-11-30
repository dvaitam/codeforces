package main

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const embeddedTestcasesD = `100
9 9 2
WCWCCLCLW
10 9 3
LCWCCWWCLL
6 5 3
CLCLLC
3 2 0
WWC
14 7 6
CCWLWLLLWLWCWC
11 7 8
WCWCCWCLWCL
10 10 10
CLCWCCCLCC
8 5 4
LLWCWLWL
15 3 0
WWWLLCCLWCCWCWC
9 1 4
LLLCCLLWW
10 3 0
WWWLWWWCWC
19 2 19
CWWCCCLWWWCWCWLWCWL
14 10 9
CLLCCWWWCWCCWC
17 1 1
CLWWCWWCCWLWLWWCW
11 7 1
LCCCLWCLCWL
12 3 10
WCCLLCWWCLWL
4 3 1
CWWL
5 1 4
LWCLW
12 2 12
CWCLWWCWWWCW
11 7 9
WLWLLLWCCWC
9 1 7
CCCWCWLLC
11 2 3
LLCCLWCLLWC
5 2 4
WLCLC
19 7 1
CLWLWCWLWLLLCLLLWLL
17 10 12
LWLWCCCWLWWLLLLLL
4 4 0
CLCC
17 6 5
WLWWCWCWWWLWWLLCL
14 2 9
LLWWCCCWCLCWLW
17 6 13
CWWLLLCWCCLWLWLLW
13 9 12
WLWCLCCCWCLCW
20 9 3
CCLWLLWLCLCLCWLLWLLL
17 5 9
LLCCCCLLCCLCLCWCL
4 2 1
CLWC
14 5 5
CWWCWLWCLWCLWC
20 10 16
CWLCWLLLWCWCCWCCLLWL
6 5 4
WLCCWC
15 10 8
LWLWWWLLCWLWWLW
17 10 6
WCCCLWCLWCLWLLCCC
4 2 2
LLWL
7 5 0
WLLWLWL
4 1 0
WWLW
16 10 10
LLWWWWWLLCCCWWLC
12 2 4
LCWLWCWLCCWC
13 8 4
CCCCWLWCCCCLC
17 9 11
LCWCCCCLCCCCLLWCW
5 1 4
LWCCC
15 9 9
CLWWWLCLLCLCCCC
14 6 1
WWWLLLWCWLCCWW
14 8 8
LWWLLWCLCLCCLL
8 7 6
CLCCLWLC
18 3 18
LLWWWLWWWWCCWCWCCC
11 8 0
CLCCLLCLWLC
16 4 6
CLLCCWLCWCLLWLCL
2 2 2
LC
4 4 4
LWWW
13 3 1
CLCLCWLLLWCLW
12 2 12
LWLLCLWWWCLC
15 4 15
WLWWCWWLWLCWWWW
18 6 12
LLLWWLWCLLLLLCWCLC
2 1 0
CL
16 8 5
LLWWWWLCCLLWLWWL
9 4 1
WWLLWCLWL
3 3 2
WWW
17 4 17
WLWCWWLLCWLCCLLLC
4 1 1
CLWL
8 6 3
WWWWCLWW
14 2 4
LWLCWCWCCWLCCW
19 6 5
WWWLWCWWWLCCCCCWWLL
7 6 0
WLCCCLW
9 1 9
LWWCCWWLL
8 7 1
WCWLWLCL
9 1 2
CLCLLLWLL
19 9 13
CCWCLWLCLLLCLLLWLLW
14 4 13
CCLCWCCWWCCWLW
16 7 6
LCLWWCCLWWLWLCCC
13 7 12
CLCCLLWLWLWWL
10 4 5
LCLCLCCWCC
15 6 6
CWWWWCLWLLLLCLW
18 9 1
CWLLCLWWWLWWCLCWLL
16 5 13
WLLLWCWWCWLWWWWC
19 4 11
WWLWWWWCLCLLLCCLCLL
3 1 1
WWW
19 1 15
LLLWWWLCCLWWWLCWWCW
12 1 0
WLCCLCLLCWCW
19 3 8
LLLLWCLWWCCCCWWCCLW
7 1 2
CLWWWWW
10 8 4
CLCWLCLWWL
16 10 8
LWCCCLWLWCWLWCLW
3 3 3
CWL
11 6 2
CWWWLWLLWWC
14 9 4
LLWWWLWLLCLCCL
5 1 2
CLWWW
16 5 9
CCLLLWLLLCLCWCWC
2 1 1
LL
7 3 3
CLLLCLC
8 7 1
WWLLWLLL
6 3 1
WWCLLL
13 3 7
LLWCWLCWWLCWC
16 10 13
WCCWWLCWWLCWCCLW
7 3 4
LLLCLWC`

type state struct {
	pos   int
	water bool
}

// solveCase mirrors the reference solution from 1992D.go.
func solveCase(n, m, k int, s string) string {
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

func expectedOutputs() ([]string, error) {
	scan := bufio.NewScanner(strings.NewReader(embeddedTestcasesD))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("missing test count")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %w", err)
	}
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("missing n for test %d", i+1)
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid n for test %d: %w", i+1, err)
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("missing m for test %d", i+1)
		}
		m, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid m for test %d: %w", i+1, err)
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("missing k for test %d", i+1)
		}
		k, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid k for test %d: %w", i+1, err)
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("missing river for test %d", i+1)
		}
		s := scan.Text()
		expected[i] = solveCase(n, m, k, s)
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return expected, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	expected, err := expectedOutputs()
	if err != nil {
		fmt.Println("could not parse embedded tests:", err)
		os.Exit(1)
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewBufferString(embeddedTestcasesD)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < len(expected); i++ {
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
