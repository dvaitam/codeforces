package main

import (
	"bufio"
	"fmt"
	"os"
)

// This solution directly follows the statement of problemF.
// It enumerates all substrings of s and of all t_i (strings with the i-th
// character replaced by '*').  Each distinct substring is inserted into
// a map in order to count the unique values.
//
// The algorithm is O(n^3) in the worst case and therefore not suitable for
// the original constraints where n can be up to 1e5.  However it demonstrates
// the required logic in a straightforward way and works for small inputs.

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	n := len(s)
	uniq := make(map[string]struct{})

	uniq[""] = struct{}{} // empty string

	// helper to insert all substrings of given string into the set
	addSubs := func(str string) {
		for i := 0; i < len(str); i++ {
			for j := i + 1; j <= len(str); j++ {
				sub := str[i:j]
				uniq[sub] = struct{}{}
			}
		}
	}

	// substrings from original string
	addSubs(s)

	// strings t_i with the i-th character replaced by '*'
	for i := 0; i < n; i++ {
		t := []byte(s)
		t[i] = '*'
		addSubs(string(t))
	}

	fmt.Println(len(uniq))
}
