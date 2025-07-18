#include <iostream>
#include <stdlib.h>
#include <algorithm>
#include <stdio.h>
#define N_MAX 150
using namespace std;

int main(void){
        double pi = 3.141592653589793, S = 0.0;
        int N;
        int A[N_MAX];

    cin >> N;
    for(int i = 0; i < N; i ++)
        cin >> A[i];

    sort(A, A+N);

    if(N%2){
        for(int i = 0; i < N; i ++)
            if(!(i%2))
                S += (pi * (A[i] * A[i]));
            else
                S -= (pi * (A[i] * A[i]));
    }
    else{
        for(int i = 0; i < N; i ++)
            if(!(i%2))
                S -= (pi * (A[i] * A[i]));
            else
                S += (pi * (A[i] * A[i]));
    }

    printf("%.15f\n", S);

    return 0;
}