#include<bits/stdc++.h>

#define ll long long
#define INF  2147483647

inline int inp(){
    char c = getchar();
    while(c < '0' || c > '9')
        c = getchar();
    int sum = 0;
    while(c >= '0' && c <= '9'){
        sum = sum * 10 + c - '0';
        c = getchar();
    }
    return sum;
}

ll f[10];
char s[100010];

inline ll min(ll a, ll b){
    return a < b ? a : b;
}

int main(){
    int n = inp();
    scanf("%s", s + 1);
    for(int i = 1; i <= 7; i++)
        f[i] = 800000000000000000;
    f[0] = 0;
    for(int t = 1; t <= n; t++){
        ll a = inp();
        for(int i = 3; i >= 0; i--){
            if(i == 0 && s[t] == 'h'){
                // printf("%d %d\n", t, i);
                f[i + 1] = min(f[i + 1], f[i]);
                f[i] += a;
            } else if(i == 1 && s[t] == 'a'){
                // printf("%d %d\n", t, i);
                f[i + 1] = min(f[i + 1], f[i]);
                f[i] += a;
            } else if(i == 2 && s[t] == 'r'){
                // printf("%d %d\n", t, i);
                f[i + 1] = min(f[i + 1], f[i]);
                f[i] += a;
            } else if(i == 3 && s[t] == 'd'){
                // printf("%d %d\n", t, i);
                f[i + 1] = min(f[i + 1], f[i]);
                f[i] += a;
            }
        }
    }
    printf("%I64d\n", min(min(f[0], f[1]), min(f[2], f[3])));
}