package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var x, k int
   fmt.Fscan(reader, &x, &k)
   // used[i] indicates identifier i is already occupied by known rounds
   used := make([]bool, x)
   for i := 0; i < k; i++ {
       var t int
       fmt.Fscan(reader, &t)
       if t == 1 {
           var num2, num1 int
           fmt.Fscan(reader, &num2, &num1)
           if num2 < x {
               used[num2] = true
           }
           if num1 < x {
               used[num1] = true
           }
       } else {
           var num int
           fmt.Fscan(reader, &num)
           if num < x {
               used[num] = true
           }
       }
   }
   var sumMin, sumMax int
   for i := 1; i < x; {
       if used[i] {
           i++
           continue
       }
       j := i
       for j < x && !used[j] {
           j++
       }
       L := j - i
       // In a gap of length L, maximum Div2 rounds = L (all solo)
       // minimum Div2 rounds = ceil(L/2) (max simultaneous)
       sumMax += L
       sumMin += (L + 1) / 2
       i = j
   }
   fmt.Printf("%d %d\n", sumMin, sumMax)
}
