#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
#define sqr(x) ((x)*(x))
#define mp make_pair
inline char gc(){
    static char buf[100000],*p1=buf,*p2=buf;
    return p1==p2&&(p2=(p1=buf)+fread(buf,1,100000,stdin),p1==p2)?EOF:*p1++;
}
#define gc getchar
inline ll read(){
	ll x = 0; char ch = gc(); bool positive = 1;
	for (; !isdigit(ch); ch = gc())	if (ch == '-')	positive = 0;
	for (; isdigit(ch); ch = gc())	x = x * 10 + ch - '0';
	return positive ? x : -x;
}
inline void write(ll a){
    if(a<0){
    	a=-a; putchar('-');
	}
    if(a>=10)write(a/10);
    putchar('0'+a%10);
}
inline void wri(ll a){write(a); putchar(' ');}
inline void writeln(ll a){write(a); puts("");}
const int N=100005;
int n,m,k,sz[N],fa[N];
pair<int,pair<int,int> > a[N];
inline int gf(int x){
	return fa[x]==x?x:fa[x]=gf(fa[x]);
}
signed main(){
	cin>>n>>m>>k;
	for(int i=1;i<=k;i++)sz[read()]=1;
	for(int i=1;i<=n;i++)fa[i]=i;
	for(int i=1;i<=m;i++){
		int x=read(),y=read(),z=read();
		a[i]=mp(z,mp(x,y));
	}
	sort(&a[1],&a[m+1]);
	for(int i=1;i<=m;i++){
		int x=a[i].second.first,y=a[i].second.second,z=a[i].first;
		int t1=gf(x),t2=gf(y);
		if(t1!=t2){
			fa[t1]=t2; sz[t2]+=sz[t1];
			if(sz[t2]==k){
				for(int i=1;i<=k;i++)wri(z); puts(""); return 0;
			}
		}
	}
}
/*
考虑最后一条边连接两个联通块 
*/