package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   type pair struct{ v, d int }
   stack := make([]pair, 0, n)
   answer := 0

   for _, val := range a {
       currDeath := 0
       for len(stack) > 0 && stack[len(stack)-1].v <= val {
           if stack[len(stack)-1].d > currDeath {
               currDeath = stack[len(stack)-1].d
           }
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           currDeath = 0
       } else {
           currDeath++
       }
       if currDeath > answer {
           answer = currDeath
       }
       stack = append(stack, pair{v: val, d: currDeath})
   }

   fmt.Fprintln(writer, answer)
}
