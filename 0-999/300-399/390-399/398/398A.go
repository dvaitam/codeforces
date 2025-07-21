package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var a, b int
   if _, err := fmt.Fscan(reader, &a, &b); err != nil {
       return
   }
   // Special cases
   if a == 0 {
       // only x's, one block
       score := -int64(b) * int64(b)
       fmt.Fprintln(writer, score)
       fmt.Fprintln(writer, strings.Repeat("x", b))
       return
   }
   if b == 0 {
       // only o's, one block
       score := int64(a) * int64(a)
       fmt.Fprintln(writer, score)
       fmt.Fprintln(writer, strings.Repeat("o", a))
       return
   }
   // Try splitting o's into k blocks (1 <= k <= a)
   var bestS int64 = math.MinInt64
   bestK := 1
   for k := 1; k <= a; k++ {
       // sum over o-blocks: one of size (a-k+1), k-1 of size 1
       big := int64(a - k + 1)
       sumO := big*big + int64(k-1)
       // total x-block positions = k+1, but cannot exceed b
       t := k + 1
       m := t
       if b < t {
           m = b
       }
       // distribute b x's into m blocks to minimize sum of squares
       base := b / m
       r := b % m
       var sumX int64
       sumX = int64(r) * int64(base+1) * int64(base+1)
       sumX += int64(m-r) * int64(base) * int64(base)
       // score
       S := sumO - sumX
       if S > bestS {
           bestS = S
           bestK = k
       }
   }
   // Reconstruct with bestK
   k := bestK
   // Prepare o-block lengths
   oLens := make([]int, k)
   oLens[0] = a - k + 1
   for i := 1; i < k; i++ {
       oLens[i] = 1
   }
   // Prepare x-slot lengths (k+1 slots)
   slots := make([]int, k+1)
   t := k + 1
   m := t
   if b < t {
       m = b
   }
   if m > 0 {
       base := b / m
       r := b % m
       for i := 0; i < m; i++ {
           if i < r {
               slots[i] = base + 1
           } else {
               slots[i] = base
           }
       }
   }
   // Output result
   fmt.Fprintln(writer, bestS)
   var sb strings.Builder
   for i := 0; i <= k; i++ {
       if slots[i] > 0 {
           sb.WriteString(strings.Repeat("x", slots[i]))
       }
       if i < k {
           sb.WriteString(strings.Repeat("o", oLens[i]))
       }
   }
   fmt.Fprintln(writer, sb.String())
}
