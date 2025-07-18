#include <bits/stdc++.h>

using namespace std;
typedef long long int lli;
typedef pair<int,int> pii;

const int MAXN = 200010;

int n;

int v[MAXN];

bool findPath[MAXN];

vector< pii > ans;

int main()
{
	scanf("%d",&n);

	for(int i = 1 ; i < n ; i++)
		scanf("%d",&v[i]);

	printf("%d\n",v[1]);

	findPath[ v[1] ] = true;

	int p = 2;
	int last = n;

	while( true )
	{
		while( last > 0 && findPath[last] ) last--;
		findPath[last] = true;

		if( last == 0 ) break;

		while( p < n && !findPath[ v[p] ] )
		{
			ans.push_back( { v[p - 1] , v[p] } );
			findPath[ v[p] ] = true;
			p++;
		}

		ans.push_back( { last , v[p - 1] } );
		p++;
	}

	for(int i = 0 ; i < n - 1 ; i++)
		printf("%d %d\n",ans[i].first,ans[i].second);
}