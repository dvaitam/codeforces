package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)
   // Pattern1: start with 'r', Pattern2: start with 'b'
   count1_b_in_r, count1_r_in_b := 0, 0
   count2_b_in_r, count2_r_in_b := 0, 0
   for i := 0; i < n; i++ {
       c := s[i]
       // Pattern1: even->'r', odd->'b'
       if i%2 == 0 {
           if c == 'b' {
               count1_b_in_r++
           }
       } else {
           if c == 'r' {
               count1_r_in_b++
           }
       }
       // Pattern2: even->'b', odd->'r'
       if i%2 == 0 {
           if c == 'r' {
               count2_r_in_b++
           }
       } else {
           if c == 'b' {
               count2_b_in_r++
           }
       }
   }
   cost1 := max(count1_b_in_r, count1_r_in_b)
   cost2 := max(count2_b_in_r, count2_r_in_b)
   fmt.Println(min(cost1, cost2))
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}
