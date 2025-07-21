package main

import "fmt"

func main() {
   var s, t string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   fmt.Scan(&t)
   // reverse s
   bs := []byte(s)
   for i, j := 0, len(bs)-1; i < j; i, j = i+1, j-1 {
       bs[i], bs[j] = bs[j], bs[i]
   }
   // check if reversed s equals t
   if string(bs) == t {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
