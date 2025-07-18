// LUOGU_RID: 157850466
#include<bits/stdc++.h>
#define F(i,a,b) for(int i=a,i##end=b;i<=i##end;i++)
using namespace std;typedef long long ll;
#define gc() getchar()
int read() {
	int s=0,w=0;char ch=gc();
	while(ch<'0'||ch>'9') w|=(ch=='-'),ch=gc();
	while(ch>='0'&&ch<='9') s=(s<<3)+(s<<1)+(ch^48),ch=gc();
	return w?-s:s;
} const int N=1e5+5;
int n,R,f[N],h[N],ans,s[N],t,col;
struct E {int v;bool b;};
basic_string<E> G[N],uG[N];
void d(int x,int fa=0) {
    for(auto [v,b]:G[x]) if(v^fa) {
        d(v,x);f[x]+=f[v]+1;
    }
    for(auto [v,b]:uG[x]) if(v^fa) {
        d(v,x);f[x]+=f[v]+!b;
    }
}
void d2(int x,int fa=0) {
    f[x]>ans&&(R=x,ans=f[x]);
    for(auto [v,b]:G[x]) if(v^fa) {
        f[v]=f[x]-1+(!b);
        d2(v,x);
    }
    for(auto [v,b]:uG[x]) if(v^fa) {
        f[v]=f[x]-(!b)+1;
        d2(v,x);
    }
}\
void d3(int x,int fa=0) {
    for(auto [v,b]:G[x]) if(v^fa) {
        printf("%d %d %d\n",x,v,s[++t]=++col);
        d3(v,x);--t;
    }
    for(auto [v,b]:uG[x]) if(v^fa) {
        if(!b) {
            printf("%d %d %d\n",x,v,s[++t]=++col);
            d3(v,x);--t;
        } else {int tmp=s[t];
            printf("%d %d %d\n",v,x,s[t--]),d3(v,x),s[++t]=tmp;
        }
}
    }
int main() {
    F(i,2,n=read()) {
        int x=read(),y=read(),b=read();
        G[x]+=E{y,b};uG[y]+=E{x,b};
    } d(1);d2(1);
    cout<<f[R]<<'\n';
    d3(R);
	return 0;
}