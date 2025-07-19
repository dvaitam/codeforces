package main

import (
   "bufio"
   "fmt"
   "os"
)

func nextInt(r *bufio.Reader) int {
   var c byte
   var err error
   // Skip non-numeric characters
   for {
       c, err = r.ReadByte()
       if err != nil {
           return 0
       }
       if (c >= '0' && c <= '9') || c == '-' {
           break
       }
   }
   sign := 1
   if c == '-' {
       sign = -1
       c, _ = r.ReadByte()
   }
   x := 0
   for ; c >= '0' && c <= '9'; c, _ = r.ReadByte() {
       x = x*10 + int(c-'0')
   }
   return x * sign
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   T := nextInt(reader)
   for t := 0; t < T; t++ {
       n := nextInt(reader)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           a[i] = nextInt(reader)
       }
       maxx, minn := -1, int(1e9+1)
       // Find min and max neighbors of missing values
       for i := 0; i < n; i++ {
           if a[i] == -1 {
               if i > 0 && a[i-1] != -1 {
                   if a[i-1] > maxx {
                       maxx = a[i-1]
                   }
                   if a[i-1] < minn {
                       minn = a[i-1]
                   }
               }
               if i+1 < n && a[i+1] != -1 {
                   if a[i+1] > maxx {
                       maxx = a[i+1]
                   }
                   if a[i+1] < minn {
                       minn = a[i+1]
                   }
               }
           }
       }
       // Determine replacement value
       k := 0
       if maxx != -1 {
           k = (maxx + minn) / 2
       }
       // Replace missing values
       for i := 0; i < n; i++ {
           if a[i] == -1 {
               a[i] = k
           }
       }
       // Compute maximum adjacent difference
       answer := 0
       for i := 0; i+1 < n; i++ {
           d := abs(a[i+1] - a[i])
           if d > answer {
               answer = d
           }
       }
       fmt.Fprintf(writer, "%d %d\n", answer, k)
   }
}
