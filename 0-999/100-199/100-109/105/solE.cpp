#include<cstdio>
#include<iostream>
#include<cmath>
#include<cstring>
#include<cstdlib>
#include<algorithm>
#include<map>
using namespace std;
bool bst;
int i,j,k,l,n,m,s,an,tot,len,L,P;
class arr
{
    public:
        int d,m,r,son,fa;
        bool f1,f2,f3;
}a[4];
bool ff[47][47][47][8][8];
void dfs(int n1,int n2)
{
	if (ff[a[1].d][a[2].d][a[3].d][n1][n2]) return;
	for (int i=1;i<=3;i++)
	if (a[i].d<=42) an=max(an,a[i].d);
	ff[a[1].d][a[2].d][a[3].d][n1][n2]=true;
	int mm=9999;
	for (int i=1;i<=3;i++)
	mm=min(mm,a[i].d);
	int ju[4];
	memset(ju,0,sizeof(ju));
	for (int i=1;i<=3;i++) if (a[i].d>42) ju[a[i].d-42]=i;
	for (int i=1;i<=3;i++)
	if (a[i].d<=42)
	if (!(n1&(1<<(i-1))))
	if (!ju[i])
	{
		for (int j=max(mm,a[i].d-a[i].m);j<=a[i].d+a[i].m;j++)
		{
			int pp=a[i].d;
			a[i].d=j;
			dfs(n1|(1<<(i-1)),n2);
			a[i].d=pp;
		}
	}
	for (int i=1;i<=3;i++)
	if (a[i].d<=42)
	if (!(n2&(1<<(i-1))))
	if (!ju[i])
	{
		for (int j=1;j<=3;j++)
		if (abs(a[j].d-a[i].d)==1)
		{
			int pp=a[j].d;
			a[j].d=42+i;
			dfs(n1,n2|(1<<(i-1)));
			a[j].d=pp;
		}
	}
	for (int i=1;i<=3;i++)
	if (a[i].d<=42)
	if (ju[i])
	{
		int y=ju[i];
		for (int j=max(mm,a[i].d-a[i].r);j<=a[i].d+a[i].r;j++)
		{
			a[y].d=j;
			dfs(n1,n2);
			a[y].d=42+i;
		}
	}
}
int main()
{
   // freopen("throw.out","w",stdout);
    //freopen("a.in","r",stdin);
    n=3;
    while (scanf("%d%d%d",&a[1].d,&a[1].m,&a[1].r)==3)
    {
		memset(ff,0,sizeof(ff));
       // memset(a,0,sizeof(a));
        for (int i=2;i<=n;i++)
        scanf("%d%d%d",&a[i].d,&a[i].m,&a[i].r);
        an=0;
        dfs(0,0);
        printf("%d\n",an);
    }
    return 0;
}