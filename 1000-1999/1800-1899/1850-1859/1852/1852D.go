package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = int(1e9)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var tt int
   if _, err := fmt.Fscan(reader, &tt); err != nil {
       return
   }
   for tt > 0 {
       tt--
       var n, k int
       fmt.Fscan(reader, &n, &k)
       var foo string
       fmt.Fscan(reader, &foo)
       s := make([]int, n)
       for i := 0; i < n; i++ {
           // convert character to 0 (A) or 1 (B)
           s[i] = int(foo[i]) - int('A')
       }
       initVal := 0
       for i := 0; i+1 < n; i++ {
           if s[i] != s[i+1] {
               initVal++
           }
       }
       // dpMin[i][c][p]: min cost up to i ending with c and parity p
       dpMin := make([][2][2]int, n)
       dpMax := make([][2][2]int, n)
       // initialize
       for i := 0; i < n; i++ {
           for c := 0; c < 2; c++ {
               for p := 0; p < 2; p++ {
                   dpMin[i][c][p] = INF
                   dpMax[i][c][p] = -INF
               }
           }
       }
       for c := 0; c < 2; c++ {
           val := initVal
           if c != s[0] {
               val++
           }
           parity := val & 1
           dpMin[0][c][parity] = val
           dpMax[0][c][parity] = val
       }
       // fill dp
       for i := 1; i < n; i++ {
           for c := 0; c < 2; c++ {
               for t := 0; t < 2; t++ {
                   add := 0
                   if c != t {
                       add++
                   }
                   if c != s[i] {
                       add++
                   }
                   for prevP := 0; prevP < 2; prevP++ {
                       if dpMin[i-1][t][prevP] < INF {
                           newP := (prevP + add) & 1
                           dpMin[i][c][newP] = min(dpMin[i][c][newP], dpMin[i-1][t][prevP] + add)
                       }
                       if dpMax[i-1][t][prevP] > -INF {
                           newP := (prevP + add) & 1
                           dpMax[i][c][newP] = max(dpMax[i][c][newP], dpMax[i-1][t][prevP] + add)
                       }
                   }
               }
           }
       }
       // reconstruct
       res := make([]int, n)
       found := false
       lastParity := k & 1
       for c := 0; c < 2; c++ {
           if dpMin[n-1][c][lastParity] <= k && k <= dpMax[n-1][c][lastParity] {
               res[n-1] = c
               found = true
               break
           }
       }
       if !found {
           writer.WriteString("NO\n")
           continue
       }
       writer.WriteString("YES\n")
       // backtrack
       for i := n - 1; i > 0; i-- {
           c := res[i]
           parity := k & 1
           for t := 0; t < 2; t++ {
               add := 0
               if c != t {
                   add++
               }
               if c != s[i] {
                   add++
               }
               prevParity := (parity + add) & 1
               prevMin := dpMin[i-1][t][prevParity]
               prevMax := dpMax[i-1][t][prevParity]
               if prevMin <= k-add && k-add <= prevMax {
                   res[i-1] = t
                   k -= add
                   break
               }
           }
       }
       // output result
       buf := make([]byte, n)
       for i := 0; i < n; i++ {
           // reconstruct characters
           buf[i] = byte('A') + byte(res[i])
       }
       writer.Write(buf)
       writer.WriteByte('\n')
   }
}
