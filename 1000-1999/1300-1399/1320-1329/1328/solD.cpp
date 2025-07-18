//Head File
#include<cstdio>
#include<algorithm>
#include<cstring>
#include<map>
#include<ctime>
#include<set>
#include<cstdlib>
#include<vector>
#include<cassert>
#include<cmath>
#include<queue>
#include<iostream>
using namespace std;

//Stand
#define il inline
#define ll long long

//Debug
#define B cerr<<"Break Point"<<endl;
#define P(x) cerr<<#x<<' '<<"="<<' '<<(x)<<endl;
#define p(x) cerr<<#x<<' '<<"="<<' '<<(x)<<' ';

//File
void fio()
{
    #ifndef ONLINE_JUDGE
    freopen("sample.in","r",stdin);
    freopen("sample.out","w",stdout);
    #endif
}
void pti()
{
    double timeuse = clock()*1000.0/CLOCKS_PER_SEC;
    cerr<<"Timeuse "<<timeuse<<"ms"<<endl;
}
void end()
{
    #ifndef ONLINE_JUDGE
    pti();
    #endif
    exit(0);
}

//IO
#define pc(s) putchar(s)
#define say(s) cout<<s<<endl
namespace io
{
    const int SIZ=55;int que[SIZ],op,qr;char ch;
    template<class I>
    il void gi(I &w)
    {
        ch=getchar(),op=1,w=0;
	    while(!isdigit(ch)){if(ch=='-') op=-1;ch=getchar();}
        while(isdigit(ch)){w=w*10+ch-'0';ch=getchar();}w*=op;
    }
    template<typename T,typename... Args>
    il void gi(T& t,Args&... args){gi(t);gi(args...);}
    template<class I>
    il void print(I w)
    {
        qr=0;if(!w) putchar('0');if(w<0) putchar('-'),w=-w;
	    while(w) que[++qr]=w%10+'0',w/=10;
        while(qr) putchar(que[qr--]);
    }
}
using io::gi;
using io::print;

const int N=2e5+5;

int T,n;
int a[N],c[N];

int main()
{
	fio();
	gi(T);
	while(T--)
	{
		gi(n);
		for(int i=1;i<=n;++i) gi(a[i]);
		bool flag=false,tag=true;
		for(int i=2;i<=n;++i) 
		{
			if(a[i]==a[i-1]) flag=true;
			else tag=false;
		}
		if(a[1]==a[n]) flag=true;
		if(tag)
		{
			print(1),pc(10);
			for(int i=1;i<=n;++i) print(1),pc(' ');
		}
		else
		{
			if(flag)
			{
				bool vis=false,op=1;
				print(2),pc(10);
				if(n&1)
				{
					print(2),pc(' ');
					for(int i=2;i<=n;++i)
					{
						if(!vis&&a[i]==a[i-1]) vis=true;
						else op^=1;
						print(op+1),pc(' ');
					}
				}
				else
				{
					print(2),pc(' ');
					for(int i=2;i<=n;++i)
					{
						op^=1;
						print(op+1),pc(' ');
					}
				}
			}
			else
			{
				bool op=1;
				if(n&1)
				{
					print(3),pc(10);
					print(2),pc(' ');
					for(int i=2;i<=n-1;++i)
					{
						op^=1;
						print(op+1),pc(' ');
					}
					print(3);
				}
				else
				{
					print(2),pc(10);
					print(2),pc(' ');
					for(int i=2;i<=n;++i)
					{
						op^=1;
						print(op+1),pc(' ');
					}
				}
			}
		}
		pc(10);
	}
	end();
}