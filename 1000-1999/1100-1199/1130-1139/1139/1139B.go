package main

import (
   "bufio"
   "fmt"
   "os"
)

// fast integer reader
func readInt(r *bufio.Reader) int {
   num := 0
   sign := 1
   b, err := r.ReadByte()
   for err == nil && (b < '0' || b > '9') && b != '-' {
       b, err = r.ReadByte()
   }
   if err != nil {
       return 0
   }
   if b == '-' {
       sign = -1
       b, err = r.ReadByte()
   }
   for err == nil && b >= '0' && b <= '9' {
       num = num*10 + int(b-'0')
       b, err = r.ReadByte()
   }
   return num * sign
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   n := readInt(reader)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       a[i] = readInt(reader)
   }
   if n == 0 {
       fmt.Print(0)
       return
   }
   h := a[n-1]
   var ans int64 = int64(h)
   for i := n - 2; i >= 0 && h > 0; i-- {
       if h > a[i] {
           h = a[i]
       } else {
           h--
       }
       ans += int64(h)
   }
   fmt.Fprint(os.Stdout, ans)
}
