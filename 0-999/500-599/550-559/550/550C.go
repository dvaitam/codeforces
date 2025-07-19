package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   // Iterate over multiples of 8 from 0 to 992
   for i := 0; i < 1000; i += 8 {
       cand := strconv.Itoa(i)
       k := 0
       // Check if cand is a subsequence of s
       for j := 0; j < len(s) && k < len(cand); j++ {
           if s[j] == cand[k] {
               k++
           }
       }
       if k == len(cand) {
           fmt.Println("YES")
           fmt.Println(cand)
           return
       }
   }
   fmt.Println("NO")
}
