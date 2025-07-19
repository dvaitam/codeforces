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
	fmt.Fscan(reader, &n)
	if n%2 == 0 {
		writer.WriteString("NO\n")
		return
	}
	writer.WriteString("YES\n")
	for i := 1; i <= n; i++ {
		x := i*2 - i%2
		writer.WriteString(strconv.Itoa(x))
		writer.WriteByte(' ')
	}
	for i := 1; i <= n; i++ {
		x := i*2 + i%2 - 1
		writer.WriteString(strconv.Itoa(x))
		writer.WriteByte(' ')
	}
}
