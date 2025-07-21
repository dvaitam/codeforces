package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var m, n int64
   if _, err := fmt.Fscan(in, &m, &n); err != nil {
       return
   }
   // Impossible cases
   if n < 2 {
       fmt.Println(-1)
       return
   }
   // Small boat capacities only allow trivial m=1
   if n <= 3 && m > 1 {
       fmt.Println(-1)
       return
   }
   // Single pair trivial
   if m == 1 {
       fmt.Println(1)
       return
   }
   // Can carry all at once
   if n >= 2*m {
       fmt.Println(1)
       return
   }
   // State counts
   gl, gr := m, int64(0)
   wl, wr := m, int64(0)
   crossings := int64(0)
   // Initial send one wolf to balance
   wl--
   wr++
   crossings++
   // Shuttle goats in batches of up to n-1, returning one wolf each time
   for gl > 0 {
       // send goats
       sendG := min(gl, n-1)
       gl -= sendG
       gr += sendG
       crossings++
       if gl == 0 {
           break
       }
       // return one wolf
       wl++
       wr--
       crossings++
   }
   // Shuttle wolves in batches of up to n, returning one goat as needed
   for wl > 0 {
       sendW := min(wl, n)
       wl -= sendW
       wr += sendW
       crossings++
       if wl == 0 {
           break
       }
       // return one goat
       gr--
       gl++
       crossings++
   }
   fmt.Println(crossings)
}
