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
#define I __attribute__((always_inline))inline
#define mset(a,b) memset(a,b,sizeof(a))
#define mcpy(a,b) memcpy(a,b,sizeof(a))

#define fo0(i,n) for(int i=0,i##end=n;i<i##end;i++)
#define fo1(i,n) for(int i=1,i##end=n;i<=i##end;i++)
#define fo(i,a,b) for(int i=a,i##end=b;i<=i##end;i++)
#define fd0(i,n) for(int i=(n)-1;~i;i--)
#define fd1(i,n) for(int i=n;i;i--)
#define fd(i,a,b) for(int i=a,i##end=b;i>=i##end;i--)
#define foe(i,x)for(__typeof((x).end())i=(x).begin();i!=(x).end();++i)
#define fre(i,x)for(__typeof((x).rend())i=(x).rbegin();i!=(x).rend();++i)

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

std::mt19937 ran(11);

int n,m,k,p,c[105],id[105][105],ir[105][105],co[105][105];
bool inq[105];
std::queue<int>q;
pii e[1005];

struct seg
{
	int mi[231],ma[231],v[105];
	inline void init(){mset(mi,0x3f);fo1(i,k)mi[i+p]=0;}
	inline void modify(int x,int y)
	{
		mi[x+p]=ma[x+p]=v[x]+=y;
		for(x+=p;x^1;x>>=1)
			mi[x>>1]=min(mi[x],mi[x^1]),
			ma[x>>1]=max(ma[x],ma[x^1]);
	}
}f[105];

inline bool chk(int x)
{
	return f[x].ma[1]-f[x].mi[1]>2;
}

inline void add(int x)
{
	if(!inq[x]&&chk(x))q.push(x);
}

inline void _modify(int x,int y,int v)
{
	f[x].modify(co[x][id[x][y]],-1);
	co[x][id[x][y]]=v;
	f[x].modify(co[x][id[x][y]],1);
}

inline void modify(int x,int y,int v)
{
	_modify(x,y,v),_modify(y,x,v);
	add(x),add(y);
}

inline void work(int x)
{
	if(!chk(x))return;
	int mi=f[x].mi[1],ma=f[x].ma[1];
	//out,'/',mi,' ',ma,'\n';
	if(ran()&1)
	{
		int tc=0,ta,tb;
		fo1(i,c[x])if(f[x].v[co[x][i]]>=mi+2)tc++;
		tc=ran()%tc;
		fo1(i,c[x])if(f[x].v[co[x][i]]>=mi+2){if(!tc--){ta=i;break;}}
		tc=0;
		fo1(i,k)if(f[x].v[i]==mi)tc++;
		tc=ran()%tc;
		fo1(i,k)if(f[x].v[i]==mi){if(!tc--){tb=i;break;}}
		modify(x,ir[x][ta],tb);
	}
	else
	{
		int tc=0,ta,tb;
		fo1(i,c[x])if(f[x].v[co[x][i]]==ma)tc++;
		tc=ran()%tc;
		fo1(i,c[x])if(f[x].v[co[x][i]]==ma){if(!tc--){ta=i;break;}}
		tc=0;
		fo1(i,k)if(f[x].v[i]<=ma-2)tc++;
		tc=ran()%tc;
		fo1(i,k)if(f[x].v[i]<=ma-2){if(!tc--){tb=i;break;}}
		modify(x,ir[x][ta],tb);
	}
}

int main()
{
	//freopen("in.txt","r",stdin);
	in,n,m,k;k=min(n,k);
	for(p=1;p<k+2;p<<=1);
	fo1(i,n)(int)in;
	fo1(i,m)
	{
		int x,y;
		in,x,y;
		id[x][y]=id[y][x]=1;
		co[x][y]=co[y][x]=ran()%k+1;
		e[i]=mp(x,y);
	}
	fo1(i,n)f[i].init();
	fo1(i,n)fo1(j,k)f[i].modify(j,0);
	fo1(i,n)fo1(j,n)if(id[i][j])
	{
		id[i][j]=++c[i];
		ir[i][c[i]]=j;
		co[i][c[i]]=co[i][j];
		f[i].modify(co[i][j],1);
	}
	fo1(i,n)add(i);
	while(!q.empty())
	{
		int x=q.front();q.pop();
		work(x);
	}
	//fo1(i,n)out,c[i],' ';out,'\n';
	//fo1(i,4)out,co[6][i],' ';out,'\n';
	//fo1(i,k)out,f[6].v[i],' ';out,'\n';
	//out,"==============\n";
	fo1(i,m)
	{
		out,co[e[i].xx][id[e[i].xx][e[i].yy]],'\n';
		//out,e[i].xx,' ',e[i].yy,' ',co[e[i].xx][id[e[i].xx][e[i].yy]],'\n';
	}
}