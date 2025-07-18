#include <bits/stdc++.h>

#define P(x,y) make_pair(x,y)

#define sqr(x)  ( 1ll * (x) * (x))

using namespace std;

const int MX=(1<<20);

int main(){

    printf("302 \n 0 800000\n");

    int s = 60000;

    for(int j=300;j>0;j--){

        printf("%d %d\n", s, j);

        s+=2*j-1;

    }

    printf("%d %d\n", s+(1<<17) , 800000);

    return 0;

}