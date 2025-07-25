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

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		monsters := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &monsters[i])
		}
		var m int
		fmt.Fscan(reader, &m)
		// For each endurance value, store max power of heroes with at least that endurance
		maxPower := make([]int, n+2) // index by endurance
		for i := 0; i < m; i++ {
			var p, s int
			fmt.Fscan(reader, &p, &s)
			if maxPower[s] < p {
				maxPower[s] = p
			}
		}
		// propagate maxima backwards so that maxPower[s] holds max power for >=s
		for i := n - 1; i >= 1; i-- {
			if maxPower[i] < maxPower[i+1] {
				maxPower[i] = maxPower[i+1]
			}
		}
		// check if each monster can be beaten individually
		possible := true
		for i := 0; i < n; i++ {
			if monsters[i] > maxPower[1] {
				possible = false
				break
			}
		}
		if !possible {
			fmt.Fprintln(writer, -1)
			continue
		}
		days := 0
		i := 0
		for i < n {
			days++
			maxA := monsters[i]
			j := i
			for j < n {
				if monsters[j] > maxA {
					maxA = monsters[j]
				}
				length := j - i + 1
				if maxPower[length] < maxA {
					break
				}
				j++
			}
			i = j
		}
		fmt.Fprintln(writer, days)
	}
}
