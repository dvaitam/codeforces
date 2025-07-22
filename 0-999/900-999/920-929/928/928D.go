package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type node struct {
	child [26]int
	cnt   int
}

var trie []node

func insert(word string) {
	idx := 0
	for i := 0; i < len(word); i++ {
		c := int(word[i] - 'a')
		if trie[idx].child[c] == 0 {
			trie[idx].child[c] = len(trie)
			trie = append(trie, node{})
		}
		idx = trie[idx].child[c]
		trie[idx].cnt++
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	data, _ := io.ReadAll(reader)
	s := string(data)

	trie = make([]node, 1)
	seen := make(map[string]bool)

	var total int64
	n := len(s)
	for i := 0; i < n; {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			j := i
			for j < n && s[j] >= 'a' && s[j] <= 'z' {
				j++
			}
			word := s[i:j]
			cost := len(word)
			if seen[word] {
				idx := 0
				for k := 0; k < len(word); k++ {
					ch := int(word[k] - 'a')
					idx = trie[idx].child[ch]
					if trie[idx].cnt == 1 {
						if k+2 < cost {
							cost = k + 2
						}
						break
					}
				}
			}
			total += int64(cost)
			if !seen[word] {
				seen[word] = true
				insert(word)
			}
			i = j
		} else {
			total++
			i++
		}
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, total)
	out.Flush()
}
