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

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genTest() []byte {
	n := rand.Intn(20) + 1 // 1..20
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('a' + rand.Intn(26)))
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

// Embedded correct solver for 1789F
var dp2 [81][81]int
var dp3 [81][81][81]int8

type TrieNode struct {
	child   int32
	sibling int32
	char    byte
	visited bool
}

var trie []TrieNode

func getChild(u int32, ch byte) int32 {
	c := trie[u].child
	var prev int32 = -1
	for c != 0 {
		if trie[c].char == ch {
			return c
		}
		prev = c
		c = trie[c].sibling
	}
	v := int32(len(trie))
	trie = append(trie, TrieNode{char: ch})
	if prev == -1 {
		trie[u].child = v
	} else {
		trie[prev].sibling = v
	}
	return v
}

func countMatches(sub []byte, s string) int {
	count := 0
	i := 0
	n := len(s)
	m := len(sub)
	for i < n {
		j := 0
		for i < n && j < m {
			if s[i] == sub[j] {
				j++
			}
			i++
		}
		if j == m {
			count++
		}
	}
	return count
}

func solve1789F(input string) string {
	S := strings.TrimSpace(input)
	n := len(S)
	if n <= 1 {
		return "0"
	}

	ans := 0

	for i := 1; i < n; i++ {
		s1 := S[:i]
		s2 := S[i:]
		n1, n2 := len(s1), len(s2)
		for a := 0; a <= n1; a++ {
			dp2[a][0] = 0
		}
		for b := 0; b <= n2; b++ {
			dp2[0][b] = 0
		}
		for a := 1; a <= n1; a++ {
			for b := 1; b <= n2; b++ {
				if s1[a-1] == s2[b-1] {
					dp2[a][b] = dp2[a-1][b-1] + 1
				} else {
					mx := dp2[a-1][b]
					if dp2[a][b-1] > mx {
						mx = dp2[a][b-1]
					}
					dp2[a][b] = mx
				}
			}
		}
		l := dp2[n1][n2]
		if l*2 > ans {
			ans = l * 2
		}
	}

	for i := 1; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			s1 := S[:i]
			s2 := S[i:j]
			s3 := S[j:]
			n1, n2, n3 := len(s1), len(s2), len(s3)
			for a := 0; a <= n1; a++ {
				for b := 0; b <= n2; b++ {
					dp3[a][b][0] = 0
				}
			}
			for a := 0; a <= n1; a++ {
				for c := 0; c <= n3; c++ {
					dp3[a][0][c] = 0
				}
			}
			for b := 0; b <= n2; b++ {
				for c := 0; c <= n3; c++ {
					dp3[0][b][c] = 0
				}
			}
			for a := 1; a <= n1; a++ {
				for b := 1; b <= n2; b++ {
					for c := 1; c <= n3; c++ {
						if s1[a-1] == s2[b-1] && s2[b-1] == s3[c-1] {
							dp3[a][b][c] = dp3[a-1][b-1][c-1] + 1
						} else {
							mx := dp3[a-1][b][c]
							if dp3[a][b-1][c] > mx {
								mx = dp3[a][b-1][c]
							}
							if dp3[a][b][c-1] > mx {
								mx = dp3[a][b][c-1]
							}
							dp3[a][b][c] = mx
						}
					}
				}
			}
			l := int(dp3[n1][n2][n3])
			if l*3 > ans {
				ans = l * 3
			}
		}
	}

	WLen := 16
	if n < 16 {
		WLen = n
	}
	trie = make([]TrieNode, 1, 100000)
	var buf [20]byte

	for start := 0; start <= n-WLen; start++ {
		w := S[start : start+WLen]
		wLen := len(w)
		var nextOcc [17][26]int
		for c := 0; c < 26; c++ {
			nextOcc[wLen][c] = -1
		}
		for i := wLen - 1; i >= 0; i-- {
			for c := 0; c < 26; c++ {
				nextOcc[i][c] = nextOcc[i+1][c]
			}
			nextOcc[i][w[i]-'a'] = i
		}

		var dfs func(wIndex int, trieIdx int32, subLen int)
		dfs = func(wIndex int, trieIdx int32, subLen int) {
			if !trie[trieIdx].visited && subLen > 0 {
				trie[trieIdx].visited = true
				c := countMatches(buf[:subLen], S)
				if c >= 2 {
					if c*subLen > ans {
						ans = c * subLen
					}
				}
			}
			for ch := byte(0); ch < 26; ch++ {
				i := nextOcc[wIndex][ch]
				if i != -1 {
					v := getChild(trieIdx, ch)
					buf[subLen] = 'a' + ch
					dfs(i+1, v, subLen+1)
				}
			}
		}
		dfs(0, 0, 0)
	}

	return fmt.Sprint(ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genTest()
		want := solve1789F(string(input))
		got, err := run(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\n", i+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("expected:\n", want)
			fmt.Println("got:\n", got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
