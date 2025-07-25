package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strconv"
	"strings"
)

func uniqueInts(a []int) []int {
	if len(a) == 0 {
		return a
	}
	j := 0
	for i := 1; i < len(a); i++ {
		if a[i] != a[j] {
			j++
			a[j] = a[i]
		}
	}
	return a[:j+1]
}

func intsToKey(a []int) string {
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	sort.Ints(arr)
	arr = uniqueInts(arr)
	n = len(arr)

	mp := make(map[string]int)
	for mask := 0; mask < (1 << 15); mask++ {
		d := make([]int, n)
		for i := 0; i < n; i++ {
			d[i] = bits.OnesCount(uint(arr[i]&0x7fff ^ mask))
		}
		base := d[0]
		for i := 0; i < n; i++ {
			d[i] -= base
		}
		key := intsToKey(d)
		mp[key] = mask
	}

	for mask := 0; mask < (1 << 15); mask++ {
		d := make([]int, n)
		for i := 0; i < n; i++ {
			d[i] = 30 - bits.OnesCount(uint(arr[i]>>15^mask))
		}
		base := d[0]
		for i := 0; i < n; i++ {
			d[i] -= base
		}
		key := intsToKey(d)
		if low, ok := mp[key]; ok {
			fmt.Fprintln(writer, (mask<<15)^low)
			return
		}
	}
	fmt.Fprintln(writer, -1)
}
