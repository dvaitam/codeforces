package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF = int(1e9)

func isWinner(votes []int, last []int, i, k int) bool {
	type cand struct{ id, v, t int }
	arr := make([]cand, 0)
	n := len(votes) - 1
	for j := 1; j <= n; j++ {
		if votes[j] > 0 {
			arr = append(arr, cand{j, votes[j], last[j]})
		}
	}
	sort.Slice(arr, func(a, b int) bool {
		if arr[a].v != arr[b].v {
			return arr[a].v > arr[b].v
		}
		if arr[a].t != arr[b].t {
			return arr[a].t < arr[b].t
		}
		return arr[a].id < arr[b].id
	})
	if k > len(arr) {
		k = len(arr)
	}
	for idx := 0; idx < k; idx++ {
		if arr[idx].id == i {
			return true
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k, m, a int
	if _, err := fmt.Fscan(reader, &n, &k, &m, &a); err != nil {
		return
	}
	votes := make([]int, n+1)
	last := make([]int, n+1)
	for i := 1; i <= n; i++ {
		last[i] = INF
	}
	for j := 1; j <= a; j++ {
		var g int
		fmt.Fscan(reader, &g)
		votes[g]++
		last[g] = j
	}
	rem := m - a
	results := make([]int, n+1)
	for i := 1; i <= n; i++ {
		// best case: give all remaining votes to candidate i
		bestVotes := make([]int, n+1)
		copy(bestVotes, votes)
		bestLast := make([]int, n+1)
		copy(bestLast, last)
		if rem > 0 {
			bestVotes[i] += rem
			bestLast[i] = a + rem
		}
		if isWinner(bestVotes, bestLast, i, k) {
			// candidate has a chance
			ci := votes[i]
			ti := last[i]
			if ci == 0 {
				ti = INF
			}
			costs := make([]int, 0, n-1)
			for j := 1; j <= n; j++ {
				if j == i {
					continue
				}
				cj := votes[j]
				tj := last[j]
				if cj == 0 {
					tj = INF
				}
				var cost int
				if cj > ci {
					cost = 0
				} else if cj == ci {
					if cj == 0 {
						cost = 1
					} else {
						if tj < ti {
							cost = 0
						} else {
							cost = 1
						}
					}
				} else {
					cost = ci - cj + 1
				}
				costs = append(costs, cost)
			}
			sort.Ints(costs)
			left := rem
			count := 0
			for _, c := range costs {
				if c <= left {
					left -= c
					count++
				} else {
					break
				}
			}
			if ci == 0 {
				// adversary can leave candidate without votes
				if rem == 0 {
					results[i] = 3
				} else {
					results[i] = 2
				}
			} else if count < k {
				results[i] = 1
			} else {
				results[i] = 2
			}
		} else {
			results[i] = 3
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, results[i])
	}
	fmt.Fprintln(writer)
}
