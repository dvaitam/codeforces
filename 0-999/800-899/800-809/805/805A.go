package main

import (
   "fmt"
)

func main() {
   var a, b int64
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   if a == b {
       fmt.Println(a)
       return
   }
	// if range is larger than 10, always choose 2
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	if diff > 10 {
		fmt.Println(2)
		return
	}
   // ensure l <= r
   l, r := a, b
   if l > r {
       l, r = r, l
   }
   var sa, sb int64
   for i := l; i <= r; i++ {
       if i%2 == 0 {
           sa++
       }
       if i%3 == 0 {
           sb++
       }
   }
   if sa > sb {
       fmt.Println(2)
   } else {
       fmt.Println(3)
   }
}
