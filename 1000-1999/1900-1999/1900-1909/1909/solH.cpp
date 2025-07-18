// LUOGU_RID: 143532239
#include<bits/stdc++.h>
#define rep(i,j,k) for(int i=j;i<=k;i++)
#define pi pair<int,int>
using namespace std;
const int N=3e5+5;
struct IN {
    bool any;
    operator bool() const {return any;}
} CIN;
inline IN& operator>>(IN& in, int& x) {
    x = 0, in.any = 0; char c;
    while ((c = getchar()) > 0) {
        if (c < 48 || c > 57)
            if (in.any) return in; else continue;
        in.any = 1, x = x * 10 + c - '0';
    }
    return in;
}
#define cin CIN
#define istream IN


int n,p[N],id[N],xl1[N],xl2[N];
vector<pi> ans;
void sol(int *a,int opt){
    vector<int> dy[N];
    rep(i,0,n-1) if(a[i]!=i){
        int d1=i,d2=a[i],d3=0,d4=0;
        if(opt) d1=n-1-d1,d2=n-1-d2;
        if(d1>d2) swap(d1,d2);
        if((d2-d1)%2) d3=(d2-d1+1)/2; else d3=(2*n-d2-d1)/2;
        if(d1%2) d3=n+1-d3;
        d4=d3;
        if(d1%2==0) d4+=d1; else d4-=d1+1;
        dy[d3-1].push_back(d4-1),swap(a[i],a[a[i]]);
    }
    rep(i,0,n-1){
        int d1=i%2,d2=n; if((d2-d1)%2) d2--;
        if(d1<d2) ans.push_back({d1,d2});
        for(auto dq:dy[i]) ans.push_back({dq,dq+2});
    } 
}
signed main(){
    ios::sync_with_stdio(false);
    cin>>n;
    rep(i,0,n-1) cin>>p[i],p[i]--,id[p[i]]=i,xl1[i]=i;
    rep(i,0,n-1){
        int zz=i;
        while(p[p[zz]]!=zz) swap(id[p[zz]],id[p[p[p[zz]]]]),swap(xl1[zz],xl1[p[p[zz]]]),swap(p[zz],p[p[p[zz]]]),zz=id[zz];
    }
    rep(i,0,n-1) xl2[i]=p[i];
    sol(xl1,0),sol(xl2,1),cout<<ans.size()<<'\n';
    for(auto [l,r]:ans) cout<<l+1<<' '<<r<<'\n';
}