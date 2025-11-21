package main

import (
	"bufio"
	"fmt"
	"os"
)

func possible(a, b int) bool {
	if a == 0 && b == 0 {
		return true
	}
	if a == 0 || b == 0 {
		return a <= 2 && b <= 2
	}
	if a+b <= 6 {
		queue := [][]int{{}}
		for len(queue) > 0 {
			seq := queue[0]
			queue = queue[1:]
			if len(seq) == 6 {
				continue
			}
			for team := 0; team < 2; team++ {
				sz := len(seq)
				cnt := 0
				for i := sz - 1; i >= 0 && seq[i] == team; i-- {
					cnt++
				}
				if cnt < 2 {
					newSeq := append([]int(nil), seq...)
					newSeq = append(newSeq, team)
					ca, cb := 0, 0
					for _, x := range newSeq {
						if x == 0 {
							ca++
						} else {
							cb++
						}
					}
					if ca == a && cb == b {
						return true
					}
					queue = append(queue, newSeq)
				}
			}
		}
		return false
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var a, b, c, d int
		fmt.Fscan(in, &a, &b, &c, &d)
		if possible(a, b) && possible(c-a, d-b) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
