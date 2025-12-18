package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	scanInt := func() int {
		if !scanner.Scan() {
			return 0
		}
		x, _ := strconv.Atoi(scanner.Text())
		return x
	}

	// Read number of tests
	numTests := scanInt()

	for t := 0; t < numTests; t++ {
		n := scanInt()
		if n == 0 {
			break
		}
		s := make([]byte, n+1)
		s[1] = '0'

		// For n <= 120, querying every position uses n-1 queries <= 119 < 789 limit.
		for i := 2; i <= n; i++ {
			fmt.Printf("1 %d\n", i)
			
			resp := scanInt()
			if resp > 0 {
				s[i] = s[resp]
			} else {
				s[i] = '1'
			}
		}

		// Output answer
		fmt.Printf("0 %s\n", string(s[1:]))
		
		// Read result
		res := scanInt()
		if res != 1 {
			// Failed
			os.Exit(0)
		}
	}
}
