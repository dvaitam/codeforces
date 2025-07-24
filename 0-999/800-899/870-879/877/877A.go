package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	names := []string{"Danil", "Olya", "Slava", "Ann", "Nikita"}
	count := 0

	for i := 0; i < len(s); i++ {
		for _, name := range names {
			if strings.HasPrefix(s[i:], name) {
				count++
			}
		}
	}

	if count == 1 {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}
