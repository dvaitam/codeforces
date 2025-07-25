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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
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
