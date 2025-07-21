package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read binary string s
   s, err := reader.ReadString('\n')
   if err != nil && len(s) == 0 {
       return
   }
   // remove newline if present
   if n := len(s); n > 0 && (s[n-1] == '\n' || s[n-1] == '\r') {
       // trim all trailing \r and \n
       end := n
       for end > 0 && (s[end-1] == '\n' || s[end-1] == '\r') {
           end--
       }
       s = s[:end]
   }
   N := len(s)
   // prepare bits reversed, with extra space
   b := make([]byte, N+2)
   for i := 0; i < N; i++ {
       // s[0] is MSB, so reverse to LSB at b[0]
       b[N-1-i] = s[i] - '0'
   }
   var carry int
   var ans int64
   // process bits including one extra for possible carry
   for i := 0; i < N+1; i++ {
       bit := int(b[i]) + carry
       if bit&1 == 0 {
           // even
           if bit == 2 {
               carry = 1
           } else {
               carry = 0
           }
       } else {
           // odd, choose +1 or -1
           if b[i+1] == 1 {
               // subtract: digit = -1
               carry = 1
           } else {
               // add: digit = +1
               carry = 0
           }
           ans++
       }
   }
   // output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
