#include<cstdio>
#include<cstring>
#include<algorithm>
using namespace std;


const int MAXN=200000+20;
void Read(int &p)
{
	p=0;
	char c=getchar();
	while(c<'0' || c>'9') c=getchar();
	while(c>='0' && c<='9')
		p=p*10+c-'0',c=getchar();
}
int N,pre,pos,dp,ans=-1;

int main()
{
	Read(N);
	for(int i=1;i<=N;i++)
	{
		Read(pos);
		if(pos<=pre*2) dp++;
		else dp=0;
		ans=max(ans,dp);
		pre=pos;
	}
	printf("%d\n",ans+1);
}