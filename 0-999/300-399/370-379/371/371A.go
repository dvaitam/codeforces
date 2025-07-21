package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   changes := 0
   // For each position modulo k, count values
   for i := 0; i < k; i++ {
       cnt1, cnt2 := 0, 0
       // traverse group
       for j := i; j < n; j += k {
           if a[j] == 1 {
               cnt1++
           } else {
               cnt2++
           }
       }
       // size of this group
       groupSize := cnt1 + cnt2
       // need to change all except the majority
       if cnt1 > cnt2 {
           changes += groupSize - cnt1
       } else {
           changes += groupSize - cnt2
       }
   }
   fmt.Println(changes)
}
