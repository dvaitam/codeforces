package main

import (
	"bufio"
	"fmt"
	"os"
)

func minMoves(s string, x, y byte) int {
	arr := []byte(s)
	n := len(arr)
	const INF = int(1e9)
	best := INF
	for j := 0; j < n; j++ {
		if arr[j] != y {
			continue
		}
		// move digit y at position j to the end
		costY := n - 1 - j
		arr1 := append(append([]byte{}, arr[:j]...), arr[j+1:]...)
		arr1 = append(arr1, y)
		for i := 0; i < n-1; i++ {
			if arr1[i] != x {
				continue
			}
			// move digit x at position i to the second last place
			costX := n - 2 - i
			arr2 := append(append([]byte{}, arr1[:i]...), arr1[i+1:n-1]...)
			arr2 = append(arr2, x, y)

			idx := -1
			for k := 0; k < n-2; k++ {
				if arr2[k] != '0' {
					idx = k
					break
				}
			}
			if idx == -1 {
				if n == 2 {
					if arr2[0] == '0' {
						continue
					}
					idx = 0
				} else {
					continue
				}
			}
			total := costY + costX + idx
			if total < best {
				best = total
			}
		}
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	const INF = int(1e9)
	ans := INF
	pairs := [][2]byte{{'0', '0'}, {'2', '5'}, {'5', '0'}, {'7', '5'}}
	for _, p := range pairs {
		mv := minMoves(s, p[0], p[1])
		if mv < ans {
			ans = mv
		}
	}
	if ans == INF {
		fmt.Println(-1)
	} else {
		fmt.Println(ans)
	}
}
