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
   n := len(s)
   // if length is odd, answer is minimal super lucky of next even length
   if n%2 == 1 {
       k := (n + 1) / 2
       // k '4's and k '7's
       for i := 0; i < k; i++ {
           fmt.Print('4')
       }
       for i := 0; i < k; i++ {
           fmt.Print('7')
       }
       return
   }
   L := n
   half := L / 2
   rem4, rem7 := half, half
   type state struct {
       pos    int
       rem4   int
       rem7   int
       alt    byte
   }
   var stack []state
   i := 0
   // greedy equal pass, record alternatives
   for i = 0; i < L; i++ {
       c := s[i]
       // record smallest alt > c
       var alt byte
       if c < '4' {
           if rem4 > 0 {
               alt = '4'
           } else if rem7 > 0 {
               alt = '7'
           }
       } else if c < '7' {
           if rem7 > 0 {
               alt = '7'
           }
       }
       if alt != 0 {
           stack = append(stack, state{i, rem4, rem7, alt})
       }
       // try equal
       if c == '4' && rem4 > 0 {
           rem4--
           continue
       }
       if c == '7' && rem7 > 0 {
           rem7--
           continue
       }
       // equal fails
       break
   }
   // if fully matched and exact rem exhausted
   if i == L && rem4 == 0 && rem7 == 0 {
       fmt.Print(s)
       return
   }
   // backtrack if possible
   for len(stack) > 0 {
       st := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       // build result
       res := make([]byte, L)
       // prefix matched equal
       for j := 0; j < st.pos; j++ {
           res[j] = s[j]
       }
       // place alternative
       res[st.pos] = st.alt
       rem4 = st.rem4
       rem7 = st.rem7
       if st.alt == '4' {
           rem4--
       } else {
           rem7--
       }
       // fill minimal suffix
       idx := st.pos + 1
       for r := 0; r < rem4; r++ {
           res[idx] = '4'
           idx++
       }
       for r := 0; r < rem7; r++ {
           res[idx] = '7'
           idx++
       }
       fmt.Print(string(res))
       return
   }
   // no candidate of same length, output next even length
   k2 := (L/2 + 1)
   for i := 0; i < k2; i++ {
       fmt.Print('4')
   }
   for i := 0; i < k2; i++ {
       fmt.Print('7')
   }
}
