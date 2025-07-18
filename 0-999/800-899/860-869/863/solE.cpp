#include <stdio.h>
#include <algorithm>
using namespace std;
#define N (200010)
int n;
int l[N], r[N], p[N];
bool cmp (int a, int b){
    return l[a] == l[b] ? r[a] < r[b] : l[a] < l[b];
}
int main(){
    scanf("%d", &n);
    for (int i = 1; i <= n; ++i){
        scanf("%d%d", l+i, r+i);
        p[i] = i;
    }
    sort(p+1,p+n+1, cmp);
    for (int i = 2; i <= n; ++i){
        if (l[p[i-1]] == l[p[i]]){
            printf("%d\n", p[i-1]);
            return 0;
        }
        if (r[p[i-1]] >= r[p[i]]){
            printf("%d\n", p[i]);
            return 0;
        }
    }
    for (int i = 2; i < n; ++i){
        if (r[p[i-1]]+1 >= l[p[i+1]]){
            printf("%d\n", p[i]);
            return 0;
        }
    }
    printf("-1\n");
}