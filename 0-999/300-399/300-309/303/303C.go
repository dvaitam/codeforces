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
   maxA := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > maxA {
           maxA = a[i]
       }
   }
   // stamps for marking residues, timestamp technique
   maxM := maxA + 2
   stamps := make([]int, maxM)
   stampID := 1
   // minimal m must be >= max(1, n-k)
   start := n - k
   if start < 1 {
       start = 1
   }
   // search m
   for m := start; ; m++ {
       if m >= maxM {
           // for m > maxA, all residues are unique
           fmt.Println(m)
           return
       }
       coll := 0
       for _, v := range a {
           r := v % m
           if stamps[r] != stampID {
               stamps[r] = stampID
           } else {
               coll++
               if coll > k {
                   break
               }
           }
       }
       stampID++
       if coll <= k {
           fmt.Println(m)
           return
       }
   }
}
