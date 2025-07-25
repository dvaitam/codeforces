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
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   N := 2 * n
   // prev0[b], prev1[b] correspond to dp at pos+1
   prev0 := make([]int, n+2)
   prev1 := make([]int, n+2)
   // initialize at pos = N: only b=0 valid, dp values are zero (slices are zero-initialized)
   // iterate pos from N-1 down to 0
   for pos := N - 1; pos >= 0; pos-- {
       maxb := min(pos, N-pos)
       cur0 := make([]int, n+2)
       cur1 := make([]int, n+2)
       for b := 0; b <= maxb; b++ {
           sumAll := 0
           bestVal := 0
           // left child: add '('
           // child state: (pos+1, b+1) valid if b+1 <= N-(pos+1)
           if b+1 <= N-pos-1 {
               p0 := prev0[b+1]
               p1 := prev1[b+1]
               mx := p0
               if p1 > mx {
                   mx = p1
               }
               sumAll += mx
               // if matching this child edge
               val := 1 + p1 - mx
               if val > bestVal {
                   bestVal = val
               }
           }
           // right child: add ')'
           if b > 0 {
               p0 := prev0[b-1]
               p1 := prev1[b-1]
               mx := p0
               if p1 > mx {
                   mx = p1
               }
               sumAll += mx
               val := 1 + p1 - mx
               if val > bestVal {
                   bestVal = val
               }
           }
           cur1[b] = sumAll
           // in cur0, can choose to match one child or none
           if bestVal > 0 {
               cur0[b] = sumAll + bestVal
           } else {
               cur0[b] = sumAll
           }
       }
       // swap prev and cur
       prev0 = cur0
       prev1 = cur1
   }
   // result at pos=0, b=0
   res := prev0[0]
   const MOD = 1000000007
   if res >= MOD {
       res %= MOD
   }
   fmt.Println(res)
}
