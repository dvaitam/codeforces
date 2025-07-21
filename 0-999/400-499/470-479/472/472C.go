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
   first := make([]string, n)
   last := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &first[i], &last[i])
   }
   perm := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &perm[i])
       perm[i]--
   }
   // Precompute small and large handle for each person
   small := make([]string, n)
   large := make([]string, n)
   for i := 0; i < n; i++ {
       if first[i] < last[i] {
           small[i], large[i] = first[i], last[i]
       } else {
           small[i], large[i] = last[i], first[i]
       }
   }
   prev := ""
   for _, idx := range perm {
       opt1 := small[idx]
       opt2 := large[idx]
       var chosen string
       ok1 := opt1 > prev
       ok2 := opt2 > prev
       if ok1 && ok2 {
           if opt1 < opt2 {
               chosen = opt1
           } else {
               chosen = opt2
           }
       } else if ok1 {
           chosen = opt1
       } else if ok2 {
           chosen = opt2
       } else {
           fmt.Println("NO")
           return
       }
       prev = chosen
   }
   fmt.Println("YES")
}
