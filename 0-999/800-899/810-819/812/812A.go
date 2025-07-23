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

	var a [4][4]int
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			fmt.Fscan(reader, &a[i][j])
		}
	}

	for i := 0; i < 4; i++ {
		if a[i][3] == 1 {
			if a[i][0] == 1 || a[i][1] == 1 || a[i][2] == 1 ||
				a[(i+1)%4][2] == 1 || a[(i+2)%4][1] == 1 || a[(i+3)%4][0] == 1 {
				fmt.Fprintln(writer, "YES")
				return
			}
		}
	}
	fmt.Fprintln(writer, "NO")
}
