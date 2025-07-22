package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const maxK = 52

var (
	cost  [maxK][maxK]int
	u, v  [maxK]int
	dist  [maxK]int
	dad   [maxK]int
	seen  [maxK]bool
	Lmate [maxK]int
	Rmate [maxK]int
	K, N  int
)

func getId(ch byte) int {
	if ch >= 'a' && ch <= 'z' {
		return int(ch - 'a')
	}
	return int(ch - 'A' + 26)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func MinCostMatching() (int, []int) {
	for i := 0; i < K; i++ {
		u[i] = cost[i][0]
		for j := 1; j < K; j++ {
			if cost[i][j] < u[i] {
				u[i] = cost[i][j]
			}
		}
	}
	for j := 0; j < K; j++ {
		v[j] = cost[0][j] - u[0]
		for i := 1; i < K; i++ {
			if cost[i][j]-u[i] < v[j] {
				v[j] = cost[i][j] - u[i]
			}
		}
	}
	for i := 0; i < K; i++ {
		Lmate[i] = -1
	}
	for j := 0; j < K; j++ {
		Rmate[j] = -1
	}
	mated := 0
	for i := 0; i < K; i++ {
		for j := 0; j < K; j++ {
			if Rmate[j] != -1 {
				continue
			}
			if cost[i][j]-u[i]-v[j] == 0 {
				Lmate[i] = j
				Rmate[j] = i
				mated++
				break
			}
		}
	}
	for mated < K {
		var s int
		for s = 0; s < K; s++ {
			if Lmate[s] == -1 {
				break
			}
		}
		for k := 0; k < K; k++ {
			dad[k] = -1
			seen[k] = false
			dist[k] = cost[s][k] - u[s] - v[k]
		}
		var j int
		for {
			j = -1
			for k := 0; k < K; k++ {
				if seen[k] {
					continue
				}
				if j == -1 || dist[k] < dist[j] {
					j = k
				}
			}
			seen[j] = true
			if Rmate[j] == -1 {
				break
			}
			i := Rmate[j]
			for k := 0; k < K; k++ {
				if seen[k] {
					continue
				}
				newDist := dist[j] + cost[i][k] - u[i] - v[k]
				if dist[k] > newDist {
					dist[k] = newDist
					dad[k] = j
				}
			}
		}
		for k := 0; k < K; k++ {
			if k == j || !seen[k] {
				continue
			}
			i := Rmate[k]
			v[k] += dist[k] - dist[j]
			u[i] -= dist[k] - dist[j]
		}
		u[s] += dist[j]
		for pj := j; dad[pj] >= 0; pj = dad[pj] {
			d := dad[pj]
			Rmate[pj] = Rmate[d]
			Lmate[Rmate[pj]] = pj
		}
		Rmate[j] = s
		Lmate[s] = j
		mated++
	}
	value := 0
	for i := 0; i < K; i++ {
		value += cost[i][Lmate[i]]
	}
	match := make([]int, K)
	for i := 0; i < K; i++ {
		match[i] = Lmate[i]
	}
	return value, match
}

func expected(k int, s1, s2 string) int {
	K = k
	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			cost[i][j] = 0
		}
	}
	for i := 0; i < len(s1); i++ {
		a := getId(s1[i])
		b := getId(s2[i])
		cost[a][b]--
	}
	val, _ := MinCostMatching()
	return -val
}

func scoreFromMapping(mapping string, s1, s2 string) int {
	mp := make([]byte, len(mapping))
	copy(mp, mapping)
	count := 0
	for i := 0; i < len(s1); i++ {
		id := getId(s1[i])
		if id >= len(mp) {
			return -1
		}
		if mp[id] == s2[i] {
			count++
		}
	}
	return count
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const tests = 100
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for t := 0; t < tests; t++ {
		K = rand.Intn(10) + 1
		N = rand.Intn(20) + 1
		s1 := make([]byte, N)
		s2 := make([]byte, N)
		for i := 0; i < N; i++ {
			ch := letters[rand.Intn(K)]
			s1[i] = ch
			s2[i] = letters[rand.Intn(K)]
		}
		input := fmt.Sprintf("%d %d\n%s\n%s\n", N, K, string(s1), string(s2))
		want := expected(K, string(s1), string(s2))
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nOutput:\n%s\n", t+1, err, out)
			return
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) < 2 {
			fmt.Printf("Test %d invalid output: expected two lines, got %q\n", t+1, out)
			return
		}
		var gotVal int
		if _, err := fmt.Sscan(lines[0], &gotVal); err != nil {
			fmt.Printf("Test %d invalid score: %v\n", t+1, err)
			return
		}
		mapping := strings.TrimSpace(lines[1])
		if len(mapping) != K {
			fmt.Printf("Test %d invalid mapping length. Expected %d got %d\n", t+1, K, len(mapping))
			return
		}
		// check permutation
		seenLetters := make(map[byte]bool)
		for i := 0; i < K; i++ {
			ch := mapping[i]
			if ch < 'A' || (ch > 'Z' && ch < 'a') || ch > 'z' {
				fmt.Printf("Test %d invalid character in mapping: %c\n", t+1, ch)
				return
			}
			if getId(ch) >= K {
				fmt.Printf("Test %d mapping uses out-of-range letter %c\n", t+1, ch)
				return
			}
			if seenLetters[ch] {
				fmt.Printf("Test %d mapping not a permutation\n", t+1)
				return
			}
			seenLetters[ch] = true
		}
		score := scoreFromMapping(mapping, string(s1), string(s2))
		if score != gotVal || score != want {
			fmt.Printf("Test %d failed. Input:\n%sExpected: %d\nGot: %d %s\n", t+1, input, want, gotVal, mapping)
			return
		}
	}
	fmt.Println("All tests passed.")
}
