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

	var guest, host, pile string
	if _, err := fmt.Fscan(reader, &guest); err != nil {
		return
	}
	fmt.Fscan(reader, &host)
	fmt.Fscan(reader, &pile)

	// check total length
	if len(pile) != len(guest)+len(host) {
		fmt.Fprintln(writer, "NO")
		return
	}

	// count letters
	var cnt [26]int
	for _, c := range guest {
		cnt[c-'A']++
	}
	for _, c := range host {
		cnt[c-'A']++
	}
	for _, c := range pile {
		cnt[c-'A']--
	}

	// verify counts
	for _, v := range cnt {
		if v != 0 {
			fmt.Fprintln(writer, "NO")
			return
		}
	}

	fmt.Fprintln(writer, "YES")
}
