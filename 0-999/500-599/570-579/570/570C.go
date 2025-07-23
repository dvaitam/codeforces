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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)
	arr := []byte(s)

	pairs := 0
	for i := 0; i < n-1; i++ {
		if arr[i] == '.' && arr[i+1] == '.' {
			pairs++
		}
	}

	for i := 0; i < m; i++ {
		var x int
		var cs string
		fmt.Fscan(reader, &x, &cs)
		ch := cs[0]
		x--
		if arr[x] == ch {
			fmt.Fprintln(writer, pairs)
			continue
		}
		if arr[x] == '.' {
			if x > 0 && arr[x-1] == '.' {
				pairs--
			}
			if x+1 < n && arr[x+1] == '.' {
				pairs--
			}
		}
		arr[x] = ch
		if arr[x] == '.' {
			if x > 0 && arr[x-1] == '.' {
				pairs++
			}
			if x+1 < n && arr[x+1] == '.' {
				pairs++
			}
		}
		fmt.Fprintln(writer, pairs)
	}
}
