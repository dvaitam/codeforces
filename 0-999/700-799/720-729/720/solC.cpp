#include <cstdio>

#include <cstring>

#include <algorithm>

#include <vector>

#include <assert.h>

using namespace std;

inline int read()

{

	int s = 0; char c; while((c=getchar())<'0'||c>'9');

	do{s=s*10+c-'0';}while((c=getchar())>='0'&&c<='9');

	return s;

}

const int N = 100010;

int n,m,t,rev,rt;

//char mp[N][N];

vector<char> mp[N];

bool work()

{

	int i,j;

	if(t>(n-1)*(m-1)*4) return 0;

	//if(n<=m) rev = 0; else rev = 1, swap(n,m);

	for(i=1;i<=n;i++) mp[i].resize(m+1);

	for(i=1;i<=n;i++) for(j=1;j<=m;j++) mp[i][j] = '.';

	mp[1][1] = mp[2][1] = '*';

	for(i=2;i<=n;i++)

	{

		if(i>2) t--, mp[i][1] = '*'; j = 2;

		if(t<4) break;

		for(j=2;j<=m&&t>=4;j++)

			mp[i-1][j-1] = mp[i-1][j] = mp[i][j-1] = mp[i][j] = '*', t -= (i>2&&j==m)?3:4;

		if(t<4) break;

	}

	if(i==2&&j<=m)

	{

		switch(t)

		{

			case 1: mp[1][j] = '*'; break;

			case 2: if(j==2) mp[2][2] = mp[3][1] = '*'; else mp[1][j] = mp[3][1] = '*'; break;

			case 3: if(j==2) mp[3][3] = mp[2][3] = '*'; mp[1][j] = mp[3][j-1] = mp[3][j] = '*'; break;

		}

	}

	else

	{

		switch(t)

		{

			case 1: if(j<m) mp[i][m] = '*'; else if(i<n) mp[i+1][1] = '*'; else return 0; break;

			case 2: 

				if(j<m-1) mp[i][m-1] = '*'; 

				else if(i<n&&j>3) mp[i+1][2] = '*'; 

				else if(i<n&&j<=m) mp[i+1][j] = mp[i+1][j-1] = '*'; 

				else return 0; 

				break;

			case 3:

				if(j<m-2) mp[i][m] = mp[i][m-2] = '*';

				else if(j==m) mp[i][j] = '*';

				else if(i<n)

				{

					if(j<m)

					{

						mp[i][m] = mp[i+1][m-1] = mp[i+1][j] = mp[i+1][j-1] = '*';

						if(j==2) mp[i+1][m] = '*';

					}

					else if(m>=4)

					{

						mp[i+1][1] = mp[i+1][m-1] = '*';

						if(j==m) mp[i+1][m] = '*';

					}

					else return 0;

				}

				else return 0;

				break;

		}

	}

	/*if(rev)

	{

		for(i=1;i<=n;i++) for(j=1;j<i;j++) swap(mp[i][j],mp[j][i]);

		swap(n,m);

	}*/

	return 1;

}

int main()

{

	#ifndef ONLINE_JUDGE

	freopen("in.txt","r",stdin);

	#endif

	int i,j,y;

	for(int T=read();T--;)

	{

		n = read(); m = read(); y = t = read();

		if(!work())

		{

			swap(n,m); t = y;

			if(!work()) puts("-1");

			else

			{

				for(i=1;i<=m;i++)

				{

					for(j=1;j<=n;j++) putchar(mp[j][i]);

					puts("");

				}

			}

		}

		else 

			for(i=1;i<=n;i++)

			{

				for(j=1;j<=m;j++) putchar(mp[i][j]);

				puts("");

			}

		if(T) puts("");

	}

	return 0;

}