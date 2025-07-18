#include <iostream>
#include <algorithm>
#include <cstring>
#include <cmath>
#include <cstdlib>
#include <cstdio>
using namespace std;
typedef long long ll;
int get_num(){
    int num = 0;
    char c;
    bool flag = false;
    while((c = getchar()) == ' ' || c == '\r' || c == '\n');
    if(c == '-')
        flag = true;
    else num = c - '0';
    while(isdigit(c = getchar()))
        num = num * 10 + c - '0';
    return (flag ? -1 : 1) * num;
}
const int maxn = 2e5+5;
int a[maxn],opt[maxn];
int cnt = 0;
int n;
int ans = 0;
int calc(int x,int pos){
	if(pos == 1)return x;
    if(opt[pos-1] < x / 2 || x % 2 == 1 || pos == 0)return x;
    return calc(x / 2,pos-1) + x;
}
int main(){
    n = get_num();
    for(int i = 1;i <= n;++i)
        a[i] = get_num();
    sort(a+1,a+n+1);
    int p = a[1];
    opt[++cnt] = 1;
    for(int i = 2;i <= n;++i){
        if(a[i] != p){
            opt[++cnt] = 1;
        }
        else opt[cnt]++;
        p = a[i];
    }
    if(cnt == 1){
        printf("%d\n",n);
        return 0;
    }
    sort(opt+1,opt+cnt+1);
    ans = max(ans,opt[cnt]);
    if(opt[cnt] % 2)opt[cnt]--;
    while(opt[cnt]){
        if(ans > opt[cnt] * 2)break;
        int p = calc(opt[cnt],cnt);
        ans = max(ans,p);
        opt[cnt] -= 2;
    }
    printf("%d\n",ans);
    return 0;
}