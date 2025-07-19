package main

import (
   "fmt"
)

func main() {
   var s string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   n := len(s)
   // positions of each character
   c := make([][]int, 26)
   for i := 0; i < n; i++ {
       ch := s[i] - 'a'
       if ch >= 0 && ch < 26 {
           c[ch] = append(c[ch], i)
       }
   }
   total := 0
   // temporary count array
   var book [26]int
   for ch := 0; ch < 26; ch++ {
       if len(c[ch]) == 0 {
           continue
       }
       maxRes := 0
       // try all non-zero shifts
       for l := 1; l < n; l++ {
           // reset counts
           for i := 0; i < 26; i++ {
               book[i] = 0
           }
           // count occurrences after shift
           for _, pos := range c[ch] {
               np := pos + l
               if np >= n {
                   np -= n
               }
               book[s[np]-'a']++
           }
           // count letters appearing exactly once
           res := 0
           for i := 0; i < 26; i++ {
               if book[i] == 1 {
                   res++
               }
           }
           if res > maxRes {
               maxRes = res
           }
       }
       total += maxRes
   }
   // output answer with 12 decimal places
   fmt.Printf("%.12f\n", float64(total)/float64(n))
}
