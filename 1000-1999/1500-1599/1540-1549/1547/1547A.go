package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var xA, yA int
		var xB, yB int
		var xF, yF int
		fmt.Fscan(reader, &xA, &yA)
		fmt.Fscan(reader, &xB, &yB)
		fmt.Fscan(reader, &xF, &yF)

		dist := abs(xA-xB) + abs(yA-yB)
		if xA == xB && xA == xF {
			minY := yA
			maxY := yB
			if minY > maxY {
				minY, maxY = maxY, minY
			}
			if yF > minY && yF < maxY {
				dist += 2
			}
		} else if yA == yB && yA == yF {
			minX := xA
			maxX := xB
			if minX > maxX {
				minX, maxX = maxX, minX
			}
			if xF > minX && xF < maxX {
				dist += 2
			}
		}
		fmt.Fprintln(writer, dist)
	}
}
