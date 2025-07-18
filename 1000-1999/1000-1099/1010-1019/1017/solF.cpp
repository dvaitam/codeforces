#include <bits/stdc++.h>

using namespace std;
typedef unsigned int ll;
const ll maxn=4e8;
const ll maxp=sqrt(maxn)+10;
ll f[4][maxp],g[4][maxp];
ll n,m,A,B,C,D,ans,res,tmp,rt;

ll p1(ll x){
    ll y=x+1; if (x%2==0) x/=2; else y/=2;
    return x*y;
}
ll p2(ll x){
    long long y=x+1,z=x*2+1;
    if (x%2==0) x/=2; else y/=2;
    if (z%3==0) z/=3; else if (x%3==0) x/=3; else y/=3;
    return (ll)x*y*z;
}
ll p3(ll x){
    ll y=p1(x);
    return y*y;
}
bool is_prime(ll x){
    for (ll i=2;i*i<=x;i++) if (x%i==0) return false;
    return true;
}

ll solve(ll n)
{
    ll i,j,rt,o[4];
    for (m=1;m*m<=n;m++) {
        f[0][m]=(n/m-1);
        f[1][m]=(p1(n/m)-1)*C;
        f[2][m]=(p2(n/m)-1)*B;
        f[3][m]=(p3(n/m)-1)*A;
    }
    for (i=1;i<=m;i++) {
        g[0][i]=(i-1);
        g[1][i]=(p1(i)-1)*C;
        g[2][i]=(p2(i)-1)*B;
        g[3][i]=(p3(i)-1)*A;
    }
    for (i=2;i<=m;i++) {
        if (g[0][i]==g[0][i-1]) continue;
        o[0]=1; for (int w=1;w<4;w++) o[w]=o[w-1]*i;
        for (j=1;j<=min(m-1,n/i/i);j++)
            for (int w=0;w<4;w++)
                if (i*j<m) f[w][j]-=o[w]*(f[w][i*j]-g[w][i-1]);
                else f[w][j]-=o[w]*(g[w][n/i/j]-g[w][i-1]);
        for (j=m;j>=i*i;j--)
            for (int w=0;w<4;w++)
                g[w][j]-=o[w]*(g[w][j/i]-g[w][i-1]);
    }
    for (int i=1;i<=m+1;i++) f[0][i]*=D,g[0][i]*=D;
    rt=0;
    for (ll i=1;n/i>m;i++) {
        for (int w=0;w<4;w++) rt+=(f[w][i]-g[w][m]);
        //cout<<"H"<<rt<<endl;
    }
    return rt;
}

int main()
{
    //freopen("input.txt","r",stdin);
    ll n; cin >> n >> A >> B >> C >> D;
    ans=solve(n);
    for (ll i=2;i<=m;i++){
        if (is_prime(i)){
            res=A*i*i*i+B*i*i+C*i+D; tmp=n;
            while (tmp){
                ans+=res*(tmp/i);
                tmp/=i;
            }
        }
    }
    cout << ans << endl;
    return 0;
}