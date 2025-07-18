#include<set>
#include<map>
#include<queue>
#include<cmath>
#include<string>
#include<cstdio>
#include<vector>
#include<cstring>
#include<iostream>
#include<algorithm>
#define rep(i,a,b) for (int i=a; i<=b; i++)
#define per(i,a,b) for (int i=a; i>=b; i--)
using namespace std;
typedef long long LL;

inline int read() {
    int x=0,f=1; char ch=getchar();
    while (!(ch>='0'&&ch<='9')) {if (ch=='-')f=-1;ch=getchar();}
    while (ch>='0'&&ch<='9') {x=x*10+(ch-'0'); ch=getchar();}
    return x*f;
}

const int N = 100005;
const int INF = 2100000000;

int n,l,r;
int a[N],p[N],b[N];
int mn,mx;

int main() {

	#ifndef ONLINE_JUDGE
	//	freopen("data.in","r",stdin);
	//	freopen("data.out","w",stdout);
	#endif

	n=read(),l=read(),r=read();
	rep(i,1,n) a[i]=read();
	rep(i,1,n) p[i]=read();
	rep(i,1,n) b[i]=a[i]+p[i];
	mn=INF; mx=-INF;
	rep(i,1,n) mn=min(mn,b[i]),mx=max(mx,b[i]);
	if (mx-mn<=r-l) {
		int delta=l-mn;
		rep(i,1,n) b[i]+=delta;
		rep(i,1,n-1) printf("%d ",b[i]); printf("%d\n",b[n]);
	} else printf("-1\n");

	return 0;
}