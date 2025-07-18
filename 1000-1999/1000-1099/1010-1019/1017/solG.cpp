#include<bits/stdc++.h>

typedef unsigned int uint;
typedef long long ll;
typedef unsigned long long ull;
typedef double lf;
typedef long double llf;
typedef std::pair<int,int> pii;

#define xx first
#define yy second

template<typename T> inline T max(T a,T b){return a>b?a:b;}
template<typename T> inline T min(T a,T b){return a<b?a:b;}
template<typename T> inline T abs(T a){return a>0?a:-a;}
template<typename T> inline bool repr(T &a,T b){return a<b?a=b,1:0;}
template<typename T> inline bool repl(T &a,T b){return a>b?a=b,1:0;}
template<typename T> inline T gcd(T a,T b){T t;if(a<b){while(a){t=a;a=b%a;b=t;}return b;}else{while(b){t=b;b=a%b;a=t;}return a;}}
template<typename T> inline T sqr(T x){return x*x;}
#define mp(a,b) std::make_pair(a,b)
#define pb push_back
#define I inline
#define mset(a,b) memset(a,b,sizeof(a))
#define mcpy(a,b) memcpy(a,b,sizeof(a))

#define fo0(i,n) for(int i=0,i##end=n;i<i##end;i++)
#define fo1(i,n) for(int i=1,i##end=n;i<=i##end;i++)
#define fo(i,a,b) for(int i=a,i##end=b;i<=i##end;i++)
#define fd0(i,n) for(int i=(n)-1;~i;i--)
#define fd1(i,n) for(int i=n;i;i--)
#define fd(i,a,b) for(int i=a,i##end=b;i>=i##end;i--)
#define foe(i,x)for(__typeof((x).end())i=(x).begin();i!=(x).end();++i)

struct Cg{I char operator()(){return getchar();}};
struct Cp{I void operator()(char x){putchar(x);}};
#define OP operator
#define RT return *this;
#define UC unsigned char
#define RX x=0;UC t=P();while((t<'0'||t>'9')&&t!='-')t=P();bool f=0;\
if(t=='-')t=P(),f=1;x=t-'0';for(t=P();t>='0'&&t<='9';t=P())x=x*10+t-'0'
#define RL if(t=='.'){lf u=0.1;for(t=P();t>='0'&&t<='9';t=P(),u*=0.1)x+=u*(t-'0');}if(f)x=-x
#define RU x=0;UC t=P();while(t<'0'||t>'9')t=P();x=t-'0';for(t=P();t>='0'&&t<='9';t=P())x=x*10+t-'0'
#define TR *this,x;return x;
I bool IS(char x){return x==10||x==13||x==' ';}template<typename T>struct Fr{T P;I Fr&OP,(int&x)
{RX;if(f)x=-x;RT}I OP int(){int x;TR}I Fr&OP,(ll &x){RX;if(f)x=-x;RT}I OP ll(){ll x;TR}I Fr&OP,(char&x)
{for(x=P();IS(x);x=P());RT}I OP char(){char x;TR}I Fr&OP,(char*x){char t=P();for(;IS(t);t=P());if(~t){for(;!IS
(t)&&~t;t=P())*x++=t;}*x++=0;RT}I Fr&OP,(lf&x){RX;RL;RT}I OP lf(){lf x;TR}I Fr&OP,(llf&x){RX;RL;RT}I OP llf()
{llf x;TR}I Fr&OP,(uint&x){RU;RT}I OP uint(){uint x;TR}I Fr&OP,(ull&x){RU;RT}I OP ull(){ull x;TR}};Fr<Cg>in;
#define WI(S) if(x){if(x<0)P('-'),x=-x;UC s[S],c=0;while(x)s[c++]=x%10+'0',x/=10;while(c--)P(s[c]);}else P('0')
#define WL if(y){lf t=0.5;for(int i=y;i--;)t*=0.1;if(x>=0)x+=t;else x-=t,P('-');*this,(ll)(abs(x));P('.');if(x<0)\
x=-x;while(y--){x*=10;x-=floor(x*0.1)*10;P(((int)x)%10+'0');}}else if(x>=0)*this,(ll)(x+0.5);else *this,(ll)(x-0.5);
#define WU(S) if(x){UC s[S],c=0;while(x)s[c++]=x%10+'0',x/=10;while(c--)P(s[c]);}else P('0')
template<typename T>struct Fw{T P;I Fw&OP,(int x){WI(10);RT}I Fw&OP()(int x){WI(10);RT}I Fw&OP,(uint x){WU(10);RT}
I Fw&OP()(uint x){WU(10);RT}I Fw&OP,(ll x){WI(19);RT}I Fw&OP()(ll x){WI(19);RT}I Fw&OP,(ull x){WU(20);RT}I Fw&OP()
(ull x){WU(20);RT}I Fw&OP,(char x){P(x);RT}I Fw&OP()(char x){P(x);RT}I Fw&OP,(const char*x){while(*x)P(*x++);RT}
I Fw&OP()(const char*x){while(*x)P(*x++);RT}I Fw&OP()(lf x,int y){WL;RT}I Fw&OP()(llf x,int y){WL;RT}};Fw<Cp>out;

const int N=100007;

