package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

// readInt reads next integer from stdin
func readInt() (int, error) {
   var x int
   _, err := fmt.Fscan(reader, &x)
   return x, err
}

func main() {
   defer writer.Flush()
   // number of test cases
   t, err := readInt()
   if err != nil {
       return
   }
   for ; t > 0; t-- {
       n, _ := readInt()
       // frequencies of elements
       freq := make([]int, 101)
       for i := 0; i < n; i++ {
           x, _ := readInt()
           if x >= 0 && x < len(freq) {
               freq[x]++
           }
       }
       // mex of A: first i with freq[i] == 0
       mexA := 0
       for mexA < len(freq) && freq[mexA] > 0 {
           mexA++
       }
       // mex of B: first i with freq[i] <= 1
       mexB := 0
       for mexB < len(freq) && freq[mexB] > 1 {
           mexB++
       }
       // output result
       fmt.Fprintln(writer, mexA+mexB)
   }
}
