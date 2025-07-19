package main

import (
   "bufio"
   "fmt"
   "os"
)

const maxn = 201000

// Pair holds two integers
type Pair struct{ first, second int }

var (
   firstArr [maxn]int
   nxt      [maxn]int
   mm       [maxn]Pair
)

func tj(x, y int) {
   nxt[x] = firstArr[y]
   firstArr[y] = x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, s int
   if _, err := fmt.Fscan(reader, &n, &s); err != nil {
       return
   }
   if s&1 == 1 {
       fmt.Fprintln(writer, "No")
       return
   }
   // read values and build lists
   for i := 1; i <= n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       tj(i, x)
   }
   ans := make([]Pair, 0, maxn)
   // process
   for i := s; i >= 1; i-- {
       for firstArr[i] != 0 {
           x := firstArr[i]
           firstArr[i] = nxt[x]
           j := i
           // collect mm and build ans
           for k := 1; k <= i; k++ {
               for j > 0 && firstArr[j] == 0 {
                   j--
               }
               if j == 0 {
                   fmt.Fprintln(writer, "No")
                   return
               }
               mid := firstArr[j]
               // store
               mm[k] = Pair{mid, j}
               // remove mid
               firstArr[j] = nxt[mid]
               ans = append(ans, Pair{x, mid})
           }
           // reinsert with decremented second
           for k := 1; k <= i; k++ {
               p := mm[k]
               tj(p.first, p.second-1)
           }
       }
   }
   // output
   fmt.Fprintln(writer, "Yes")
   fmt.Fprintln(writer, len(ans))
   for _, p := range ans {
       fmt.Fprintln(writer, p.first, p.second)
   }
}
