package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	n := len(s)

	// dpPal[i][j] == true if s[i:j+1] is a palindrome
	dpPal := make([][]bool, n)
	for i := range dpPal {
		dpPal[i] = make([]bool, n)
	}
	for i := n - 1; i >= 0; i-- {
		for j := i; j < n; j++ {
			if s[i] == s[j] {
				if j-i < 2 || dpPal[i+1][j-1] {
					dpPal[i][j] = true
				}
			}
		}
	}

	// characteristic value for each substring
	charVal := make([][]int, n)
	for i := range charVal {
		charVal[i] = make([]int, n)
	}
	freq := make([]int, n+2)

	for i := n - 1; i >= 0; i-- {
		for j := i; j < n; j++ {
			if dpPal[i][j] {
				if i == j {
					charVal[i][j] = 1
				} else {
					mid := i + (j-i)/2
					charVal[i][j] = 1 + charVal[i][mid]
				}
				freq[charVal[i][j]]++
			}
		}
	}

	for k := n - 1; k >= 1; k-- {
		freq[k] += freq[k+1]
	}

	writer := bufio.NewWriter(os.Stdout)
	for k := 1; k <= n; k++ {
		if k > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, freq[k])
	}
	writer.WriteByte('\n')
	writer.Flush()
}
