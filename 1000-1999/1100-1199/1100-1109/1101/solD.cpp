//#include<math.h>
#include<algorithm>
#include<stdlib.h>
#include<time.h>
#include<stdio.h>
#include<string.h>
#define srd srand(time(0))
#define ll long long
#define con continue
#define gtc getchar()
#define ptc putchar
#define dou double
#define eps 0.00000000001
#define opr operator
#define cl(x,a) memset(x,a,sizeof(x))
#define fo0(i,k) for(i=fr[k];i;i=nx[i])
#define fo1(i,l,r) for(i=l;i<=r;i++)
#define fo2(i,l,r) for(i=l;i>=r;i--)
#define fo(i,n) for(i=1;i<=n;i++)
#define ret return
#define x first
#define cint const int
#define y second
#define opi(x) freopen(x,"r",stdin)
#define opo(x) freopen(x,"w",stdout)
#define tpl template<class T>
#define priq priority_queue
#define mp make_pair
#define use using namespace
#define WT while(T--)
use std;
typedef pair<int,int> pii;typedef pair<int,ll> pil;typedef pair<ll,int> pli;typedef pair<ll,ll> pll;
namespace io
{
	void _(int &k){char c;int e=1;k=0;while((c=gtc)>'9'||c<'0')if(c=='-')e=-1;k=c-'0';while((c=gtc)<='9'&&c>='0'){k*=10;k+=c-'0';}k*=e;}
	void _(ll &k){char c;int e=1;k=0;while((c=gtc)>'9'||c<'0')if(c=='-')e=-1;k=c-'0';while((c=gtc)<='9'&&c>='0'){k*=10;k+=c-'0';}k*=e;}
	void _(char &c){while((c=gtc)==' '||c=='\n');}void _(dou &c){scanf("%lf",&c);}void _(char *s){char c;while((c=gtc)!=EOF&&c!=' '&&c!=10)*s++=c;}
	template<class t1,class t2>void _(t1 &a,t2 &b){_(a);_(b);}template<class t1,class t2,class t3>void _(t1 &a,t2 &b,t3 &c){_(a);_(b);_(c);}
	template<class t1,class t2,class t3,class t4>void _(t1 &a,t2 &b,t3 &c,t4 &d){_(a);_(b);_(c);_(d);}
	template<class t1,class t2,class t3,class t4,class t5>void _(t1 &a,t2 &b,t3 &c,t4 &d,t5 &e){_(a);_(b);_(c);_(d);_(e);}
	void _p(dou k){printf("%.6lf",k);}void _p(char *c){for(;*c;ptc(*c++));}void _p(const char *c){for(;*c;ptc(*c++));}void _p(char c){ptc(c);}
	tpl void _p0(T k){if(k>=10)_p0(k/10);ptc(k%10+'0');}tpl void _p(T k){if(k<0){ptc('-');_p0(-k);}else _p0(k);}tpl void __p(T k){_p(k);ptc(' ');}
	tpl void _pn(T k){_p(k);ptc('\n');}template<class t1,class t2>void _p(t1 a,t2 b){__p(a);_pn(b);}
	template<class t1,class t2,class t3>void _p(t1 a,t2 b,t3 c){__p(a);__p(b);_pn(c);}
	template<class t1,class t2,class t3,class t4>void _p(t1 a,t2 b,t3 c,t4 d){__p(a);__p(b);__p(c);_pn(d);}
	tpl void op(T *a,int n){int i;n--;fo(i,n)__p(a[i]);_pn(a[n+1]);}int gi(){int x;_(x);ret x;}ll gll(){ll x;_(x);ret x;}
}
int gcd(int a,int b){ret b?gcd(b,a%b):a;}void fcl(){fclose(stdin);fclose(stdout);}
void fop(const char *s){char c[256],d[256];cl(c,0);cl(d,0);strcpy(c,s);strcpy(d,s);opi(strcat(c,".in"));opo(strcat(d,".out"));}
int eq(dou a,dou b){return a+eps>=b&&b+eps>=a;}tpl void _ma(T &a,T b){if(a<b)a=b;}tpl void _mi(T &a,T b){if(a>b)a=b;}
cint N=444444,EE=100000000,GG=1000000000,ima=2147483647;
use io;
int n,m,a[N],f[N][10],f2[N][10],g0[N][10],g[N][10],z[N],z0[N],v[N],an,ma,T,j2,fr[N],nx[N],t[N];
void me(int s1,int t1)
{
	nx[++j2]=fr[s1];
	fr[s1]=j2;
	t[j2]=t1;
}
void dfs(int k)
{
	int i,j,k1,y,t1,m=z[k];
	v[k]=1;
	fo0(i,k)
		if(!v[y=t[i]])
		{
			dfs(y);
			fo(j,m)
				fo(k1,z[y])
					if(g[y][k1]==g[k][j])
						if(f[y][k1]+1>f[k][j])
						{
							f2[k][j]=f[k][j];
							f[k][j]=f[y][k1]+1;
						}
						else if(f[y][k1]+1>f2[k][j])
							f2[k][j]=f[y][k1]+1;
		}
	fo(i,m)
		_ma(an,f[k][i]+f2[k][i]+1);
}
int main()
{
	int i,j,a1,a2;
	_(n);
	fo(i,n)
		_(a[i]),_ma(ma,a[i]);
	fo1(i,2,ma)
		if(!z0[i])
			for(j=i;j<=ma;j+=i)
				g0[j][++z0[j]]=i;
	fo(i,n)
	{
		z[i]=z0[a[i]];
		fo(j,z[i])
			g[i][j]=g0[a[i]][j];
	}
	fo(i,n-1)
	{
		_(a1,a2);
		me(a1,a2);
		me(a2,a1);
	}
	dfs(1);
	_pn(an);
}