package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// reverse returns the reverse of the input string
func reverse(s string) string {
   b := []byte(s)
   i, j := 0, len(b)-1
   for i < j {
       b[i], b[j] = b[j], b[i]
       i++
       j--
   }
   return string(b)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k, n int
   if _, err := fmt.Fscan(reader, &k, &n); err != nil {
       return
   }
   // map from string to list of beauties
   m := make(map[string][]int)
   for i := 0; i < k; i++ {
       var s string
       var a int
       fmt.Fscan(reader, &s, &a)
       m[s] = append(m[s], a)
   }
   // sort each list in descending order
   for s, list := range m {
       sort.Slice(list, func(i, j int) bool { return list[i] > list[j] })
       m[s] = list
   }
   crossSum := 0
   palSum := 0
   bestGain := 0
   // process each string group
   for s, A := range m {
       rev := reverse(s)
       if s < rev {
           B, ok := m[rev]
           if !ok {
               continue
           }
           maxPairs := len(A)
           if len(B) < maxPairs {
               maxPairs = len(B)
           }
           for i := 0; i < maxPairs; i++ {
               sum := A[i] + B[i]
               if sum > 0 {
                   crossSum += sum
               } else {
                   break
               }
           }
       } else if s == rev {
           // palindromic strings
           // try pairing and track possible gain from taking single
           i := 0
           for ; i+1 < len(A); i += 2 {
               sum := A[i] + A[i+1]
               if sum > 0 {
                   palSum += sum
                   // gain from dropping this pair and using A[i] as center
                   if gain := A[i] - sum; gain > bestGain {
                       bestGain = gain
                   }
               } else {
                   break
               }
           }
           // leftover element can be used as center
           if i < len(A) && A[i] > bestGain {
               bestGain = A[i]
           }
       }
   }
   result := crossSum + palSum
   if bestGain > 0 {
       result += bestGain
   }
   fmt.Println(result)
}
