package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   // Read the number as string to process its digits
   fmt.Fscan(reader, &s)
   n := len(s)
   digits := make([]int, n)
   maxd := 0
   for i, c := range s {
       d := int(c - '0')
       digits[i] = d
       if d > maxd {
           maxd = d
       }
   }
   // The minimum count is max digit value
   results := make([]int, 0, maxd)
   // Greedily build each quasibinary number
   for k := 0; k < maxd; k++ {
       num := 0
       for i := 0; i < n; i++ {
           num *= 10
           if digits[i] > 0 {
               num++
               digits[i]--
           }
       }
       results = append(results, num)
   }
   // Output
   fmt.Println(maxd)
   for i, v := range results {
       if i > 0 {
           fmt.Print(" ")
       }
       fmt.Print(v)
   }
   fmt.Println()
}
