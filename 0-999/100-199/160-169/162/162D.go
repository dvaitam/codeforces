package main

import (
	"fmt"
)

func main() {
   var s string
   // Read the input string (no spaces as per problem constraints)
   if _, err := fmt.Scan(&s); err != nil {
		return
   }
   // Build result without digits
   res := make([]rune, 0, len(s))
   for _, r := range s {
	if r < '0' || r > '9' {
		res = append(res, r)
	}
   }
   // Output the filtered string
   fmt.Println(string(res))
}
