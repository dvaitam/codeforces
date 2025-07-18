#include <stdio.h>


int main(){
    int q;
    int r, l, d;

    scanf("%d", &q);
    for(int i=0; i<q; i++){
        scanf("%d %d %d", &l, &r, &d);
        for(long long j=d; j<2000000010; j+=d){
            if(j < l || j > r){
                printf("%lld\n", j);
                break;
            } if(j >= l && j <= r ){
                long long n=r/d;
                j=n*d;
            }
        }
    }


}