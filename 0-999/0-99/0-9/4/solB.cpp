#include <bits/stdc++.h>
/* Asyamov Igor
e-mail: igor9669@gmail.com*/

#include <iostream>
#include <deque>
#include <string>
#include <vector>
#include <cmath>
#include <algorithm>
#include <cstdio>
#include <map>
#include <fstream>
#include <cstdlib>
#include <queue>
#include <bitset>
#include <set>
#include <stack>
#include <utility>
#include<cassert>
using namespace std;
#define FR(i,a,b) for(int i=(a);i<(b);++i)
#define FOR(i,n) FR(i,0,n)
#define CLR(x,a) memset(x,a,sizeof(x))
#define MP make_pair
#define PB push_back
#define A first
#define B second
#define Len(a) (int)a.length()
#define Sz(a) (int)a.size()
typedef long long LL;
typedef long double LD;
typedef pair<int,int> PII;
typedef vector<int> VI;
typedef vector<VI > VVI;
#define MAXN 1001
const double Eps=1e-9;
const double Pi=2*acos(0.0);
const int inf=1000*1000*1000;

int main()
{
	//freopen("input.txt","r",stdin);freopen("output.txt","w",stdout);
	int d,time;
	scanf("%d%d",&d,&time);
	PII a[35];
	FOR(i,d)scanf("%d%d",&a[i].A,&a[i].B);
	int cur[35];
	CLR(cur,0);
	int x=0,y=0;
	FOR(i,d)cur[i]=a[i].A,x+=cur[i],y+=a[i].B;
	if(x>time )
	{
		printf("NO\n");
		return 0;
	}
	int ans[35],prev[35];
	CLR(ans,0);
	CLR(prev,0);
	bool ok=false,first=true;
	while(time>0)
	{
		ok=false;
	
		FOR(i,d)
		{
			if(time-(cur[i]-prev[i])>=0 && (cur[i]-prev[i])>0)
			{
				time-=(cur[i]-prev[i]);
				ans[i]+=(cur[i]-prev[i]);
				ok=true;
			}
			if(time==0)break;
		}
		if(!ok)break;
		if(time==0)break;
		memmove(prev,cur,sizeof(cur));
		FOR(i,d)
		{
			if(cur[i]+1<=a[i].B)cur[i]++;
		}
	}
	if(time>0)printf("NO\n");
	else
	{
		printf("YES\n");
		FOR(i,d)printf("%d ",ans[i]);
	}
	return 0;
}