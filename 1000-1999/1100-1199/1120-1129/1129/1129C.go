package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   N   = 4600005
   mod = 1000000007
)

var (
   goArr  [N][2]int32
   dp     [N]int32
   have   [5][]bool
   id     int
   curAns int32
   s      []int
   node   [3001]int
   kArr   [3001]int
)

func sum(a, b int32) int32 {
   c := a + b
   if c >= mod {
       c -= mod
   }
   return c
}

// add processes the current bit string s into the trie and updates dp and curAns
func add() {
   v := 0
   it := 0
   for _, c := range s {
       kArr[it] = c
       node[it] = v
       if goArr[v][c] == 0 {
           // new node
           j := it
           val := 0
           var dpi int32
           for t := 1; t <= 4 && j >= 0; t++ {
               val = (val<<1) + kArr[j]
               if have[t][val] {
                   dpi = sum(dpi, dp[node[j]])
               }
               j--
           }
           goArr[v][c] = int32(id)
           dp[id] = dpi
           curAns = sum(curAns, dpi)
           id++
       }
       v = int(goArr[v][c])
       it++
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   dp[0] = 1
   id = 1
   // initialize allowed patterns
   for t := 1; t <= 4; t++ {
       size := 1 << t
       have[t] = make([]bool, size)
       for i := 0; i < size; i++ {
           // build LSB-first pattern
           str := make([]byte, t)
           for k := 0; k < t; k++ {
               if (i>>k)&1 == 1 {
                   str[k] = '1'
               } else {
                   str[k] = '0'
               }
           }
           // exclude forbidden length-4 patterns
           if t == 4 {
               s4 := string(str)
               if s4 == "1100" || s4 == "1010" || s4 == "0111" || s4 == "1111" {
                   continue
               }
           }
           // compute MSB-first value
           val := 0
           for k := t - 1; k >= 0; k-- {
               val = (val << 1) + int(str[k]-'0')
           }
           have[t][val] = true
       }
   }

   var n int
   fmt.Fscan(reader, &n)
   for i := 0; i < n; i++ {
       // read next bit
       b, err := reader.ReadByte()
       for err == nil && b != '0' && b != '1' {
           b, err = reader.ReadByte()
       }
       bit := int(b - '0')
       // prepend bit to s
       s = append([]int{bit}, s...)
       add()
       // output current answer
       fmt.Fprintln(writer, curAns)
   }
}
