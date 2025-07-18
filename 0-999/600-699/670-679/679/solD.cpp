#include<algorithm>
#include <iostream>
#include <stdlib.h>
#include <string.h>
#include  <stdio.h>
#include   <math.h>
#include   <time.h>
#include   <vector>
#include   <bitset>
#include    <queue>
#include      <set>
#include      <map>
using namespace std;

const int N=405;

int n,m,g[N][N],tag1,tag2,vis1[N],vis2[N];
vector<int> G[N],d[N],po,Le;
double Ans,P[N],Maxp[N];

int main()
{
	cin>>n>>m;
	for(int i=1;i<=n;i++)
		for(int j=i+1;j<=n;j++)
			g[i][j]=g[j][i]=N;
	for(int i=1,u,v;i<=m;i++)
		scanf("%d%d",&u,&v),g[u][v]=g[v][u]=1,G[u].push_back(v),G[v].push_back(u);
	for(int k=1;k<=n;k++)
		for(int i=1;i<=n;i++)
			for(int j=1;j<=n;j++)
				g[i][j]=min(g[i][j],g[i][k]+g[k][j]);
	for(int i=1;i<=n;i++)
	{
		double p=0;
		for(int l=0;l<n;l++)
			d[l].clear();
		for(int j=1;j<=n;j++)
			d[g[i][j]].push_back(j);
		for(int l=0;l<n;l++)
			if(d[l].size())
			{
				if(d[l].size()==1)
				{
					p+=1./n;continue;
				}
				double Max=0;
				tag1++;po.clear();
				for(int u=0;u<d[l].size();u++)
				{
					int v=d[l][u];double _p=1./d[l].size()/G[v].size();
					for(int k=0;k<G[v].size();k++)
					{
						int pos=G[v][k];
						if(vis1[pos]!=tag1)
							vis1[pos]=tag1,po.push_back(pos),P[pos]=0;
						P[pos]+=_p;
					}
				}
				for(int j=1;j<=n;j++)
				{
					double q=0;tag2++;Le.clear();
					for(int u=0;u<po.size();u++)
					{
						int v=po[u];
						if(vis2[g[v][j]]!=tag2)
							vis2[g[v][j]]=tag2,Le.push_back(g[v][j]),Maxp[g[v][j]]=0;
						Maxp[g[v][j]]=max(Maxp[g[v][j]],P[v]);
					}
					for(int u=0;u<Le.size();u++)
						q+=Maxp[Le[u]];
					Max=max(Max,q);
				}
				p+=max(1./n,Max*d[l].size()/n);
			}
		Ans=max(Ans,p);
	}
	printf("%.10lf",Ans);
	return 0;
}