std::vector<int>p[N];
int n,m,idm,fa[N],dep[N],sz[N],ws[N],top[N],id[N],idr[N],ir[N],bv[N];

inline void dfs(int x)
{
	sz[x]=1;int y=0;
	foe(i,p[x])
	{
		dfs(*i);
		sz[x]+=sz[*i];
		if(sz[*i]>y)y=sz[*i],ws[x]=*i;
	}
}

inline void dfs2(int x,int y)
{
	top[x]=y;
	id[x]=++idm;
	if(ws[x])dfs2(ws[x],y);
	foe(i,p[x])if(*i!=ws[x])
		dfs2(*i,*i);
	idr[x]=idm;
}

struct node
{
	int mi,ui,ua,tag;
	//mi=min dep[i]+sum(i to root)
	//ui=min sum(i to root)
	//ua=min sum(i to root)
}f[262333];

inline void st(int x,int v)
{
	f[x].mi+=v;
	f[x].tag+=v;
}

inline void pd(int x)
{
	if(f[x].tag)
	{
		st(x<<1,f[x].tag);
		st(x<<1|1,f[x].tag);
		f[x].tag=0;
	}
}

inline void pu(int x)
{
	f[x].mi=min(f[x<<1].mi,f[x<<1|1].mi);
	f[x].ui=min(f[x<<1].ui,f[x<<1|1].ui);
	f[x].ua=max(f[x<<1].ua,f[x<<1|1].ua);
}

inline void build(int x,int l,int r)
{
	if(l==r)f[x].mi=bv[l];
	else
	{
		int t=(l+r)>>1;
		build(x<<1,l,t);
		build(x<<1|1,t+1,r);
		pu(x);
	}
}

inline void _mod(int x,int l,int r,int ql,int qr,int v)
{
	if(l>=ql&&r<=qr)return st(x,v);
	int t=(l+r)>>1;pd(x);
	if(ql<=t)_mod(x<<1,l,t,ql,qr,v);
	if(qr>t)_mod(x<<1|1,t+1,r,ql,qr,v);
	pu(x);
}

inline void mod(int x,int l,int r,int ql,int qr,int v)
{
	if(ql<=qr)_mod(x,l,r,ql,qr,v);
}

inline void modu(int x,int l,int r,int p,int v)
{
	if(l==r)
	{
		f[x].ui+=v;
		f[x].ua+=v;
		return;
	}
	int t=(l+r)>>1;pd(x);
	if(p<=t)modu(x<<1,l,t,p,v);
	else modu(x<<1|1,t+1,r,p,v);
	pu(x);
}

int qa;

inline void _qry(int x,int l,int r,int ql,int qr)
{
	if(l>=ql&&r<=qr){repl(qa,f[x].mi);return;}
	int t=(l+r)>>1;pd(x);
	if(ql<=t)_qry(x<<1,l,t,ql,qr);
	if(qr>t)_qry(x<<1|1,t+1,r,ql,qr);
}

inline int qry(int x,int l,int r,int ql,int qr)
{
	qa=1e9;
	_qry(x,l,r,ql,qr);
	return qa;
}

inline int qryu(int x,int l,int r,int p)
{
	if(l==r)return f[x].ua;
	int t=(l+r)>>1;
	if(p<=t)return qryu(x<<1,l,t,p);
	return qryu(x<<1|1,t+1,r,p);
}

pii ru[N];
int rc;

inline void rst(int x,int l,int r,int ql,int qr)
{
	if(f[x].ua==0&&f[x].ui==0)return;
	if(l==r)
	{
		ru[rc++]=mp(ir[l],f[x].ua);
		f[x].ua=f[x].ui=0;
		return;
	}
	int t=(l+r)>>1;pd(x);
	if(ql<=t)rst(x<<1,l,t,ql,qr);
	if(qr>t)rst(x<<1|1,t+1,r,ql,qr);
	pu(x);
}

int main()
{
	in,n,m;
	fo(i,2,n)in,fa[i];
	fo(i,2,n)p[fa[i]].pb(i);
	dfs(1);
	dfs2(1,1);
	fo1(i,n)ir[id[i]]=i;
	fo(i,2,n)dep[i]=dep[fa[i]]+1;
	fo1(i,n)bv[id[i]]=-dep[i];
	build(1,1,n);
	while(m--)
	{
		int op,x;
		in,op,x;
		if(op==1)
		{
			modu(1,1,n,id[x],1);
			mod(1,1,n,id[x]+1,idr[x],1);
		}
		else if(op==2)
		{
			rc=0;
			rst(1,1,n,id[x],idr[x]);
			fo0(i,rc)mod(1,1,n,id[ru[i].xx]+1,idr[ru[i].xx],-ru[i].yy);
			int y=x,r=1e9;
			while(y)
			{
				repl(r,qry(1,1,n,id[top[y]],id[y]));
				y=fa[top[y]];
			}
			r-=qry(1,1,n,id[x],id[x]);
			modu(1,1,n,id[x],r);
			mod(1,1,n,id[x]+1,idr[x],r);
		}
		else
		{
			int v=qry(1,1,n,id[x],id[x])+qryu(1,1,n,id[x]),r=v;
			while(x)
			{
				repl(r,qry(1,1,n,id[top[x]],id[x]));
				x=fa[top[x]];
			}
			out,r<v?"black":"white",'\n';
		}
	}
}