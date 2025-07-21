package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([][]byte, n)
   var s string
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &s)
       a[i] = []byte(s)
   }
   if n == 0 {
       fmt.Println(0)
       return
   }
   L := len(a[0])
   var sb, sc []byte
   var lenB, lenC int
   // process each string
   for _, str := range a {
       // compute overlap for B
       tb := 0
       maxT := max(len(sb), 0)
       if maxT > L {
           maxT = L
       }
       for t := maxT; t >= 0; t-- {
           // check sb suffix of length t equals str prefix length t
           ok := true
           // sb suffix start index
           start := len(sb) - t
           for i := 0; i < t; i++ {
               if sb[start+i] != str[i] {
                   ok = false
                   break
               }
           }
           if ok {
               tb = t
               break
           }
       }
       // compute overlap for C
       tc := 0
       maxT = max(len(sc), 0)
       if maxT > L {
           maxT = L
       }
       for t := maxT; t >= 0; t-- {
           ok := true
           start := len(sc) - t
           for i := 0; i < t; i++ {
               if sc[start+i] != str[i] {
                   ok = false
                   break
               }
           }
           if ok {
               tc = t
               break
           }
       }
       costB := L - tb
       costC := L - tc
       if costB <= costC {
           // assign to B
           lenB += costB
           // update suffix B
           newAdded := str[tb:]
           sb = append(sb, newAdded...)
           if len(sb) > L {
               sb = sb[len(sb)-L:]
           }
       } else {
           // assign to C
           lenC += costC
           newAdded := str[tc:]
           sc = append(sc, newAdded...)
           if len(sc) > L {
               sc = sc[len(sc)-L:]
           }
       }
   }
   fmt.Println(lenB + lenC)
}
