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

using namespace std;

pair<int,int> a[300000];
set<pair<int,int> > st;

int main()
{
	int i,j,k,n,ll,rr,x1,y1,l1,r1,b1,tmp;
	scanf("%d%d",&n,&k);
	st.clear();
	for (i=0;i<n;i++)
	{
		scanf("%d%d",&a[i].first,&a[i].second);
	}
	if (k>=n)
	{
		printf("-1\n");
		return 0;
	}
	sort(a,a+n);
	for (ll=0;ll<=k;ll++)
		for (rr=0;rr+ll<=k;rr++)
		{
			if (ll+rr==n-1)
			{
				x1=a[ll].first;
				y1=a[ll].second;
				st.insert(make_pair(x1+x1,y1+y1));
			}
			else
			{
				x1=a[ll].first+a[n-1-rr].first;
				y1=a[ll].second+a[n-1-rr].second;
				tmp=ll+rr;
				b1=1;
				l1=ll;r1=n-1-rr;
				for (i=0;l1<=r1;i++)
				{
					if (a[l1].first+a[r1].first!=x1)
					{
						if (tmp==k)
						{
							b1=0;
							break;
						}
						tmp++;
						if (a[l1].first+a[r1].first>x1)
						{
							r1--;
						}
						else l1++;
						continue;
					}
					if (a[l1].second+a[r1].second!=y1)
					{
						if (tmp==k)
						{
							b1=0;
							break;
						}
						tmp++;
						if (a[l1].second+a[r1].second>y1)
						{
							r1--;
						}
						else l1++;
						continue;
					}
					l1++;r1--;

				}
				if (b1==1)
					st.insert(make_pair(x1,y1));
			}
		}
	printf("%d\n",st.size());
	for (set<pair<int,int> >::iterator it=st.begin();it!=st.end();it++)
	{
		printf("%.1lf %.1lf\n",it->first/2.0,it->second/2.0);
	}
	return 0;
}