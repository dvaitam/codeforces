#include<cstdio>
#include<cstring>
#include<algorithm>
#include<queue>
#include<vector>
#include<string>
#define LL long long
#define PII pair<int,int>
#define PIII pair<PII,int>
#define MP make_pair
using namespace std;
int n,m,x,y;
bool vis[200010];
PIII p[200010];
void read(int &x)
{
	char ch=getchar();int mark=1;for(;ch!='-'&&(ch<'0'||ch>'9');ch=getchar());if (ch=='-') mark=-1,ch=getchar();
	for(x=0;ch>='0'&&ch<='9';ch=getchar()) x=x*10+ch-48;x*=mark;
}
int main()
{
	//freopen("a.txt","r",stdin);
	read(n);read(m);
	for(int i=1;i<=m;i++)
	{
		read(x);read(y);p[i]=MP(MP(x,-(x+y-1)),i);
	}
	sort(p+1,p+m+1);
	int now=-1,pre=-1,cnt=m,id=0;
	for(int i=1;i<=m;i++)
		if (p[i].first.first<=pre+1)
		{
			if (now<-p[i].first.second) now=-p[i].first.second,id=p[i].second;
		}
		else 
		{
			if (now>pre) pre=now,cnt--,vis[id]=1,now=-p[i].first.second,id=p[i].second;
			else pre=-p[i].first.second,now=-1,cnt--,vis[p[i].second]=1,id=0;
		}
	if (now>pre) vis[id]=1,cnt--;
	printf("%d\n",cnt);
	for(int i=1;i<=m;i++) if (!vis[i]) printf("%d ",i);printf("\n");
	return 0;
}