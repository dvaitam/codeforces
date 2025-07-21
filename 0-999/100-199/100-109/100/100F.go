package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // coeffs[k] is coefficient for X^k
   coeffs := []int{1}
   for _, ai := range a {
       m := len(coeffs)
       nxt := make([]int, m+1)
       for j := 0; j < m; j++ {
           // ai * coeffs[j] contributes to X^j
           nxt[j] += ai * coeffs[j]
           // coeffs[j] contributes to X^(j+1)
           nxt[j+1] += coeffs[j]
       }
       coeffs = nxt
   }
   // build polynomial string
   var sb strings.Builder
   first := true
   // degrees from highest to 0
   deg := len(coeffs) - 1
   for k := deg; k >= 0; k-- {
       c := coeffs[k]
       if c == 0 {
           continue
       }
       if first {
           // sign for first term
           if c < 0 {
               sb.WriteString("-")
               c = -c
           }
           first = false
       } else {
           if c < 0 {
               sb.WriteString(" - ")
               c = -c
           } else {
               sb.WriteString(" + ")
           }
       }
       // term body
       if k == 0 {
           sb.WriteString(strconv.Itoa(c))
       } else {
           if c != 1 {
               sb.WriteString(strconv.Itoa(c))
               sb.WriteString("*")
           }
           sb.WriteString("X")
           if k != 1 {
               sb.WriteString("^")
               sb.WriteString(strconv.Itoa(k))
           }
       }
   }
   // prefix p(x) = 
   result := "p(x) = " + sb.String()
   fmt.Println(result)
}
