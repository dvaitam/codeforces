package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type problem struct {
	t   int64
	typ int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var cases int
	fmt.Fscan(reader, &cases)
	for ; cases > 0; cases-- {
		var n int
		var T, a, b int64
		fmt.Fscan(reader, &n, &T, &a, &b)

		tp := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &tp[i])
		}
		ti := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &ti[i])
		}

		tasks := make([]problem, n)
		totalEasy := 0
		totalHard := 0
		for i := 0; i < n; i++ {
			tasks[i] = problem{t: ti[i], typ: tp[i]}
			if tp[i] == 0 {
				totalEasy++
			} else {
				totalHard++
			}
		}
		sort.Slice(tasks, func(i, j int) bool {
			if tasks[i].t == tasks[j].t {
				return tasks[i].typ < tasks[j].typ
			}
			return tasks[i].t < tasks[j].t
		})
		// append sentinel
		tasks = append(tasks, problem{t: T + 1})

		prefixEasy, prefixHard := 0, 0
		ans := 0
		idx := 0
		m := len(tasks)
		for idx < m {
			curr := tasks[idx].t
			candidate := curr - 1
			if candidate >= 0 {
				spent := int64(prefixEasy)*a + int64(prefixHard)*b
				if spent <= candidate && candidate <= T {
					remain := candidate - spent
					remEasy := totalEasy - prefixEasy
					remHard := totalHard - prefixHard
					addEasy := int(minInt64(int64(remEasy), remain/a))
					remain -= int64(addEasy) * a
					addHard := int(minInt64(int64(remHard), remain/b))
					totalSolved := prefixEasy + prefixHard + addEasy + addHard
					if totalSolved > ans {
						ans = totalSolved
					}
				}
			}
			for idx < m && tasks[idx].t == curr {
				if tasks[idx].typ == 0 {
					prefixEasy++
				} else {
					prefixHard++
				}
				idx++
			}
		}

		fmt.Fprintln(writer, ans)
	}
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
