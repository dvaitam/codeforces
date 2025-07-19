package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var cnt [3]int
   var l int
   // read three counts and l
   fmt.Fscan(reader, &cnt[0], &cnt[1], &cnt[2], &l)
   lt := [3]int{0, 0, 0}
   // ans entries: machine, start, end
   var ans [][3]int
   for t := 1; t <= l; t++ {
       mi := -1
       thenTask := -1
       for i := 0; i < 3; i++ {
           ht := t - lt[i]
           can := -1
           // j from min(2, ht-2) down to 0
           for j := min(2, ht-2); j >= 0; j-- {
               if cnt[j] > 0 {
                   can = j
                   break
               }
           }
           if can == -1 {
               continue
           }
           if mi == -1 {
               mi = i
               thenTask = can
           } else {
               if thenTask < can {
                   mi = i
                   thenTask = can
               } else if thenTask > can {
                   // keep existing
               } else if lt[mi] < lt[i] {
                   mi = i
               }
           }
       }
       if mi == -1 {
           continue
       }
       // assign
       cnt[thenTask]--
       lt[mi] = t
       start := t - thenTask - 2
       ans = append(ans, [3]int{mi + 1, start, t})
   }
   // output
   fmt.Fprintln(writer, len(ans))
   for _, v := range ans {
       fmt.Fprintln(writer, v[0], v[1], v[2])
   }
}
