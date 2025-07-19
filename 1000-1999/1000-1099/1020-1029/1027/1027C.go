package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   te := readInt(reader)
   for ; te > 0; te-- {
       n := readInt(reader)
       // Read array
       a := make([]int, n)
       maxVal := 0
       for i := 0; i < n; i++ {
           v := readInt(reader)
           a[i] = v
           if v > maxVal {
               maxVal = v
           }
       }
       // Count frequencies
       cnt := make([]int, maxVal+1)
       for _, v := range a {
           cnt[v]++
       }
       // Collect unique values in ascending order
       keys := make([]int, 0, len(cnt))
       for v, f := range cnt {
           if f > 0 {
               keys = append(keys, v)
           }
       }
       // Since iteration over map or slice index order: keys from 0..maxVal, ascending

       var lst, A, B int
       minRatio := math.MaxFloat64
       for _, v := range keys {
           f := cnt[v]
           if f > 3 {
               minRatio = 1.0
               A, B = v, v
           }
           if f > 1 {
               if lst != 0 {
                   ratio := float64(v) / float64(lst)
                   if ratio < minRatio {
                       minRatio = ratio
                       A, B = lst, v
                   }
               }
               lst = v
           }
       }
       fmt.Fprintf(writer, "%d %d %d %d\n", A, A, B, B)
   }
}

// readInt reads next integer from bufio.Reader
func readInt(reader *bufio.Reader) int {
   var x int
   var neg bool
   for {
       b, err := reader.ReadByte()
       if err != nil {
           return 0
       }
       if b == '-' {
           neg = true
           break
       }
       if b >= '0' && b <= '9' {
           x = int(b - '0')
           break
       }
   }
   for {
       b, err := reader.ReadByte()
       if err != nil {
           break
       }
       if b < '0' || b > '9' {
           break
       }
       x = x*10 + int(b-'0')
   }
   if neg {
       x = -x
   }
   return x
}
