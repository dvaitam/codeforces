#include <bits/stdc++.h>
//Was yea ra,rra yea ra synk sphilar yor en me exec hymme METAFALICA waath!

#pragma GCC optimize("Ofast")

#pragma GCC optimize("unroll-loops")

#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,avx,avx2,tune=native")

#include<bits/stdc++.h>

using namespace std;

#define rg register

#define ll long long

#define ull unsigned ll

#define lowbit(x) (x&(-x))

#define djq 998244353

const double eps=1e-8;

const int inf=0x3f3f3f3f;

const ll linf=0x3f3f3f3f3f3f3f3f;

const double alpha=0.73;

inline void file(){

	freopen("1.in","r",stdin);

	freopen("1.out","w",stdout);

}

char buf[1<<21],*p1=buf,*p2=buf;

inline int getc(){

    return p1==p2&&(p2=(p1=buf)+fread(buf,1,(1<<20)+5,stdin),p1==p2)?EOF:*p1++;

}

//#define getc getchar

inline ll read(){

	rg ll ret=0,f=0;char ch=getc();

    while(!isdigit(ch)){if(ch==EOF)exit(0);if(ch=='-')f=1;ch=getc();}

    while(isdigit(ch)){ret=ret*10+ch-48;ch=getc();}

    return f?-ret:ret;

}

inline ll read2(){

	rg ll ret=0,f=0;char ch=getc();

    while(!isdigit(ch)){if(ch==EOF)exit(0);if(ch=='-')f=1;ch=getc();}

    while(isdigit(ch)){ret=ret*2+ch-48;ch=getc();}

    return f?-ret:ret;

}

#define epb emplace_back

#define all(x) (x).begin(),(x).end()

#define fi first

#define se second

#define it iterator

#define mkp make_pair

#define naive return 0*puts("YES")

#define angry return 0*puts("NO")

#define fls fflush(stdout)

#define rep(i,a) for(rg int i=1;i<=a;++i)

#define per(i,a) for(rg int i=a;i;--i)

typedef vector<int> vec;

typedef pair<int,int> pii;

struct point{ int x,y; point(int x=0,int y=0):x(x),y(y) {} inline bool operator<(const point& T)const{ return x^T.x?x<T.x:y<T.y; }; };

inline int ksm(int base,int p){int ret=1;while(p){if(p&1)ret=1ll*ret*base%djq;base=1ll*base*base%djq,p>>=1;}return ret;}



inline void fmtmx(int* A,int n){

	for(rg int i=1;i<n;i<<=1){

		const int len=i<<1;

		for(rg int j=0;j<n;j+=len) for(rg int k=0;k<i;++k) A[j+k]=max(A[j+k],A[i+j+k]);

	}

}

inline void fmtmn(int* A,int n){

	for(rg int i=1;i<n;i<<=1){

		const int len=i<<1;

		for(rg int j=0;j<n;j+=len) for(rg int k=0;k<i;++k) A[j+k]=min(A[j+k],A[i+j+k]);

	}

}

int n,m,k,na,nb,ans,x[1000005],y[1000005];

inline void chg(int& a,int x,int y){

	const int ax=((a>>(k-x))&1),ay=((a>>(k-y))&1);

	a=a-(ax<<(k-x))-(ay<<(k-y))+(ax<<(k-y))+(ay<<(k-x));

}

int a[1<<20],b[1<<20],pos[25];

signed main(){

	n=read(),m=read(),k=read(),na=read2(),nb=read2(); int nna=na;

	memset(b,~0x3f,sizeof(b)),memset(a,0x3f,sizeof(a));

	a[na]=n+1,b[nb]=n+1;

	rep(i,k) pos[i]=i;

	rep(i,n) x[i]=read(),y[i]=read(); 

	per(i,n){

		chg(na,pos[x[i]],pos[y[i]]),chg(nb,pos[x[i]],pos[y[i]]);

		swap(pos[x[i]],pos[y[i]]);

		a[na]=min(a[na],i),b[nb]=max(b[nb],i);

	}

	fmtmn(a,1<<k),fmtmx(b,1<<k);

	for(rg int i=0;i<(1<<k);++i) if(b[i]-a[i]>=m) ans=max(ans,__builtin_popcount(i));

	printf("%d\n",ans*2+k-__builtin_popcount(na)-__builtin_popcount(nb));

	//printf("%d %d\n",__builtin_popcount(na),__builtin_popcount(nb));

	for(rg int i=0;i<(1<<k);++i) if(b[i]-a[i]>=m&&__builtin_popcount(i)==ans) return 0*printf("%d %d\n",a[i],b[i]-1);

	return 0;

}

/*

A B C

k-(A-C)-(B-C)

*/