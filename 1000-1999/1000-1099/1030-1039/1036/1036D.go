package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read ints
   n := int(readInt(reader))
   a := make([]int64, n)
   var sumA int64
   for i := 0; i < n; i++ {
       a[i] = readInt(reader)
       sumA += a[i]
   }
   m := int(readInt(reader))
   b := make([]int64, m)
   var sumB int64
   for i := 0; i < m; i++ {
       b[i] = readInt(reader)
       sumB += b[i]
   }
   if sumA != sumB {
       fmt.Println(-1)
       return
   }
   var i, j int
   var curA, curB int64
   var ans int
   // two pointers
   for i < n || j < m {
       if curA <= curB {
           if i < n {
               curA += a[i]
               i++
           } else {
               // only b remains
               curB += b[j]
               j++
           }
       } else {
           if j < m {
               curB += b[j]
               j++
           } else {
               curA += a[i]
               i++
           }
       }
       if curA == curB {
           ans++
       }
   }
   fmt.Println(ans)
}

// readInt reads next integer (possibly negative) from bufio.Reader
func readInt(r *bufio.Reader) int64 {
   var neg bool
   var c byte
   // skip non-numeric
   for {
       b, err := r.ReadByte()
       if err != nil {
           if err == io.EOF {
               return 0
           }
           panic(err)
       }
       c = b
       if c == '-' || (c >= '0' && c <= '9') {
           break
       }
   }
   if c == '-' {
       neg = true
       // read next byte
       b, err := r.ReadByte()
       if err != nil {
           panic(err)
       }
       c = b
   }
   var x int64
   for c >= '0' && c <= '9' {
       x = x*10 + int64(c-'0')
       b, err := r.ReadByte()
       if err != nil {
           break
       }
       c = b
   }
   if neg {
       x = -x
   }
   return x
}
