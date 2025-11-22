package main

import (
	"bufio"
	"fmt"
	"os"
)

type participant struct {
	a int64
	b int64
	s int64
}

type problem struct {
	d int64
	t int64
}

const (
	inf    int64 = 4e18
	maxVal       = int64(1_000_000_000)
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		parts := make([]participant, n)
		specSet := make(map[int64]bool)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &parts[i].a, &parts[i].b, &parts[i].s)
			specSet[parts[i].s] = true
		}

		cities := make([][]int, m)
		for i := 0; i < m; i++ {
			var k int
			fmt.Fscan(in, &k)
			cities[i] = make([]int, k)
			for j := 0; j < k; j++ {
				fmt.Fscan(in, &cities[i][j])
				cities[i][j]--
			}
		}

		// prefix data
		prefMin1 := make([]int64, m)
		prefMin2 := make([]int64, m)
		prefCnt1 := make([]int, m)
		prefCntTopic1 := make([]int, m)
		prefTopic1 := make([]int64, m)
		prefMinBTopic := make([]int64, m)

		min1, min2 := inf, inf
		cnt1, cntTopic1 := 0, 0
		topic1 := int64(-1)
		minBTopic := inf
		minBMap := make(map[int64]int64)

		for idx := 0; idx < m; idx++ {
			for _, id := range cities[idx] {
				p := parts[id]
				if v, ok := minBMap[p.s]; !ok || p.b < v {
					minBMap[p.s] = p.b
				}

				if p.a < min1 {
					min2 = min1
					min1 = p.a
					cnt1 = 1
					topic1 = p.s
					cntTopic1 = 1
				} else if p.a == min1 {
					cnt1++
					if p.s == topic1 {
						cntTopic1++
					}
				} else if p.a < min2 {
					min2 = p.a
				}
			}

			if topic1 != -1 {
				minBTopic = minBMap[topic1]
			}

			prefMin1[idx] = min1
			prefMin2[idx] = min2
			prefCnt1[idx] = cnt1
			prefCntTopic1[idx] = cntTopic1
			prefTopic1[idx] = topic1
			prefMinBTopic[idx] = minBTopic
		}

		// suffix data
		suffixMaxA := int64(-1)
		suffixMaxB := make(map[int64]int64)
		for _, id := range cities[m-1] {
			p := parts[id]
			if p.a > suffixMaxA {
				suffixMaxA = p.a
			}
			if v, ok := suffixMaxB[p.s]; !ok || p.b > v {
				suffixMaxB[p.s] = p.b
			}
		}

		usedTopics := make(map[int64]bool)
		problems := make([]problem, 0, m-1)
		newTopic := maxVal
		fail := false

		for i := m - 2; i >= 0; i-- {
			boundA := suffixMaxA
			minA := prefMin1[i]
			found := false
			var chosenTopic int64
			var diff int64

			if minA > boundA {
				for specSet[newTopic] || usedTopics[newTopic] {
					newTopic--
				}
				chosenTopic = newTopic
				diff = boundA + 1
				if diff > maxVal {
					fail = true
					break
				}
				found = true
				usedTopics[chosenTopic] = true
				newTopic--
			} else if prefCnt1[i] == prefCntTopic1[i] {
				topic := prefTopic1[i]
				if !usedTopics[topic] {
					bound := boundA
					if v, ok := suffixMaxB[topic]; ok && v > bound {
						bound = v
					}
					capacity := prefMin2[i]
					if capacity > prefMinBTopic[i] {
						capacity = prefMinBTopic[i]
					}
					if capacity > bound {
						diff = bound + 1
						if diff > maxVal {
							fail = true
							break
						}
						found = true
						chosenTopic = topic
						usedTopics[chosenTopic] = true
					}
				}
			}

			if !found {
				fail = true
				break
			}

			problems = append(problems, problem{d: diff, t: chosenTopic})

			for _, id := range cities[i] {
				p := parts[id]
				if p.a > suffixMaxA {
					suffixMaxA = p.a
				}
				if v, ok := suffixMaxB[p.s]; !ok || p.b > v {
					suffixMaxB[p.s] = p.b
				}
			}
		}

		if fail {
			fmt.Fprintln(out, -1)
			continue
		}

		fmt.Fprintln(out, len(problems))
		for _, pr := range problems {
			fmt.Fprintln(out, pr.d, pr.t)
		}
	}
}
