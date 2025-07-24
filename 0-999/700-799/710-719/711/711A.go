package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   rows := make([]string, n)
   for i := 0; i < n; i++ {
       // read each row string
       fmt.Fscan(reader, &rows[i])
   }
   found := false
   for i := 0; i < n && !found; i++ {
       r := []rune(rows[i])
       // check first pair
       if r[0] == 'O' && r[1] == 'O' {
           r[0], r[1] = '+', '+'
           rows[i] = string(r)
           found = true
           break
       }
       // check second pair
       if r[3] == 'O' && r[4] == 'O' {
           r[3], r[4] = '+', '+'
           rows[i] = string(r)
           found = true
           break
       }
   }
   if !found {
       fmt.Println("NO")
       return
   }
   fmt.Println("YES")
   for _, row := range rows {
       fmt.Println(row)
   }
}
