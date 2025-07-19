package main

import (
	"bufio"
	"fmt"
	"os"
)

type op struct {
	pos int
	typ byte
}

var (
	n, k int
	a, b []int64
	flag []int
	nxt  []int
	ans  []op
)

func initNxt(left, right int) {
	for i := left; i < right; i++ {
		nxt[i] = i + 1
	}
	if right < len(nxt) {
		nxt[right] = -1
	}
}

func check(left, right int) bool {
	return nxt[left] == -1
}

// findIndex returns 1-based position of index in current linked list starting at left
func findIndex(left, index int) int {
	cnt := 0
	for i := left; i != -1; i = nxt[i] {
		cnt++
		if index == i {
			return cnt
		}
	}
	return -1
}

// solve merges one optimal adjacent pair in [left..right], pre as offset
// returns sum of merged pair or -1 if no merge possible
func solve(left, right, pre int) int64 {
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
		// merge
		a[index] += a[nxt[index]]
		nxt[index] = nxt[nxt[index]]
	}
	return mx
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a = make([]int64, n)
	var suma int64
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		suma += a[i]
	}
	fmt.Fscan(reader, &k)
	b = make([]int64, k)
	var sumb int64
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &b[i])
		sumb += b[i]
	}
	if suma != sumb {
		fmt.Println("NO")
		return
	}
	flag = make([]int, n)
	nxt = make([]int, n)
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
			fmt.Println("NO")
			return
		}
	}
	if cur != k {
		fmt.Println("NO")
		return
	}
	ans = make([]op, 0)
	var noFail bool
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
		fmt.Println("NO")
	} else {
		fmt.Println("YES")
		for _, t := range ans {
			fmt.Printf("%d %c\n", t.pos, t.typ)
		}
	}
}
