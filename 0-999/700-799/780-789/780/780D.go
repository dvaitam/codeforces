package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	first := make([]string, n)
	second := make([]string, n)
	groups := make(map[string][]int)

	for i := 0; i < n; i++ {
		var team, city string
		fmt.Fscan(reader, &team, &city)
		f1 := team[:3]
		f2 := team[:2] + city[:1]
		first[i] = f1
		second[i] = f2
		groups[f1] = append(groups[f1], i)
	}

	useSecond := make([]bool, n)
	queue := make([]string, 0)

	for _, idxs := range groups {
		if len(idxs) > 1 {
			for _, id := range idxs {
				if !useSecond[id] {
					useSecond[id] = true
					queue = append(queue, second[id])
				}
			}
		}
	}

	for len(queue) > 0 {
		x := queue[0]
		queue = queue[1:]
		if idxs, ok := groups[x]; ok {
			for _, id := range idxs {
				if !useSecond[id] {
					useSecond[id] = true
					queue = append(queue, second[id])
				}
			}
		}
	}

	res := make([]string, n)
	used := make(map[string]bool)
	for i := 0; i < n; i++ {
		if useSecond[i] {
			res[i] = second[i]
		} else {
			res[i] = first[i]
		}
		if used[res[i]] {
			fmt.Fprintln(writer, "NO")
			return
		}
		used[res[i]] = true
	}

	fmt.Fprintln(writer, "YES")
	for i := 0; i < n; i++ {
		fmt.Fprintln(writer, res[i])
	}
}
