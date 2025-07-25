package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, q int
   fmt.Fscan(reader, &n, &m, &q)
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       row := make([]int, m)
       for j := 0; j < m; j++ {
           row[j] = int(s[j] - '0')
       }
       a[i] = row
   }
   // prefix sums of original matrix
   P := make([][]int64, n+1)
   for i := range P {
       P[i] = make([]int64, m+1)
   }
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           P[i+1][j+1] = P[i+1][j] + P[i][j+1] - P[i][j] + int64(a[i][j])
       }
   }
   total := int64(n * m)
   // helper: count ones in Thue-Morse sequence up to length k (0..k-1)
   var tmOnes func(int64) int64
   tmOnes = func(k int64) int64 {
       if k <= 0 {
           return 0
       }
       // largest power of two <= k
       b := 63 - bits.LeadingZeros64(uint64(k))
       p := int64(1) << b
       if p == k {
           return p / 2
       }
       r := k - p
       return p/2 + (r - tmOnes(r))
   }
   // sum of F over [1..x]x[1..y]
   var sumRect func(int64, int64) int64
   sumRect = func(x, y int64) int64 {
       if x <= 0 || y <= 0 {
           return 0
       }
       U := x / int64(n)
       R := x % int64(n)
       V := y / int64(m)
       C := y % int64(m)
       U1 := tmOnes(U)
       U0 := U - U1
       V1 := tmOnes(V)
       V0 := V - V1
       // full blocks
       countEven := U0*V0 + U1*V1
       countOdd := U0*V1 + U1*V0
       res := countEven*P[n][m] + countOdd*(total-P[n][m])
       // partial columns (full rows)
       if C > 0 {
           sumCols := P[n][C]
           blockSize := int64(n) * C
           pv := bits.OnesCount64(uint64(V)) & 1
           var uSame int64
           if pv == 0 {
               uSame = U0
           } else {
               uSame = U1
           }
           uDiff := U - uSame
           res += uSame*sumCols + uDiff*(blockSize-sumCols)
       }
       // partial rows (full columns)
       if R > 0 {
           sumRows := P[R][m]
           blockSize := R * int64(m)
           pu := bits.OnesCount64(uint64(U)) & 1
           var vSame int64
           if pu == 0 {
               vSame = V0
           } else {
               vSame = V1
           }
           vDiff := V - vSame
           res += vSame*sumRows + vDiff*(blockSize-sumRows)
       }
       // partial corner
       if R > 0 && C > 0 {
           sumCorner := P[R][C]
           blockSize := R * C
           pu := bits.OnesCount64(uint64(U)) & 1
           pv := bits.OnesCount64(uint64(V)) & 1
           if pu == pv {
               res += sumCorner
           } else {
               res += blockSize - sumCorner
           }
       }
       return res
   }
   // answer queries
   for i := 0; i < q; i++ {
       var x1, y1, x2, y2 int64
       fmt.Fscan(reader, &x1, &y1, &x2, &y2)
       ans := sumRect(x2, y2) - sumRect(x1-1, y2) - sumRect(x2, y1-1) + sumRect(x1-1, y1-1)
       fmt.Fprintln(writer, ans)
   }
