package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(a int64) int64 {
   if a < 0 {
       return -a
   }
   return a
}

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return abs(a)
}

func gcd3(a, b, c int64) int64 {
   return gcd(gcd(a, b), c)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   xs := make([]int64, n)
   ys := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &xs[i], &ys[i])
   }
   // Collect unique lines and count per direction
   lineSet := make(map[string]struct{})
   dirCount := make(map[string]int64)
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           x1, y1 := xs[i], ys[i]
           x2, y2 := xs[j], ys[j]
           A := y2 - y1
           B := x1 - x2
           C := -(A*x1 + B*y1)
           g := gcd3(A, B, C)
           if g != 0 {
               A /= g
               B /= g
               C /= g
           }
           // Normalize sign
           if A < 0 || (A == 0 && B < 0) {
               A, B, C = -A, -B, -C
           }
           lineKey := fmt.Sprintf("%d_%d_%d", A, B, C)
           if _, exists := lineSet[lineKey]; !exists {
               lineSet[lineKey] = struct{}{}
               dirKey := fmt.Sprintf("%d_%d", A, B)
               dirCount[dirKey]++
           }
       }
   }
   // Total unique lines
   L := int64(len(lineSet))
   // Total pairs of lines
   totalPairs := L * (L - 1) / 2
   // Subtract pairs of parallel lines (no intersection)
   var parallelPairs int64
   for _, k := range dirCount {
       parallelPairs += k * (k - 1) / 2
   }
   ans := totalPairs - parallelPairs
   fmt.Println(ans)
}
