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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	// read initial states
	states := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &states[i])
	}
	// read key presses
	var k int
	fmt.Fscan(reader, &k)
	keys := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &keys[i])
	}
	// map two state words
	wordIdx := make(map[string]int, 2)
	words := make([]string, 0, 2)
	init := make([]bool, n)
	for i, s := range states {
		if _, ok := wordIdx[s]; !ok && len(words) < 2 {
			wordIdx[s] = len(words)
			words = append(words, s)
		}
		// default map existing words to bool by index (0->false, 1->true)
		init[i] = wordIdx[s] == 1
	}
	// prepare toggle parity for keys
	keyParity := make([]bool, n+1)
	for _, x := range keys {
		if x >= 1 && x <= n {
			keyParity[x] = !keyParity[x]
		}
	}
	// apply toggles for each key with odd count
	for i := 1; i <= n; i++ {
		if !keyParity[i] {
			continue
		}
		for j := i; j <= n; j += i {
			init[j-1] = !init[j-1]
		}
	}
	// output final states
	for i, st := range init {
		var out string
		if st {
			// index 1 word
			if len(words) > 1 {
				out = words[1]
			} else {
				out = words[0]
			}
		} else {
			out = words[0]
		}
		writer.WriteString(out)
		if i+1 < n {
			writer.WriteByte(' ')
		}
	}
	writer.WriteByte('\n')
}
