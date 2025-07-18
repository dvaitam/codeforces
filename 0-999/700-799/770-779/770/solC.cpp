#include <bits/stdc++.h>
#define ll long long
using namespace std;
	int n,k,t,len;
	int l[100005],a[100005],s[100005],vis[100005];
	vector <int> p[100005]; 
	bool found;
inline int read()
{
    int x=0,f=1;char ch=getchar();
    while(ch<'0'||ch>'9'){if(ch=='-')f=-1;ch=getchar();}
    while(ch>='0'&&ch<='9'){x=x*10+ch-'0';ch=getchar();}
    return x*f;
}
void dfs(int x)
{
	for (int i=0;i<s[x]&&!found;i++)
		if (!vis[p[x][i]])
		{
			vis[p[x][i]]=2;
			dfs(p[x][i]);
			vis[p[x][i]]=1;
		}
		else if (vis[p[x][i]]==2)
		{
			found=true;
			break;
		}
	l[++len]=x;
}
int main()
{
	n=read(),k=read();
	for (int i=1;i<=k;i++)
		a[i]=read();
	for (int i=1;i<=n;i++)
	{
		s[i]=read();
		for (int j=1;j<=s[i];j++)
			p[i].push_back(read());
	}
	for (int i=1;i<=k;i++)
	if (!vis[a[i]])
	{
		vis[a[i]]=2;
		dfs(a[i]);
		vis[a[i]]=1;
	}
	if (found)
	{
		puts("-1");
		return 0;
	}
	printf("%d\n",len);
	for (int i=1;i<=len;i++)
		printf("%d ",l[i]);
	return 0;
}