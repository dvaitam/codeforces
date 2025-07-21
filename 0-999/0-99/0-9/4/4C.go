package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	cnt := make(map[string]int)
	for i := 0; i < n; i++ {
		var name string
		fmt.Fscan(reader, &name)
		if cnt[name] == 0 {
			fmt.Fprintln(writer, "OK")
			cnt[name] = 1
		} else {
			res := name + strconv.Itoa(cnt[name])
			fmt.Fprintln(writer, res)
			cnt[name]++
		}
	}
}
