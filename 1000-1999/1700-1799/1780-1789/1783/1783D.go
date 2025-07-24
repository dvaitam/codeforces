package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const mod int = 998244353

func serialize(arr []int) string {
	// join numbers with comma
	sb := strings.Builder{}
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func deserialize(s string) []int {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	res := make([]int, len(parts))
	for i, p := range parts {
		fmt.Sscanf(p, "%d", &res[i])
	}
	return res
}

func countReachable(a []int) int {
	states := map[string]struct{}{serialize(a): {}}
	n := len(a)
	for i := 1; i < n-1; i++ {
		next := make(map[string]struct{})
		for s := range states {
			arr := deserialize(s)
			v := arr[i]
			if v == 0 {
				next[s] = struct{}{}
				continue
			}
			arr1 := append([]int(nil), arr...)
			arr1[i-1] += v
			arr1[i+1] -= v
			next[serialize(arr1)] = struct{}{}

			arr2 := append([]int(nil), arr...)
			arr2[i-1] -= v
			arr2[i+1] += v
			next[serialize(arr2)] = struct{}{}
		}
		states = next
	}
	return len(states) % mod
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	res := countReachable(a)
	fmt.Println(res)
}
