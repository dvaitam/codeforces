// Fresh Peach Heart Shower
#include <bits/stdc++.h>
#define reg
#define ALL(x) (x).begin(),(x).end()
#define mem(x,y) memset(x,y,sizeof x)
#define sz(x) (int)(x).size()
#define ln putchar('\n')
#define lsp putchar(32)
#define pb push_back
#define MP std::make_pair
#define MT std::make_tuple
#ifdef _LOCAL_
#define dbg(x) std::cerr<<__func__<<"\tLine:"<<__LINE__<<' '<<#x<<": "<<x<<"\n"
#define dprintf(x...) std::fprintf(stderr,x)
#else
#define dbg(x) 42
#define dprintf(x...) 42
#endif
#define rep(i,a,b) for(int i=(a);i<=(b);++i)
#define per(i,b,a) for(int i=(b);i>=(a);--i)
typedef std::pair<int,int> pii;
typedef std::vector<int> poly;
template <class t> inline void read(t &s){s=0;
signed f=1;char c=getchar();while(!isdigit(c)){if(c=='-')f=-1;c=getchar();}
while(isdigit(c))s=(s<<3)+(s<<1)+(c^48),c=getchar();;s*=f;}
template<class t,class ...A> inline void read(t &x,A &...a){read(x);read(a...);}
template <class t> inline void write(t x){if(x<0)putchar('-'),x=-x;
static int buf[50],top=0;while(x)buf[++top]=x%10,x/=10;if(!top)buf[++top]=0;
while(top)putchar(buf[top--]^'0');}
inline void setIn(std::string s){freopen(s.c_str(),"r",stdin);return;}
inline void setOut(std::string s){freopen(s.c_str(),"w",stdout);return;}
inline void setIO(std::string s=""){setIn(s+".in");setOut(s+".out");return;}
template <class t>inline bool ckmin(t&x,t y){if(x>y){x=y;return 1;}return 0;}
template <class t>inline bool ckmax(t&x,t y){if(x<y){x=y;return 1;}return 0;}
inline int lowbit(int x){return x&(-x);}
#define fi first
#define se second
#define lson (u<<1)
#define rson (u<<1|1)
typedef long long ll;
const int MaxN=2e5+50;
struct Node
{
	int C,cnt;
	ll S,sum;
	int tag;
	inline void flip(){cnt=C-cnt,sum=S-sum,tag^=1;}
	inline Node operator + (const Node &nt) const
	{
		Node ret;ret.tag=0;
		ret.C=C+nt.C,ret.cnt=cnt+nt.cnt;
		ret.S=S+nt.S,ret.sum=sum+nt.sum;
		return ret;
	}
}a[MaxN<<2];
std::vector<pii > E[MaxN];
int siz[MaxN],gson[MaxN],top[MaxN],L[MaxN],R[MaxN],F[MaxN],dep[MaxN],dw[MaxN],rev[MaxN],dfncnt,n;
inline void dfs1(int u,int fa)
{
	F[u]=fa,dep[u]=dep[fa]+1,siz[u]=1;
	for(auto v:E[u])if(v.fi!=fa)
	{
		dfs1(v.fi,u),dw[v.fi]=v.se,siz[u]+=siz[v.fi];
		if(siz[v.fi]>siz[gson[u]])gson[u]=v.fi;
	}
}
inline void dfs2(int u,int ftop)
{
	top[u]=ftop,rev[L[u]=++dfncnt]=u;
	if(!gson[u])return R[u]=dfncnt,void();
	dfs2(gson[u],ftop);
	for(auto v:E[u])if(v.fi!=F[u]&&v.fi!=gson[u])dfs2(v.fi,v.fi);
	R[u]=dfncnt;
}
inline void pushup(int u){a[u]=a[lson]+a[rson];}
inline void pushdown(int u){if(a[u].tag)a[lson].flip(),a[rson].flip(),a[u].tag=0;}
inline void upd(int u,int l,int r,int p)
{
	if(l==r)return a[u].C=1,a[u].S=dw[rev[l]],a[u].sum=a[u].S,void();
	pushdown(u);
	int mid=(l+r)>>1;
	p<=mid?upd(lson,l,mid,p):upd(rson,mid+1,r,p),pushup(u);
}
inline void Flip(int u,int l,int r,int ql,int qr)
{
	if(ql<=l&&r<=qr)return a[u].flip();
	pushdown(u);
	int mid=(l+r)>>1;
	if(ql<=mid)Flip(lson,l,mid,ql,qr);
	if(mid<qr)Flip(rson,mid+1,r,ql,qr);
	pushup(u);
}
poly ans;
int lf[MaxN],yes[MaxN];
inline int dfs3(int u,int fa)
{
	if(lf[u])return ans.pb(dw[u]),1;
	int ret=0;
	for(auto v:E[u])if(yes[v.fi]&&v.fi!=fa)ret+=dfs3(v.fi,u);
	return ret?0:(ans.pb(dw[u]),1);
}
inline void solve()
{
	ans.clear(),dfs3(1,0);
	write((int)ans.size()),lsp;
	std::sort(ALL(ans));
	for(auto i:ans)write(i),lsp;ln;
}
signed main(void)
{
	read(n);
	int u,v;
	rep(i,1,n-1)read(u,v),E[u].pb(MP(v,i)),E[v].pb(MP(u,i));
	dfs1(1,0),dfs2(1,1);
	lf[1]=1,upd(1,1,n,1);
	int opt,x,ok=0,N=1;
	while(true)
	{
		read(opt);
		if(opt==1)
		{
			read(x),upd(1,1,n,L[x]),++N;
			lf[F[x]]=0,lf[x]=1,yes[x]=1,x=F[x];
			while(top[x]!=1)Flip(1,1,n,L[top[x]],L[x]),x=F[top[x]];
			Flip(1,1,n,1,L[x]);
			ok=a[1].cnt*2==N;
			write(ok?a[1].sum:0),ln;
		}
		else if(opt==2)
		{
			if(!ok)std::puts("0");
			else solve();
		}
		else break;
		fflush(stdout);
	}
	return 0;
}

/*
 * Check List:
 * 1. Input / Output File (OI)
 * 2. long long 
 * 3. Special Test such as n=1
 * 4. Array Size
 * 5. Memory Limit (OI) int is 4 and longlong is 8
 * 6. Mod (a*b%p*c%p not a*b*c%p  ,  (a-b+p)%p not a-b )
 * 7. Name ( int k; for(int k...))
 * 8. more tests , (T=2 .. more)
 * 9. blank \n after a case
*/