package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	var M int64
	if _, err := fmt.Fscan(reader, &n, &k, &M); err != nil {
		return
	}
	t := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &t[i])
	}
	sort.Ints(t)
	prefix := make([]int64, k+1)
	for i := 0; i < k; i++ {
		prefix[i+1] = prefix[i] + int64(t[i])
	}
	fullTaskTime := prefix[k]

	ans := 0
	for full := 0; full <= n; full++ {
		spent := int64(full) * fullTaskTime
		if spent > M {
			break
		}
		points := full * (k + 1)
		left := M - spent
		tasksAvail := n - full
		curTasks := tasksAvail
		for j := 0; j < k; j++ {
			if curTasks == 0 || left < int64(t[j]) {
				break
			}
			// number of tasks we can solve j-th subtask for
			maxTasks := int(left / int64(t[j]))
			if maxTasks > curTasks {
				maxTasks = curTasks
			}
			points += maxTasks
			left -= int64(maxTasks) * int64(t[j])
			curTasks = maxTasks
		}
		if points > ans {
			ans = points
		}
	}
	fmt.Fprintln(writer, ans)
}
