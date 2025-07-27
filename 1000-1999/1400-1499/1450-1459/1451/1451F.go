package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   const MAXN, MAXM = 100, 100
   maxSize := MAXN*MAXM + 5
   // Precompute Grundy values for positions (i,j)
   g := make([][]int, MAXN+2)
   for i := range g {
       g[i] = make([]int, MAXM+2)
   }
   seen := make([]int, maxSize)
   timestamp := 1
   for i := MAXN; i >= 1; i-- {
       for j := MAXM; j >= 1; j-- {
           timestamp++
           // include move to terminal (kill), Grundy 0
           seen[0] = timestamp
           // consider moves to any (x,y) >= (i,j), excluding (i,j)
           for x := i; x <= MAXN; x++ {
               for y := j; y <= MAXM; y++ {
                   if x == i && y == j {
                       continue
                   }
                   v := g[x][y]
                   if v < maxSize {
                       seen[v] = timestamp
                   }
               }
           }
           // mex
           mex := 0
           for seen[mex] == timestamp {
               mex++
           }
           g[i][j] = mex
       }
   }
   // Process test cases
   for ; t > 0; t-- {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       xor := 0
       for i := 1; i <= n; i++ {
           for j := 1; j <= m; j++ {
               var a int
               fmt.Fscan(reader, &a)
               if a&1 == 1 {
                   xor ^= g[i][j]
               }
           }
       }
       if xor != 0 {
           fmt.Fprintln(writer, "Ashish")
       } else {
           fmt.Fprintln(writer, "Jeel")
       }
   }
