package main

import (
   "bufio"
   "bytes"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   data, _ := reader.ReadBytes('\n')
   // trim newline characters
   s := strings.TrimRight(string(data), "\r\n")
   n := len(s)
   tmp := n >> 2
   ans := make([]byte, tmp)
   for i := 0; i < tmp; i++ {
       k1 := i << 1
       k2 := k1 | 1
       k3 := n - ((i + 1) << 1)
       k4 := k3 + 1
       if s[k1] == s[k3] || s[k1] == s[k4] {
           ans[i] = s[k1]
       } else {
           ans[i] = s[k2]
       }
   }
   var buf bytes.Buffer
   buf.Write(ans)
   if n%4 != 0 {
       buf.WriteByte(s[tmp<<1])
   }
   for i := tmp - 1; i >= 0; i-- {
       buf.WriteByte(ans[i])
   }
   os.Stdout.Write(buf.Bytes())
}
