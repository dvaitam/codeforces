#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
const int maxn = 505;
const int maxm = 10000;
const int inf = 0x3f3f3f3f;
const int mod = 1e9 + 7;

int n;
int a[5005],b[5005],v[5005];
int a0,a1,a2,sumb;
int main()
{
    scanf("%d",&n);
    for(int i = 1;i<=n;i++) {
        scanf("%1d",&a[i]);
    }
    for(int i = 1;i<=n;i++) {
        scanf("%1d",&b[i]);
        if(b[i]) sumb++;
        v[i] = a[i] + b[i];
        if(v[i]==0) a0++;
        else if(v[i]==1) a1++;
        else if(v[i]==2) a2++;
    }
    int ans = 0;
    for(int i = 0;i<=a2;i++){
        for(int j = 0;j<=a1;j++){
            if(i+j>n/2) continue;
            //if(i==0&&j==0) continue;
            if(2*i+j==sumb){
                int cnt1 = j,cnt2 = i,cnt0 = n/2-i-j;
                if(cnt0>a0) continue;
                for(int k = 1;k<=n;k++){
                    if(v[k]==1&&cnt1){
                        printf("%d ",k); cnt1--;
                    }
                    else if(v[k]==2&&cnt2){
                        printf("%d ",k); cnt2--;
                    }
                    else if(v[k]==0&&cnt0){
                        printf("%d ",k); cnt0--;
                    }
                }
                return 0;
            }
        }
    }
    printf("-1");
}
/*
2
11
00

2
00
00
*/