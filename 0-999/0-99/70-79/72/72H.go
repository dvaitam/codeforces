package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   sign := ""
   if len(s) > 0 && s[0] == '-' {
       sign = "-"
       s = s[1:]
   }
   // remove leading zeros before reverse
   i := 0
   for i < len(s) && s[i] == '0' {
       i++
   }
   s = s[i:]
   if len(s) == 0 {
       fmt.Println("0")
       return
   }
   // reverse digits
   bs := []byte(s)
   for i, j := 0, len(bs)-1; i < j; i, j = i+1, j-1 {
       bs[i], bs[j] = bs[j], bs[i]
   }
   // remove leading zeros after reverse
   j := 0
   for j < len(bs) && bs[j] == '0' {
       j++
   }
   bs = bs[j:]
   if len(bs) == 0 {
       fmt.Println("0")
       return
   }
   // output result
   if sign != "" {
       fmt.Print(sign)
   }
   fmt.Println(string(bs))
}
