package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)

	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		if isConsistent(a) {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
		}
	}
}

func isConsistent(a []int) bool {
	n := len(a)
	// state: 0 = no win yet, 1 = already has win (either earlier or claimed)
	win := make([]int, n)

	// For each match between i and i+1, check if there is a possible direction consistent with claims.
	for i := 0; i < n-1; i++ {
		// If both claim no win and neither previously got a win, then impossible
		if win[i] == 0 && win[i+1] == 0 && a[i] == 0 && a[i+1] == 0 {
			return false
		}
		// Otherwise, at least one must win here. Assign win to someone consistent with claim.
		if win[i] == 0 && a[i] == 0 {
			// i cannot win anymore -> i+1 must win
			win[i+1] = 1
		} else if win[i+1] == 0 && a[i+1] == 0 {
			// i+1 cannot win -> i must win
			win[i] = 1
		} else {
			// if both can potentially win, choose one arbitrarily; prefer giving win to whoever still needs to fulfill a claim
			if a[i] == 1 && win[i] == 0 {
				win[i] = 1
			} else {
				win[i+1] = 1
			}
		}
	}

	// verify final claims: if a player claimed 1 but never got win, impossible
	for i := 0; i < n; i++ {
		if a[i] == 1 && win[i] == 0 {
			return false
		}
	}
	return true
}
