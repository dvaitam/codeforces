package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   const INF = 1 << 30
   // avail_time[v]: earliest operation after which card v is in hand
   avail := make([]int, n+1)
   for i := 1; i <= n; i++ {
       avail[i] = INF
   }
   for i := 0; i < n; i++ {
       if a[i] != 0 {
           avail[a[i]] = 0
       }
   }
   for i := 0; i < n; i++ {
       if b[i] != 0 {
           // b[i] drawn after i+1 operations
           if avail[b[i]] > i+1 {
               avail[b[i]] = i + 1
           }
       }
   }
   // check special tail case
   tail := b[n-1]
   if tail != 0 {
       ok := true
       // check suffix b[n-tail .. n-1] is 1..tail
       if tail <= n {
           for i := 0; i < tail; i++ {
               if b[n-tail+i] != i+1 {
                   ok = false
                   break
               }
           }
       } else {
           ok = false
       }
       if ok {
           // verify all cards > tail available in time
           valid := true
           for v := tail + 1; v <= n; v++ {
               // need avail[v] <= v - tail - 1
               if avail[v] > v-tail-1 {
                   valid = false
                   break
               }
           }
           if valid {
               fmt.Println(n - tail)
               return
           }
       }
   }
   // general case
   maxDiff := 0
   for v := 1; v <= n; v++ {
       diff := avail[v] - v + 1
       if diff > maxDiff {
           maxDiff = diff
       }
   }
   if maxDiff < 0 {
       maxDiff = 0
   }
   fmt.Println(maxDiff + n)
}
