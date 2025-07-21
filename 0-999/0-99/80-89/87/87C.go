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
   // Grundy numbers up to n
   g := make([]int, n+1)
   // last seen round for mex computation
   const maxG = 512
   last := make([]int, maxG)
   round := 1
   for i := 1; i <= n; i++ {
       // compute possible moves
       // use k piles, k>=2, require k*(k+1)/2 <= i
       offset := 0
       // mark mex
       for k := 2; ; k++ {
           // offset = k*(k-1)/2
           offset = k*(k-1)/2
           // minimal sum needs a1>=k: sum = k*k - k*(k-1)/2 = k*(k+1)/2
           if k*(k+1)/2 > i {
               break
           }
           // check divisibility for a1
           if (i+offset)%k != 0 {
               continue
           }
           a1 := (i + offset) / k
           if a1 < k {
               continue
           }
           // compute xor of piles
           x := 0
           for j := 0; j < k; j++ {
               x ^= g[a1-j]
           }
           if x < maxG {
               last[x] = round
           }
       }
       // mex
       mex := 0
       for ; mex < maxG; mex++ {
           if last[mex] != round {
               break
           }
       }
       g[i] = mex
       round++
   }
   // if losing position
   if g[n] == 0 {
       fmt.Println(-1)
       return
   }
   // find minimal k for first winning move: move to losing (nim xor=0)
   for k := 2; ; k++ {
       if k*(k+1)/2 > n {
           break
       }
       offset := k*(k-1)/2
       if (n+offset)%k != 0 {
           continue
       }
       a1 := (n + offset) / k
       if a1 < k {
           continue
       }
       x := 0
       for j := 0; j < k; j++ {
           x ^= g[a1-j]
       }
       if x == 0 {
           fmt.Println(k)
           return
       }
   }
   // fallback
   fmt.Println(-1)
}
