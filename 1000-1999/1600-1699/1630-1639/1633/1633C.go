package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemC.txt for 1633C.
// It checks if the hero can defeat the monster after distributing up to k coins
// between weapon upgrades (increasing attack by w each) and armor upgrades
// (increasing health by a each). We test every possible distribution and
// compare the number of turns required for the hero and the monster to win.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var hC, dC int64
		fmt.Fscan(in, &hC, &dC)
		var hM, dM int64
		fmt.Fscan(in, &hM, &dM)
		var k, w, a int64
		fmt.Fscan(in, &k, &w, &a)

		win := false
		for i := int64(0); i <= k; i++ {
			attack := dC + i*w
			health := hC + (k-i)*a
			turnsHero := (hM + attack - 1) / attack
			turnsMonster := (health + dM - 1) / dM
			if turnsHero <= turnsMonster {
				win = true
				break
			}
		}
		if win {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
