package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   // read n
   line, _ := in.ReadString('\n')
   n, _ := strconv.Atoi(line[:len(line)-1])
   // read permutation
   p := make([]int, n+1)
   pos := make([]int, n+1)
   // read rest of line tokens
   for i := 1; i <= n; i++ {
       var x int
       fmt.Fscan(in, &x)
       p[i] = x
       pos[x] = i
   }
   // B[i] = 1 if break between values i and i+1
   // Collect positions of breaks in B[1..n-1]
   var P []int
   // sentinel at 0
   P = append(P, 0)
   for i := 1; i < n; i++ {
       if pos[i+1] - pos[i] != 1 && pos[i] - pos[i+1] != 1 {
           P = append(P, i)
       }
   }
   // sentinel at n
   P = append(P, n)
   var ans int64
   // sum over zero-runs between breaks: subarrays with sum B=0
   for j := 0; j+1 < len(P); j++ {
       r := P[j+1] - P[j] - 1
       if r > 0 {
           ans += int64(r) * int64(r+1) / 2
       }
   }
   // sum over subarrays with exactly one break: sum B = 1
   for j := 1; j+1 < len(P); j++ {
       prevDist := P[j] - P[j-1]
       nextDist := P[j+1] - P[j]
       ans += int64(prevDist) * int64(nextDist)
   }
   fmt.Fprintln(out, ans)
}
