package main

import (
   "bufio"
   "os"
)

var (
   rdr = bufio.NewReader(os.Stdin)
   wtr = bufio.NewWriter(os.Stdout)
)

// readInt reads next integer from standard input
func readInt() int {
   sign := 1
   var b byte
   // skip non-digit
   for {
       x, err := rdr.ReadByte()
       if err != nil {
           break
       }
       b = x
       if (b >= '0' && b <= '9') || b == '-' {
           break
       }
   }
   if b == '-' {
       sign = -1
       b, _ = rdr.ReadByte()
   }
   val := 0
   for b >= '0' && b <= '9' {
       val = val*10 + int(b-'0')
       x, err := rdr.ReadByte()
       if err != nil {
           break
       }
       b = x
   }
   return val * sign
}

func main() {
   defer wtr.Flush()
   n := readInt()
   a := make([]int, n+1)
   for i := 2; i <= n; i++ {
       a[i] = readInt()
   }
   now := n
   path := make([]int, 0, n)
   for now != 1 {
       path = append(path, now)
       now = a[now]
   }
   // print root
   wtr.WriteString("1")
   // print path in reverse
   for i := len(path) - 1; i >= 0; i-- {
       wtr.WriteByte(' ')
       // convert integer to string
       // use simple conversion
       num := path[i]
       // write digits
       buf := [20]byte{}
       pos := len(buf)
       if num == 0 {
           pos--
           buf[pos] = '0'
       } else {
           for num > 0 {
               pos--
               buf[pos] = byte('0' + num%10)
               num /= 10
           }
       }
       wtr.Write(buf[pos:])
   }
}
