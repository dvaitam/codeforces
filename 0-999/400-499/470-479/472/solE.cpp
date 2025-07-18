#include <vector>

#include <map>

#include <numeric>

#include <iostream>

#include <cmath>

#include <cstdio>

#include <cassert>

#include <cstring>

#include <algorithm>

#include <ctime>

#include <stack>

#include <queue>

using namespace std;

typedef long long ll;

#define REP(i,n) for(int i=0;i<(n);i++)

const double EPS=1e-7;



const int MAXN=32;



int n,m;

int s[MAXN][MAXN];

int t[MAXN][MAXN];

int fix[MAXN][MAXN];



vector<pair<int,int> > res;



int sign(int x)

{

	if(x<0) return -1;

	if(x>0) return 1;

	return 0;

}



vector<int> solve1d(const vector<int> &s, const vector<int> &t)

{

	int n=s.size();

	vector<int> tmp1(s), tmp2(t);

	sort(tmp1.begin(), tmp1.end());

	sort(tmp2.begin(), tmp2.end());

	if(tmp1!=tmp2)

		return vector<int>(1,-1);

	for(int from=0;from<n;from++)

		for(int to=0;to<n;to++)

			if(from!=to)

			{

				tmp1=s;

				for(int i=from;i!=to;i+=sign(to-from))

					swap(tmp1[i], tmp1[i+sign(to-from)]);

				if(tmp1==t)

				{

					vector<int> res;

					for(int i=from;i!=to;i+=sign(to-from))

						res.push_back(i);

					res.push_back(to);

					return res;

				}

			}

	return vector<int>(1,-1);

}



int curx, cury;



void go(int nx, int ny)

{

	assert(abs(curx-nx)<=1 && abs(cury-ny)<=1);

	swap(s[curx][cury], s[nx][ny]);

	res.push_back(make_pair(nx, ny));

	curx=nx;

	cury=ny;

}



int iter;

int used[MAXN][MAXN];

int prvX[MAXN][MAXN];

int prvY[MAXN][MAXN];



void goFar(int tx, int ty)

{

	iter++;

	queue<int> q;

	used[tx][ty]=true;

	q.push(tx);

	q.push(ty);

	while(used[curx][cury]!=iter)

	{

		int x=q.front(); q.pop();

		int y=q.front(); q.pop();

		for(int dx=-1;dx<=1;dx++)

			for(int dy=-1;dy<=1;dy++)

			{

				int nx=x+dx;

				int ny=y+dy;

				if(!fix[nx][ny] && used[nx][ny]!=iter)

				{

					used[nx][ny]=iter;

					prvX[nx][ny]=x;

					prvY[nx][ny]=y;

					q.push(nx);

					q.push(ny);

				}

			}

	}

	while(curx!=tx || cury!=ty)

		go(prvX[curx][cury], prvY[curx][cury]);

}



int main()

{

	scanf("%d%d",&n,&m);

	if(n==1 || m==1)

	{

		vector<int> s,t;

		REP(i,n*m)

		{

			int tmp;

			scanf("%d",&tmp);

			s.push_back(tmp);

		}

		REP(i,n*m)

		{

			int tmp;

			scanf("%d",&tmp);

			t.push_back(tmp);

		}

		vector<int> cr=solve1d(s,t);

		if(cr[0]==-1)

			res.push_back(make_pair(-1,-1));

		else if(n==1)

		{

			for(int i:cr)

				res.push_back(make_pair(1,i+1));

		}

		else

		{

			for(int i:cr)

				res.push_back(make_pair(i+1,1));

		}

	}

	else

	{

		memset(s,-1,sizeof(s));

		REP(i,n) REP(j,m)

			scanf("%d",&s[i+1][j+1]);

		REP(i,n) REP(j,m)

			scanf("%d",&t[i+1][j+1]);



		REP(i,m+2)

			fix[0][i]=fix[n+1][i]=true;

		REP(i,n+2)

			fix[i][0]=fix[i][m+1]=true;

		int cnt[901]={};

		for(int i=1;i<=n;i++)

			for(int j=1;j<=m;j++)

			{

				cnt[s[i][j]]++;

				cnt[t[i][j]]--;

			}

		bool ok=true;

		for(int i=1;i<=900;i++)

			ok&=cnt[i]==0;



		if(!ok)

			res.push_back(make_pair(-1,-1));

		else

		{

			curx=-1;

			cury=-1;

			for(int i=1;i<=n;i++)

				for(int j=1;j<=m;j++)

					if(s[i][j]==t[n][m])

					{

						curx=i;

						cury=j;

					}

			res.push_back(make_pair(curx, cury));

			for(int i=1;i<=n-2;i++)

				for(int j=1;j<=m;j++)

				{

					int tx=-1, ty=-1;

					for(int k=1;k<=n;k++)

						for(int l=1;l<=m;l++)

							if(!fix[k][l] && (k!=curx || l !=cury) && (tx==-1 && ty==-1) && s[k][l]==t[i][j])

							{

								tx=k;

								ty=l;

							}

					while(ty!=j)

					{

						fix[tx][ty]=true;

						goFar(tx, ty+sign(j-ty));

						fix[tx][ty]=false;

						go(tx,ty);

						ty+=sign(j-ty);

					}

					while(tx!=i)

					{

						fix[tx][ty]=true;

						goFar(tx+sign(i-tx), ty);

						fix[tx][ty]=false;

						go(tx, ty);

						tx+=sign(i-tx);

					}



					fix[i][j]=true;

				}

			for(int j=1;j<=m;j++)

				for(int i=n-1;i<=n;i++)

				{

					if(j==m && i==n) continue;

					int tx=-1, ty=-1;

					for(int k=1;k<=n;k++)

						for(int l=1;l<=m;l++)

							if(!fix[k][l] && (k!=curx || l !=cury) && (tx==-1 && ty==-1) && s[k][l]==t[i][j])

							{

								tx=k;

								ty=l;

							}

					while(tx!=i)

					{

						fix[tx][ty]=true;

						goFar(tx+sign(i-tx), ty);

						fix[tx][ty]=false;

						go(tx, ty);

						tx+=sign(i-tx);

					}

					while(ty!=j)

					{

						fix[tx][ty]=true;

						goFar(tx, ty+sign(j-ty));

						fix[tx][ty]=false;

						go(tx,ty);

						ty+=sign(j-ty);

					}

					fix[i][j]=true;

				}

		}

	}

	if(res[0]==make_pair(-1,-1))

		puts("-1");

	else

	{

		printf("%u\n",res.size()-1);

		for(int i=0;i<res.size();i++)

			printf("%d %d\n",res[i].first, res[i].second);

	}

	return 0;

}