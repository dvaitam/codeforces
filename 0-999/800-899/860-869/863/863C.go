package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var k int64
	var a, b int
	if _, err := fmt.Fscan(reader, &k, &a, &b); err != nil {
		return
	}
	nextA := make([][]int, 4)
	nextB := make([][]int, 4)
	for i := 1; i <= 3; i++ {
		nextA[i] = make([]int, 4)
		for j := 1; j <= 3; j++ {
			fmt.Fscan(reader, &nextA[i][j])
		}
	}
	for i := 1; i <= 3; i++ {
		nextB[i] = make([]int, 4)
		for j := 1; j <= 3; j++ {
			fmt.Fscan(reader, &nextB[i][j])
		}
	}

	visitedStep := [4][4]int64{}
	visitedA := [4][4]int64{}
	visitedB := [4][4]int64{}
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			visitedStep[i][j] = -1
		}
	}

	var scoreA, scoreB int64
	var step int64
	curA, curB := a, b

	for step < k {
		if visitedStep[curA][curB] != -1 {
			prev := visitedStep[curA][curB]
			cycleLen := step - prev
			cycleScoreA := scoreA - visitedA[curA][curB]
			cycleScoreB := scoreB - visitedB[curA][curB]
			if cycleLen > 0 {
				times := (k - step) / cycleLen
				if times > 0 {
					scoreA += cycleScoreA * times
					scoreB += cycleScoreB * times
					step += cycleLen * times
				}
			}
		}
		if step >= k {
			break
		}
		visitedStep[curA][curB] = step
		visitedA[curA][curB] = scoreA
		visitedB[curA][curB] = scoreB

		if curA != curB {
			if curA == 1 && curB == 3 || curA == 2 && curB == 1 || curA == 3 && curB == 2 {
				scoreA++
			} else {
				scoreB++
			}
		}
		step++
		na := nextA[curA][curB]
		nb := nextB[curA][curB]
		curA, curB = na, nb
	}

	fmt.Printf("%d %d\n", scoreA, scoreB)
}
