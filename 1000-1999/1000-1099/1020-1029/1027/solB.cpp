#include <cstdio>
#include <algorithm>
#include <cstring>
#include <cctype>
#include <cmath>
#define ll long long
#define linf (1ll<<60)
#define iinf 0x3f3f3f3f
#define dinf 1e100
#define eps 1e-8
#define cl(x) memset(x,0,sizeof(x))
#define maxn 1000
using namespace std;
ll n, Q;
ll read(ll x=0)
{
	ll c, f=1;
	for(c=getchar();!isdigit(c);c=getchar())if(c=='-')f=-1;
	for(;isdigit(c);c=getchar())x=(x<<1)+(x<<3)+c-48;
	return f*x;
}
int main()
{
	ll i, j, ans;
	n=read(), Q=read();
	while(Q--)
	{
		i=read(), j=read();
//	for(i=1;i<=n;i++)for(j=1;j<=n;j++)
//	{
		if((n&1)==0)
		{
			ans=(i-1)*n/2+j/2+(j&1);
			if((i+j&1))ans+=n*n/2;
		}
		else
		{
			if((i+j&1)==0)
			{
				if(i&1)ans=(i-1)/2*n+j/2+1;
				else ans=(i-2)/2*n+n/2+1+j/2;
			}
			else
			{
				if(i&1)ans=(i-1)/2*n+j/2+n*n/2+1;
				else ans=(i-2)/2*n+n/2+j/2+1+n*n/2+1;
			}
		}
		printf("%I64d\n",ans);
	}
	return 0;
}