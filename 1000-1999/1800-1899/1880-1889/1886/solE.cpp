#include <vector>
#include <list>
#include <map>
#include <set>
#include <deque>
#include <vector>
#include <list>
#include <map>
#include <set>
#include <deque>
#include <queue>
#include <stack>
#include <bitset>
#include <algorithm>
#include <functional>
#include <numeric>
#include <utility>
#include <sstream>
#include <iostream>
#include <iomanip>
#include <cstdio>
#include <cmath>
#include <cstdlib>
#include <cctype>
#include <string>
#include <cstring>
#include <ctime>
#include <random>
#include <chrono>
 
using namespace std;
 
#define _int64 long long
#define mo 998244353

int a[210000];
int b[100];
int d[21][210000];
int d2[1<<20];
int p[1<<20][2];
int n,m;

vector<int> sol[30];

int doit(int bit)
{
	int i,tmp,ans;
	if (bit==0) return n;
	if (d2[bit]!=-1) return d2[bit];
	ans=-2;
	for (i=0;i<m;i++)
	{
		if (((1<<i)&bit)==0) continue;
		tmp=doit(bit^(1<<i));
		if (tmp==-2) continue;
		if (tmp==0) continue;
		if (d[i][tmp-1]==-1) continue;
		if (d[i][tmp-1]>ans)
		{
			//cerr<<"bit,i,tmp,d[i][tmp-1]:"<<bit<<" "<<i<<" "<<tmp<<" "<<d[i][tmp-1]<<endl;
			ans=max(ans,d[i][tmp-1]);
			p[bit][0]=tmp;
			p[bit][1]=i;
		}
	}
	d2[bit]=ans;
	return ans;
}

int main()
{
	int i,j,k,l,t,x,y,o,b1,tmp,ans,v,tmpv;
	vector<pair<int,int> > a;
	a.clear();
	scanf("%d%d",&n,&m);
	for (i=0;i<n;i++)
	{
		scanf("%d",&x);
		a.push_back(make_pair(x,i));
	}
	sort(a.begin(),a.end());
	for (i=0;i<m;i++)
		scanf("%d",&b[i]);
	for (i=0;i<m;i++)
	{
		memset(d[i],-1,sizeof(d[i]));
		for (j=0;j<n;j++)
		{
			tmp=j+(b[i]+a[j].first-1)/a[j].first-1;
			if (tmp>=n) continue;
			d[i][tmp]=max(d[i][tmp],j);
		}
		for (j=0;j+1<n;j++)
			d[i][j+1]=max(d[i][j+1],d[i][j]);
		/*
		cerr<<"d[i]:"<<i<<" ";
		for (j=0;j<n;j++)
			cerr<<d[i][j]<<" ";
		cerr<<endl;
		*/
	}
	memset(d2,-1,sizeof(d2));
	ans=doit((1<<m)-1);
	if (ans==-2)
	{
		printf("NO\n");
		return 0;
	}
	printf("YES\n");
	tmp=(1<<m)-1;
	//v=n;
	while (tmp>0)
	{
		x=p[tmp][1];
		sol[x].clear();
		v=p[tmp][0];
		tmpv=d[x][v-1];
		//cerr<<"x:"<<tmp<<" "<<x<<" "<<v<<" "<<tmpv<<" "<<endl;
		for (i=tmpv;i<v;i++)
			sol[x].push_back(a[i].second);
		tmp^=(1<<x);
		//v=tmpv;
	}
	for (i=0;i<m;i++)
	{
		printf("%d",(int)sol[i].size());
		for (j=0;j<sol[i].size();j++)
			printf(" %d",sol[i][j]+1);
		printf("\n");
	}
}