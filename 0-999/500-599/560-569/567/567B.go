package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	inside := make(map[int]bool)
	cur := 0
	ans := 0
	for i := 0; i < n; i++ {
		var op string
		var id int
		fmt.Fscan(reader, &op, &id)
		if op == "+" {
			inside[id] = true
			cur++
			if cur > ans {
				ans = cur
			}
		} else {
			if inside[id] {
				delete(inside, id)
				cur--
			} else {
				ans++
			}
		}
	}
	fmt.Println(ans)
}
