#include<bits/stdc++.h>
#define qmin(x,y) (x=min(x,y))
#define qmax(x,y) (x=max(x,y))
#define pir pair<int,int>
#define mp(x,y) make_pair(x,y)
#define fr first
#define sc second
using namespace std;

char gc() {
//	static char buf[100000],*p1,*p2;
//	return (p1==p2)&&(p2=(p1=buf)+fread(buf,1,100000,stdin))?EOF:*p1++;
	return getchar();
}

template<class T>
int read(T &ans) {
	T f=1;ans=0;
	char ch=gc();
	while(!isdigit(ch)) {
		if(ch==EOF) return EOF;
		if(ch=='-') f=-1;
		ch=gc();
	}
	while(isdigit(ch))
		ans=ans*10+ch-'0',ch=gc();
	ans*=f;return 1;
}

template<class T1,class T2>
int read(T1 &a,T2 &b) {
	return read(a)==EOF?EOF:read(b);
}

template<class T1,class T2,class T3>
int read(T1 &a,T2 &b,T3 &c) {
	return read(a,b)==EOF?EOF:read(c);
}

typedef long long ll;
const int Maxn=1100000;
const int inf=0x3f3f3f3f;
const ll mod=998244353;

int n,m,p[Maxn],u,v,ans,bj[Maxn],flag[Maxn];
int to[Maxn],nxt[Maxn],first[Maxn],tot=1,a[Maxn],num;

inline void add(int u,int v) {
	to[tot]=v;
	nxt[tot]=first[u];
	first[u]=tot++;
}

int main() {
//	freopen("test.in","r",stdin);
	read(n,m);
	for(int i=1;i<=n;i++) read(p[i]);
	for(int i=1;i<=m;i++) {
		read(u,v);
		add(u,v);
		if(v==p[n]) bj[u]=1;
	}
	for(int i=n-1;i>=1;i--)
		if(bj[p[i]]) {
			ans++;
			for(int j=first[p[i]];j;j=nxt[j])
				flag[to[j]]=1;
			for(int j=1;j<=num;j++)
				if(!flag[p[a[j]]]) {
					ans--;
					a[++num]=i;
					break;
				}
			for(int j=first[p[i]];j;j=nxt[j])
				flag[to[j]]=0;
		}
		else a[++num]=i;
	printf("%d\n",ans);
	return 0;
}