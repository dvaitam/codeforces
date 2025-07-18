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

int a[110000];

int main()
{
	int i,j,k,n,l,t,m,x,y,o,b1,prev;
	vector<int> ans;
	scanf("%d",&t);
	for (l=0;l<t;l++)
	{
		scanf("%d",&n);
		for (i=0;i<n;i++)
			scanf("%d",&a[i]);
		if (a[n-1]!=0)
		{
			printf("NO\n");
			continue;
		}
		printf("YES\n");
		ans.clear();
		prev=0;
		for (i=0;i<n;i++)
			if (a[i]==0)
			{
				ans.push_back(i-prev);
				for (j=prev;j<i;j++)
				{
					ans.push_back(0);
				}
				prev=i+1;
			}
		reverse(ans.begin(),ans.end());
		for (i=0;i<n;i++)
			printf("%d ",ans[i]);
		printf("\n");
	}
}