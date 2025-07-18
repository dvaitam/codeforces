// Fresh Peach Heart Shower

#include <bits/stdc++.h>

#define reg

#define ALL(x) (x).begin(),(x).end()

#define mem(x,y) memset(x,y,sizeof x)

#define ln putchar('\n')

#define lsp putchar(32)

#define pb push_back

#define MP std::make_pair

#ifdef _LOCAL_

#define dprintf(x...) std::fprintf(stderr,x)

#else

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

const int MaxN=3e5+50;

int a[MaxN],f[MaxN],pre[MaxN],n;

struct ST

{

	int up[21][MaxN],LG[MaxN];

	inline void init()

	{

		rep(i,2,n)LG[i]=LG[i>>1]+1;

		rep(i,1,n)up[0][i]=a[i]+i;

		rep(k,1,20)for(int i=1;i+(1<<k)-1<=n;++i)

			up[k][i]=std::max(up[k-1][i],up[k-1][i+(1<<~-k)]);

	}

	inline int ask(int l,int r)

	{

		if(l>r)return 0;

		int k=LG[r-l+1];

		return std::max(up[k][l],up[k][r-(1<<k)+1]);

	}

}w;

char ans[MaxN];

inline void work()

{

	read(n);

	rep(i,1,n)read(a[i]),f[i]=0,ans[i]='R';

	w.init();

	rep(i,1,n)

	{

		f[i]=f[i-1],pre[i]=i-1;

		if(f[i]>=i&&ckmax(f[i],a[i]+i))pre[i]=+(i-1);

		int it=std::lower_bound(f,f+i,i-a[i]-1)-f;

		if(it<i&&ckmax(f[i],std::max({i-1,f[it],w.ask(it+1,i-1)})))pre[i]=-(it);

	}

	if(f[n]<n)return std::puts("NO"),void();

	std::puts("YES");

	for(int x=n;x;)

		if(pre[x]<=0)ans[x]='L',x=-pre[x];

		else x=pre[x];

	ans[n+1]='\0',std::puts(ans+1);

}

signed main(void)

{

	int t;read(t);

	while(t--)work();

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