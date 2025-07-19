package main

import (
   "bufio"
   "os"
   "sort"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   n64 := readInt(reader)
   _ = readInt(reader) // m, unused
   n := int(n64)
   a := make([]int64, n)
   var sum int64
   for i := 0; i < n; i++ {
       a[i] = readInt(reader)
       sum += a[i]
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   var maxx int64
   // process from largest to second
   for i := int(n) - 1; i > 0; i-- {
       if a[i-1] >= a[i]-1 {
           if a[i] <= 1 {
               a[i-1] = 1
           } else {
               a[i-1] = a[i] - 1
           }
           maxx++
       } else {
           maxx += a[i] - a[i-1]
       }
   }
   maxx += a[0]
   // result = sum - maxx
   writer.WriteString(strconv.FormatInt(sum-maxx, 10))
   writer.WriteByte('\n')
}

// readInt reads an integer from bufio.Reader
func readInt(r *bufio.Reader) int64 {
   var x int64
   var neg bool
   for {
       b, err := r.ReadByte()
       if err != nil {
           return x
       }
       if b == '-' {
           neg = true
           break
       }
       if b >= '0' && b <= '9' {
           x = int64(b - '0')
           break
       }
   }
   for {
       b, err := r.ReadByte()
       if err != nil {
           break
       }
       if b < '0' || b > '9' {
           break
       }
       x = x*10 + int64(b-'0')
   }
   if neg {
       return -x
   }
   return x
}
