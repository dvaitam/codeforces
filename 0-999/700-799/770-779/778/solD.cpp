#include <stdio.h>
#include <algorithm>
#include <vector>
#include <queue>
#include <map>

using namespace std;

const int MAX_N = 55;

typedef pair<int, int> pi;

char Mp[2][MAX_N][MAX_N];
char now[MAX_N][MAX_N];
int N, M;
vector<pi> Ans[2]; int nowK;
bool changeOnce(int n, int m, char k)
{
//	printf("%d %d (%c) -> %c\n", n, m, now[n][m], k);
	if(now[n][m] == k) return true;
	else if(k == 'L' && now[n][m] == 'U')
	{
		changeOnce(n, m+1, 'U');
		now[n][m] = now[n+1][m] = 'L';
		now[n][m+1] = now[n+1][m+1] = 'R';
		Ans[nowK].push_back(make_pair(n, m));
		return true;
	}
	else if(k == 'U' && now[n][m] == 'L')
	{
		changeOnce(n+1, m, 'L');
		now[n][m] = now[n][m+1] = 'U';
		now[n+1][m] = now[n+1][m+1] = 'D';
		Ans[nowK].push_back(make_pair(n, m));
		return true;
	}
	else
	{
		printf("once ?? %d %d : %c but %c\n", n, m, k, now[n][m]);
		return false;
	}
}
bool change(int n, int m)
{
	if(n > N || (N % 2 == 1 && n == N))
	{
		return true;
	}
	if(m > M) return change(n+2, 1);
	if(now[n][m] == 'U') return change(n, m+1);
	else if(now[n][m] == 'L')
	{
		changeOnce(n+1, m, 'L');
		Ans[nowK].push_back(make_pair(n, m));
		now[n][m] = now[n][m+1] = 'U';
		now[n+1][m] = now[n+1][m+1] = 'D';
		/*
		printf("{{{\n");
		for(int i=1; i<=N; i++) printf("%s\n", now[i]+1);
		printf("}}}\n");
		*/
		return change(n, m+2);
	}
	else
	{
		printf("change ?? %d %d : %d\n", n, m, now[n][m]);
		return false;
	}
}
void calc(int k)
{
	for(int i=1; i<=N; i++) for(int j=1; j<=M; j++) now[i][j] = Mp[k][i][j];
	nowK = k;
	change(1, 1);
}
int main()
{
	scanf("%d%d", &N, &M); scanf(" ");
	for(int k=0; k<2; k++) for(int i=1; i<=N; i++) scanf("%s", Mp[k][i]+1);
	for(int k=0; k<2; k++) calc(k);
	printf("%d\n", Ans[0].size() + Ans[1].size());
	reverse(Ans[1].begin(), Ans[1].end());
	for(pi &x : Ans[0]) printf("%d %d\n", x.first, x.second);
	for(pi &x : Ans[1]) printf("%d %d\n", x.first, x.second);
	return 0;
}