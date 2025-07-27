package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for Codeforces problem 1480B - The Great Hero.
// The hero with attack power A and health B fights n monsters with attack a_i
// and health b_i. In each round, the hero chooses a living monster and they
// simultaneously deal damage. The goal is to determine if the hero can kill all
// monsters (the hero may die after killing the last one).
//
// For monster i, it requires ceil(b_i/A) rounds to defeat it, and in each round
// the hero receives a_i damage. Hence the total damage taken from monster i is
// a_i * ceil(b_i/A). The order of fights only affects which monster delivers the
// final blow, so the hero should face the monster with the largest attack last.
// The hero must have strictly more health than the damage from all other
// monsters. This translates to B > sum(a_i * ceil(b_i/A)) - max(a_i).
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var A, B int64
		var n int
		fmt.Fscan(reader, &A, &B, &n)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}

		var total int64
		var maxAttack int64
		for i := 0; i < n; i++ {
			hits := (b[i] + A - 1) / A
			total += hits * a[i]
			if a[i] > maxAttack {
				maxAttack = a[i]
			}
		}

		if B > total-maxAttack {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
