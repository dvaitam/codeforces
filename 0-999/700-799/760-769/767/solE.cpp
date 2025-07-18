#include<cstdio>
#include<algorithm>
#include<queue>
using namespace std;
#define LL long long
struct opt
{
	int v,num;
	bool operator < (const opt &oo) const
	{
		return v>oo.v;
	}
}o1;
priority_queue<opt> pq;
int a[100010],ans1[100010],ans2[100010],w[100010],n;
int main()
{
	LL m,now=0;
	scanf("%d%I64d",&n,&m);
	for (int i=1;i<=n;i++)
	{
		scanf("%d",&a[i]);
		ans1[i]=a[i]/100;
		a[i]%=100;
	}
	for (int i=1;i<=n;i++) scanf("%d",&w[i]);
	for (int i=1;i<=n;i++)
	if (a[i])
	{
		if (m>=a[i])
		{
			m-=a[i];
			pq.push((opt){w[i]*(100-a[i]),i});
		}
		else
		{
			if (pq.empty())
			{
				now+=w[i]*(100-a[i]);
				m+=100-a[i];
				ans2[i]=1;
				continue;
			}
			o1=pq.top();
			if (o1.v<w[i]*(100-a[i]))
			{
				pq.pop();
				pq.push((opt){w[i]*(100-a[i]),i});
				now+=o1.v;
				ans2[o1.num]=1;
				m+=100-a[i];
			}
			else
			{
				now+=w[i]*(100-a[i]);
				m+=100-a[i];
				ans2[i]=1;
			}
		}
	}
	printf("%I64d\n",now);
	for (int i=1;i<=n;i++)
		if (ans2[i]) printf("%d 0\n",ans1[i]+1);
		else printf("%d %d\n",ans1[i],a[i]);
}