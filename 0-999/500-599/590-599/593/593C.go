package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

// ati returns string "abs((t(sign)i))"
func ati(i int) string {
   s := "abs((t"
   if i < 0 {
       i = -i
       s += "+"
   } else {
       s += "-"
   }
   s += strconv.Itoa(i)
   s += "))"
   return s
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   f := "0"
   g := "0"
   for i := 0; i < n; i++ {
       var x, y, r int
       fmt.Fscan(reader, &x, &y, &r)
       if x%2 != 0 {
           x++
       }
       if y%2 != 0 {
           y++
       }
       halfX := strconv.Itoa(x / 2)
       halfY := strconv.Itoa(y / 2)
       // Append term to f
       f = "(" + f + "+(" + halfX + "*((" + ati(i-1) + "+" + ati(i+1) + ")-(" + ati(i) + "+" + ati(i) + "))))"
       // Append term to g
       g = "(" + g + "+(" + halfY + "*((" + ati(i-1) + "+" + ati(i+1) + ")-(" + ati(i) + "+" + ati(i) + "))))"
   }
   fmt.Println(f)
   fmt.Println(g)
}
