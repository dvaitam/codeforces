package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Test struct {
	n        int
	arr      []int
	possible bool
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func canClear(a []int) bool {
	n := len(a)
	s := make([]int, n+1)
	prefOk := make([]bool, n+1)
	prefOk[0] = true
	for i := 1; i <= n; i++ {
		if i%2 == 1 {
			s[i] = s[i-1] + a[i-1]
		} else {
			s[i] = s[i-1] - a[i-1]
		}
		if i%2 == 1 {
			prefOk[i] = prefOk[i-1] && s[i] >= 0
		} else {
			prefOk[i] = prefOk[i-1] && s[i] <= 0
		}
	}
	if prefOk[n] && s[n] == 0 {
		return true
	}
	const INF = int(1e18)
	oddMin := make([]int, n+2)
	evenMax := make([]int, n+2)
	oddMin[n+1] = INF
	evenMax[n+1] = -INF
	for i := n; i >= 0; i-- {
		if i%2 == 1 {
			if oddMin[i+1] != INF {
				oddMin[i] = min(oddMin[i+1], s[i])
			} else {
				oddMin[i] = s[i]
			}
			evenMax[i] = evenMax[i+1]
		} else {
			if evenMax[i+1] != -INF {
				evenMax[i] = max(evenMax[i+1], s[i])
			} else {
				evenMax[i] = s[i]
			}
			oddMin[i] = oddMin[i+1]
		}
	}
	for i := 1; i < n; i++ {
		if !prefOk[i-1] {
			continue
		}
		var nsI, nsIp1 int
		if i%2 == 1 {
			nsI = s[i-1] + a[i]
			nsIp1 = nsI - a[i-1]
			if nsI < 0 || nsIp1 > 0 {
				continue
			}
		} else {
			nsI = s[i-1] - a[i]
			nsIp1 = nsI + a[i-1]
			if nsI > 0 || nsIp1 < 0 {
				continue
			}
		}
		delta := nsIp1 - s[i+1]
		if delta != -s[n] {
			continue
		}
		if i+2 <= n {
			if oddMin[i+2] != INF && oddMin[i+2] < s[n] {
				continue
			}
			if evenMax[i+2] != -INF && evenMax[i+2] > s[n] {
				continue
			}
		}
		return true
	}
	return false
}

func genTests() []Test {
	r := rand.New(rand.NewSource(42))
	tests := make([]Test, 0, 100)
	for len(tests) < 100 {
		n := r.Intn(8) + 2
		arr := make([]int, n)
		for i := range arr {
			arr[i] = r.Intn(5) + 1
		}
		possible := canClear(arr)
		tests = append(tests, Test{n: n, arr: arr, possible: possible})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := genTests()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, t := range tests {
		fmt.Fprintln(&input, t.n)
		for i, v := range t.arr {
			if i > 0 {
				fmt.Fprint(&input, " ")
			}
			fmt.Fprint(&input, v)
		}
		fmt.Fprintln(&input)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("time limit exceeded")
		return
	}
	if err != nil {
		fmt.Println("execution error:", err)
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for i, t := range tests {
		if !scanner.Scan() {
			fmt.Printf("missing output for case %d\n", i+1)
			return
		}
		ans := strings.TrimSpace(scanner.Text())
		if (ans == "YES") != t.possible {
			fmt.Printf("case %d incorrect answer\n", i+1)
			return
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output detected")
		return
	}
	fmt.Println("OK")
}
