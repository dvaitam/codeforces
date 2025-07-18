#include <cstdio>
#include <iostream>
#include <cstring>
#include <algorithm>
using namespace std;
int s[30];
int a[30];
int ss[30];
int sum, lim;
int main(){
    scanf("%d%d", &sum, &lim);
    memset(s, 0, sizeof(s));
    memset(a, 0, sizeof(a));
    for(int i = 1; i <= lim; i++){
        for(int j = 0; (1 << j) <= sum; j++){
            if((1 << j) & i){
                s[j]++;
                break;
            }

        }
    }
    int num[30];
    memset(num, 0, sizeof(num));
    int f = 0;
    for(int i = 0; (1 << i) <= sum; i++){
        if(sum & (1 << i)){
            int tot = 1;
            for(int j = i; j >= 0; j--){
                if(s[j] >= tot){
                    num[j] += tot;
                    s[j] -= tot;
                    break;
                }
                else {
                    tot -= s[j];
                    num[j] += s[j];
                    s[j] = 0;
                    tot *= 2;
                    if(j == 0){
                        f = 1;
                        break;
                    }
                }
            }
            if(f == 1)break;
        }
    }
    if(f == 1)printf("-1\n");
    else {
        int tot = 0;
        for(int i = 0; (1 << i) <= sum; i++){
            tot += num[i];
        }
        printf("%d\n", tot);
        for(int i = 1; i <= lim; i++){
            for(int j = 0; (1 << j) <= sum; j++){
                if(i & (1 << j)){
                    if(num[j] > 0){
                        num[j]--;
                        printf("%d ", i);
                    }
                    break;
                }
            }
        }
        printf("\n");
    }
    return 0;
}