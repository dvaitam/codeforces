#include <iostream>
#include <algorithm>

using namespace std;

int main(int argc, char **argv)
{
    ios_base::sync_with_stdio(false);
    cin.tie(0), cout.tie(0), cout.precision(15);
    int a, b, c; //1 100
    cin >> a >> b >> c;
    if(a > b){
        swap(a, b);
    }
    //a <= b
    if(b > c){
        swap(b, c);
    }
    //b <= c
    if(a > b){
        swap(a, b);
    }
    if(c - b < a){ //a + b > c. a + c(max) > b b + c(max) > a
        cout << 0 << endl;
    }else{
        cout << c - b - a + 1 << endl;
    }
    return 0;
}