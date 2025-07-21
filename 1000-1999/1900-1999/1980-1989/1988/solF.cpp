// LUOGU_RID: 171594872
#include<bits/stdc++.h>
#define ll long long
#define rep(i,a,b) for(int i=(a);i<=(b);++i)
using namespace std;
const int N=710,mod=998244353;
int n,a[N],b[N],c[N],C[N][N];
ll f[2][N][N],pre[N][N],suf[N][N],g[N][N],ans;
void upd(ll &x,ll y){x=(x+y)%mod;}
int main(){
	scanf("%d",&n);
    rep(i,1,n) scanf("%d",a+i);rep(i,1,n) scanf("%d",b+i);rep(i,0,n-1) scanf("%d",c+i);
    C[0][0]=1;rep(i,1,n){C[i][0]=1;rep(j,1,i) C[i][j]=(C[i-1][j]+C[i-1][j-1])%mod;}
    f[1][1][0]=1;pre[0][0]=a[1];suf[0][0]=b[1];
    rep(i,1,n){
        int nw=i&1,nx=nw^1;memset(f[nx],0,sizeof(f[nx]));
        rep(j,1,i) rep(k,0,i-1) if(f[nw][j][k]){
            ll val=f[nw][j][k];
            upd(pre[i][k],val*a[j+1]);upd(suf[i][i-1-k],val*b[j+1]);
            upd(f[nx][j+1][k+1],val);upd(f[nx][j][k],val*(k+1));upd(f[nx][j][k+1],val*(i-k-1));
        }
    }
    rep(i,0,n) rep(x,0,i) rep(y,0,n-x-1) upd(g[i][y],pre[i][x]*c[x+y+(i>0)]);
    rep(i,1,n){ans=0;rep(p,1,i) rep(y,0,i-p) upd(ans,g[p-1][y]*suf[i-p][y]%mod*C[i-1][p-1]%mod);printf("%lld ",ans);}
	return 0;
}