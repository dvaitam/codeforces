package main

import (
   "fmt"
   "strings"
)

func main() {
   var a, b, k int
   if _, err := fmt.Scan(&a, &b, &k); err != nil {
       return
   }
   // Special cases: no zeros or only one one
   if a == 0 || b == 1 {
       if k > 0 {
           fmt.Println("NO")
       } else {
           fmt.Println("YES")
           s := strings.Repeat("1", b) + strings.Repeat("0", a)
           fmt.Println(s)
           fmt.Println(s)
       }
       return
   }
   // Maximum possible distance is a+b-2
   if k > a+b-2 {
       fmt.Println("NO")
       return
   }
   // Build two binary strings x and y of length a+b
   n := a + b
   x := make([]byte, n)
   y := make([]byte, n)
   for i := 0; i < n; i++ {
       x[i] = '1'
       y[i] = '1'
   }
   // Place a zero at start of x
   x[0] = '0'
   a--
   // Place zeros in both strings for first k-1 positions
   for i := 1; i < k && a > 0; i++ {
       x[i] = '0'
       y[i] = '0'
       a--
   }
   // Place a zero only in y at position k
   y[k] = '0'
   // Place remaining zeros in both strings
   idx := k
   for a > 0 {
       idx++
       x[idx] = '0'
       y[idx] = '0'
       a--
   }
   // Reverse to finalize
   reverse(x)
   reverse(y)
   fmt.Println("YES")
   fmt.Println(string(x))
   fmt.Println(string(y))
}

// reverse reverses a byte slice in place
func reverse(s []byte) {
   for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
       s[i], s[j] = s[j], s[i]
   }
}
