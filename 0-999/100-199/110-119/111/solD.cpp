#include <bits/stdc++.h>
#include <iostream>
#include <cstdio>
#include <memory.h>
#define MAX(a , b) ((a) > (b) ? (a) : (b))
using namespace std;
#define ll long long
int mod=1000000007;
int n,m,k;
int c[1005][1005];
void getc() {
    for(int i=0;i<=n;++i)  {
        c[i][0]=1;
        for(int j=1;j<=i;++j){
           c[i][j]=c[i-1][j]+c[i-1][j-1];
           if(c[i][j]>=mod)
             c[i][j]-=mod;
        }
    }
}
ll fast_mod(int64_t js,int64_t cs) {
    int64_t t=js%mod,res=1;
    while(cs) {
        if(cs&1)
           res=res*t%mod;
        t=t*t%mod;
        cs>>=1;
    }
    return res;
}
ll f[1005];
ll fz[2005],fm[1005];
int main()
{
    scanf("%d%d%d",&n,&m,&k);
    if(m==1){
       printf("%lld\n",fast_mod(k,n));
       return 0;
    }
    getc();
    for(int i=1;i<=n;++i){
        f[i]=fast_mod(i,n);
        for(int j=1;j<i;++j)
            f[i]=(f[i]-f[j]*c[i][j]%mod+mod)%mod;
    }
    fz[0]=fm[0]=1;
    for(int i=1;i<=2*n && i<=k;++i)
       fz[i]=fz[i-1]*(k-i+1)%mod;
    for(int i=1;i<=n;++i)
       fm[i]=fm[i-1]*fast_mod(i,mod-2)%mod;
    ll ans=0;
    for(int i=0;i<=n && i<=k;++i){
        ll temp=fast_mod(i,(m-2)*n);
        for(int j=0;j+i<=n && i+2*j<=k;++j){
            ans =(ans + temp* fz[j*2+i]%mod *fm[i]%mod *fm[j]%mod *fm[j]%mod *f[i+j]%mod *f[i+j]%mod);
            if(ans >=mod)
              ans -=mod;
        }
    }
    printf("%lld\n",ans);

    return 0;
}