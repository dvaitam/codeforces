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

	var n int
	var m, s, d int
	if _, err := fmt.Fscan(reader, &n, &m, &s, &d); err != nil {
		return
	}
	obstacles := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &obstacles[i])
	}
	sort.Ints(obstacles)

	if n > 0 && obstacles[0] <= s {
		fmt.Fprintln(writer, "IMPOSSIBLE")
		return
	}

	type cmd struct {
		t   string
		len int
	}
	var res []cmd
	pos := 0
	i := 0
	for i < n {
		if obstacles[i]-pos-1 < s {
			fmt.Fprintln(writer, "IMPOSSIBLE")
			return
		}
		runLen := obstacles[i] - 1 - pos
		if runLen > 0 {
			res = append(res, cmd{"RUN", runLen})
		}
		j := i
		for j+1 < n && obstacles[j+1]-obstacles[j]-2 < s {
			j++
		}
		jumpLen := obstacles[j] - obstacles[i] + 2
		if jumpLen > d || obstacles[i]-1+jumpLen > m {
			fmt.Fprintln(writer, "IMPOSSIBLE")
			return
		}
		res = append(res, cmd{"JUMP", jumpLen})
		pos = obstacles[j] + 1
		i = j + 1
	}
	if pos < m {
		res = append(res, cmd{"RUN", m - pos})
	}
	for _, c := range res {
		fmt.Fprintln(writer, c.t, c.len)
	}
}
