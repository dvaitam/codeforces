#include<stdio.h>
#include<bits/stdc++.h>
#define maxn 200010
using namespace std;
int a[maxn];
int main(){
    long long int n;
    int k;
    scanf("%lld%d",&n,&k);
    a[0]=0;
    for(long long int i=1ll;i<=n;i++){
        scanf("%d",&a[i]);
        a[i]=a[i]^a[i-1ll];
        a[i]=min(a[i],a[i]^((1<<k)-1));
    }
    n++;
    sort(a,a+n);
    long long int ans = n*(n-1ll);
    ans/=2ll;
    long long idx=0ll;
    long long cnt=0ll;
    while(idx<n){
        while(idx+cnt<n&&a[idx]==a[idx+cnt])
            cnt++;
        long long c1 = cnt/2ll;
        long long c2 = (cnt+1ll)/2ll;
        ans-=(c1*(c1-1ll))/2ll;
        ans-=(c2*(c2-1ll))/2ll;
        idx=idx+cnt;
        cnt=0ll;
    }
    printf("%lld\n",ans);
}