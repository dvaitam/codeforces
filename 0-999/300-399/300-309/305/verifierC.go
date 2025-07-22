package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Heap struct{ data []int64 }

func (h *Heap) Len() int { return len(h.data) }
func (h *Heap) Push(x int64) {
	h.data = append(h.data, x)
	i := len(h.data) - 1
	for i > 0 {
		p := (i - 1) / 2
		if h.data[p] <= h.data[i] {
			break
		}
		h.data[p], h.data[i] = h.data[i], h.data[p]
		i = p
	}
}
func (h *Heap) Pop() int64 {
	n := len(h.data)
	if n == 0 {
		return 0
	}
	ret := h.data[0]
	last := h.data[n-1]
	h.data = h.data[:n-1]
	if n-1 > 0 {
		h.data[0] = last
		i := 0
		for {
			left := 2*i + 1
			if left >= len(h.data) {
				break
			}
			minChild := left
			right := left + 1
			if right < len(h.data) && h.data[right] < h.data[left] {
				minChild = right
			}
			if h.data[minChild] >= h.data[i] {
				break
			}
			h.data[i], h.data[minChild] = h.data[minChild], h.data[i]
			i = minChild
		}
	}
	return ret
}

func solveC(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	cnt := make(map[int64]int64, n)
	var ai int64
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &ai)
		cnt[ai]++
	}
	h := &Heap{}
	inHeap := make(map[int64]bool)
	for k, v := range cnt {
		if v >= 2 {
			h.Push(k)
			inHeap[k] = true
		}
	}
	for h.Len() > 0 {
		k := h.Pop()
		inHeap[k] = false
		v := cnt[k]
		if v < 2 {
			continue
		}
		carry := v >> 1
		cnt[k] = v & 1
		nk := k + 1
		cnt[nk] += carry
		if cnt[nk] >= 2 && !inHeap[nk] {
			h.Push(nk)
			inHeap[nk] = true
		}
	}
	var H int64 = -1
	var ones int64
	for k, v := range cnt {
		if v > 0 {
			if k > H {
				H = k
			}
			ones++
		}
	}
	ans := (H + 1) - ones
	return fmt.Sprint(ans)
}

func genTestC() string {
	n := rand.Intn(20) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rand.Intn(20))
	}
	var buf bytes.Buffer
	fmt.Fprintln(&buf, n)
	for i, v := range arr {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprint(&buf, v)
	}
	buf.WriteByte('\n')
	return buf.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		in := genTestC()
		expected := solveC(in)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(in)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\noutput: %s\n", i, err, string(out))
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
