package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type SegTree struct {
	size int
	data []int
}

func NewSegTree(a []int) *SegTree {
	n := len(a) - 1
	size := 1
	for size < n {
		size <<= 1
	}
	data := make([]int, 2*size)
	for i := 1; i <= n; i++ {
		data[size+i-1] = a[i]
	}
	for i := size - 1; i >= 1; i-- {
		left, right := data[2*i], data[2*i+1]
		if left > right {
			data[i] = left
		} else {
			data[i] = right
		}
	}
	return &SegTree{size, data}
}

func (st *SegTree) Query(l, r int) int {
	if l > r {
		return 0
	}
	l += st.size - 1
	r += st.size - 1
	res := 0
	for l <= r {
		if l%2 == 1 {
			if st.data[l] > res {
				res = st.data[l]
			}
			l++
		}
		if r%2 == 0 {
			if st.data[r] > res {
				res = st.data[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func expected(a []int) int64 {
	n := len(a) - 1
	st := NewSegTree(a)
	var ans int64
	for i := 1; i < n; i++ {
		cur := i
		reach := a[i]
		dist := 1
		for {
			if reach >= n {
				ans += int64(dist) * int64(n-cur)
				break
			}
			ans += int64(dist) * int64(reach-cur)
			newReach := st.Query(cur+1, reach)
			cur = reach
			reach = newReach
			dist++
		}
	}
	return ans
}

const testcasesERaw = `100
11 6 8 4 11 9 8 8 9 10 11
9 5 6 8 5 7 7 8 9
6 3 6 4 5 6
11 11 10 6 6 6 7 9 9 10 11
6 4 4 6 5 6
13 5 9 8 5 11 10 9 10 12 11 13 13
15 11 12 4 14 11 8 12 11 12 14 14 13 15 15
13 4 3 8 5 11 13 11 9 13 12 13 13
9 2 8 5 9 7 7 8 9
7 6 5 6 6 6 7
15 6 3 10 6 9 12 13 10 12 13 12 14 15 15
4 2 3 4
9 4 8 4 5 9 7 8 9
15 6 10 14 11 8 7 8 15 13 13 13 13 14 15
3 2 3
7 3 3 7 6 6 7
11 4 11 11 8 11 9 11 10 11 11
8 3 3 7 6 8 8 8
4 2 4 4
10 10 8 4 7 10 7 9 10 10
15 9 7 15 9 11 9 8 12 14 13 14 15 15 15
6 6 5 5 6 6
13 8 8 6 12 11 9 12 10 11 11 13 13
6 2 6 4 6 6
11 10 7 4 6 7 11 11 11 10 11
14 12 5 14 5 13 10 9 9 11 11 13 13 14
5 2 4 5 5
8 7 7 4 6 6 8 8
7 5 5 7 5 7 7
10 6 4 6 9 10 9 10 10 10
15 4 9 10 14 14 12 15 10 11 15 15 15 15 15
4 4 3 4
15 12 3 10 6 11 15 10 11 15 15 15 14 15 15
9 6 6 9 8 9 7 9 9
9 6 6 4 7 8 8 9 9
9 2 3 9 9 9 7 9 9
14 4 9 4 12 14 11 9 12 10 12 13 14 14
13 6 7 11 12 7 7 9 11 12 13 13 13
13 10 3 11 10 11 12 13 13 11 11 12 13
10 9 4 9 9 7 7 9 10 10
15 15 9 11 11 9 11 15 15 10 13 12 15 15 15
12 9 7 6 7 9 12 12 12 11 11 12
11 10 3 5 10 7 9 8 9 10 11
13 8 4 10 12 6 7 8 10 10 13 12 13
13 9 5 6 7 12 13 9 9 11 11 12 13
13 3 7 4 6 12 11 9 9 13 11 13 13
9 7 8 6 5 7 9 9 9
11 10 7 7 7 6 11 9 10 11 11
14 5 8 11 11 8 9 13 9 12 11 14 13 14
15 3 11 13 12 9 13 15 12 12 11 12 15 15 15
15 4 9 14 8 7 8 15 15 11 12 15 14 14 15
7 6 5 7 6 6 7
12 6 9 9 5 10 11 9 10 11 11 12
2 2
4 2 3 4
9 9 8 6 5 9 8 9 9
7 5 6 4 5 6 7
3 3 3
15 9 5 12 9 7 11 10 15 12 12 12 14 14 15
14 7 8 4 5 13 14 14 10 10 13 13 14 14
13 4 3 7 5 9 12 8 13 13 13 13 13
4 3 4 4
3 3 3
14 5 5 14 8 6 13 13 11 10 14 12 14 14
4 2 3 4
3 2 3
14 9 4 11 5 9 7 11 12 10 11 13 14 14
6 2 3 4 5 6
5 3 3 5 5
3 2 3
15 13 6 13 5 13 15 13 12 13 15 15 15 14 15
14 6 14 11 8 14 8 9 11 13 11 13 14 14
4 2 3 4
5 2 4 4 5
2 2
13 5 3 7 9 11 7 8 12 12 12 12 13
12 4 12 4 12 7 12 8 9 12 12 12
3 3 3
7 2 5 5 5 7 7
12 7 7 4 6 10 12 12 9 10 11 12
4 3 4 4
11 2 6 10 6 8 8 10 10 10 11
11 11 5 5 7 10 11 9 11 10 11
14 8 9 5 10 11 14 9 12 14 14 12 14 14
14 12 5 12 12 13 7 8 10 13 11 12 13 14
11 5 10 9 8 9 8 9 10 10 11
10 10 8 10 6 8 9 8 9 10
8 4 4 5 6 7 8 8
3 2 3
3 3 3
14 14 8 11 12 14 8 14 10 13 11 12 13 14
2 2
12 2 12 8 11 10 10 12 12 11 11 12
14 8 8 13 5 14 10 10 11 10 11 12 14 14
6 2 5 6 5 6
2 2
12 6 8 12 8 9 11 8 9 10 12 12
8 2 8 7 8 6 7 8
8 4 3 8 7 6 7 8
6 4 5 6 5 6`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	data := []byte(testcasesERaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expect := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n+1)
		for j := 1; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			val, _ := strconv.Atoi(scan.Text())
			arr[j] = val
		}
		arr[n] = n
		ans := expected(arr)
		expect[i] = fmt.Sprintf("%d", ans)
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("execution failed: %v\n%s", err, out)
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
		if got != expect[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expect[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
