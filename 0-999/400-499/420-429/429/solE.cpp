#include<cstdio>
#include<algorithm>
using namespace std;
const int N=222222;
int n,p[N],c[N],a,b,i;
pair<int,int> s[N];
int dp(int x,int t)
{
	if(c[x])return c[x];
	c[x]=t;
	dp(p[x<<1]>>1,p[x<<1]&1^t^1);
	dp(p[x<<1|1]>>1,p[x<<1|1]&1^t);
	return t;
}
int main()
{
	scanf("%d",&n);
	for(i=0;i<n;i++)
	{
		scanf("%d%d",&a,&b);
		s[i<<1]=make_pair(a,i<<1);
		s[i<<1|1]=make_pair(b+1,i<<1|1);
	}
	sort(s,s+n+n);
	for(i=0;i<n;i++)
	{
		a=s[i<<1].second,b=s[i<<1|1].second;
		p[a]=b,p[b]=a;
	}
	for(int i=0;i<n;i++)
		printf("%d ",dp(i,2)^2);
	return 0;
}