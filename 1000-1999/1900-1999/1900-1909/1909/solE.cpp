#include<bits/stdc++.h>
using namespace std;
vector<int> imp[20];
int adj[20][20];
int main()
{
	for(int n = 5; n <= 19; n++)
	{	
		for(int i = 1; i < (1 << n); i++)
		{
			int a[20];
			for(int j = 1; j <= n; j++) a[j] = 0;
			for(int d = 0; d < n; d++)
			{
				if(i & (1 << d))
				{
					for(int j = d + 1; j <= 19; j += d + 1)
					{
						a[j] ^= 1;
					 } 
				}			
			}
			int z = 0;
			for(int j = 1; j <= n; j++)
			{
				z += a[j];
			}
			if(z <= n / 5) 
			{
				imp[n].push_back(i);
			}
	
		}
	}
	int t;
	scanf("%d", &t);
	while(t--)
	{
		int n, m;
		scanf("%d %d", &n, &m);
		for(int x = 1; x <= 19; x++) for(int y = 1; y <= 19; y++) adj[x][y] = 0;
		for(int i = 0; i < m; i++)
		{
			int x, y;
			scanf("%d %d", &x, &y);
			if(x < 20 && y < 20) adj[x][y] = 1;
		}
		if(n < 5)
		{
			printf("-1\n");
			continue;
		}
		if(n > 19)
		{
			printf("%d\n", n);
			for(int i = 0; i < n; i++) printf("%d ", i + 1);
			printf("\n");
			continue;
		}
		vector<pair<int, int> > ll;
		for(int x = 1; x <= 19; x++) for(int y = 1; y <= 19; y++) 
		{
			if(adj[x][y]) ll.push_back(make_pair(x, y));
		}
		int ans = -1;
		for(int i = 0; i < (int)imp[n].size(); i++)
		{
			int bad = 0;
			int j = imp[n][i];
			for(int k = 0; k < (int)ll.size(); k++)
			{
				int x = ll[k].first;
				int y = ll[k].second;
				int xx = (1 << (x - 1));
				int yy = (1 << (y - 1));
				if(adj[x][y])
				{
					//printf("not %d %d \n", x, y);
					if(xx & j) 
					{
						//printf("Here j = %d, yy = %d, %d \n", j, yy, j & yy);
						if((yy & j) == 0)
						{
							//printf("bad\n");
							bad = 1;
							break;
						}
					}
				}	
			}
			if(!bad)
			{
				ans = j;
				break;
			}
		}
		if(ans == -1)
		{
			printf("-1\n");
			continue;
		}
		printf("%d\n", __builtin_popcount(ans));
		for(int i = 0; i < n; i++)
		{
			if(ans & (1 << i))
			{
				printf("%d ", i + 1);
			}
		}
		printf("\n");
	}
	return 0;
}