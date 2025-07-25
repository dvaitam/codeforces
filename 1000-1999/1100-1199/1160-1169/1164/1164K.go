package main

import "fmt"

// Problem K:
// There are 2018 integers written around a circle, summing to one.
// A good sequence is a sequence of consecutive (clockwise) numbers
// whose sum is positive. We can pair each sequence of length <2018
// with its complement; since sums are integers, exactly one in each
// pair has positive sum, giving 2018*(2018-1)/2 good sequences from
// lengths 1..2017, plus 2018 good full-circle sequences. Thus total =
// 2018*(2018+1)/2 = 2037171.
func main() {
    fmt.Println(2037171)
}
