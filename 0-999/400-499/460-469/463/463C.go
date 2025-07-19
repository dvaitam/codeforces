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
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // fast integer read
   readInt := func() int {
       var x int
       var neg bool
       b, err := reader.ReadByte()
       for err == nil && (b < '0' || b > '9') && b != '-' {
           b, err = reader.ReadByte()
       }
       if err != nil {
           return 0
       }
       if b == '-' {
           neg = true
           b, _ = reader.ReadByte()
       }
       for err == nil && b >= '0' && b <= '9' {
           x = x*10 + int(b-'0')
           b, err = reader.ReadByte()
       }
       if neg {
           x = -x
       }
       return x
   }

   n := readInt()
   // allocate arrays with padding
   a := make([][]int64, n+2)
   d1 := make([][]int64, n+2)
   d2 := make([][]int64, n+2)
   for i := 0; i <= n+1; i++ {
       a[i] = make([]int64, n+2)
       d1[i] = make([]int64, n+3) // extra for j+1
       d2[i] = make([]int64, n+2)
   }
   for i := 1; i <= n; i++ {
       for j := 1; j <= n; j++ {
           a[i][j] = int64(readInt())
       }
   }
   // build diagonal prefix sums
   for i := 1; i <= n; i++ {
       for j := 1; j <= n; j++ {
           d2[i][j] = d2[i-1][j-1] + a[i][j]
           d1[i][j] = d1[i-1][j+1] + a[i][j]
       }
   }
   var bsum, wsum int64
   var bx, by, wx, wy int
   bx, by = 1, 1
   wx, wy = 2, 1
   for i := 1; i <= n; i++ {
       for j := 1; j <= n; j++ {
           // sum along both diagonals through (i,j)
           k := min(n-i, n-j)
           tmp := d2[i+k][j+k]
           k = min(n-i, j-1)
           tmp += d1[i+k][j-k]
           tmp -= a[i][j]
           if (i+j)%2 == 0 {
               if tmp > bsum {
                   bsum = tmp
                   bx, by = i, j
               }
           } else {
               if tmp > wsum {
                   wsum = tmp
                   wx, wy = i, j
               }
           }
       }
   }
   fmt.Fprintln(writer, bsum+wsum)
   fmt.Fprintf(writer, "%d %d %d %d", bx, by, wx, wy)
}
