#include <iostream>
#include <fstream>
#include <vector>

using namespace std;

const int inf = 100000000;

int a[2000] = { };
int b[2000] = { };
int c[2000] = { };

char ans[2000][2000] = { } ;
vector <int> q[2000] = { } ;
int q1[2000] = {};
int n;

int all;

bool END = false;

void dfs(int depth)
{
	if (b[depth] != 0)
	{
		b[depth] --;
		all --;
		if (all == 0) END = true;
		for (int i = 0; i < depth; i ++) 
		{
//			cout << (char) (c[i] + '0');
			ans[q[depth][q1[depth]]][i] = c[i] + '0';
		}
//		cout << endl;
		q1[depth] ++;
		return;
	}
	for (int i = 0; i <= 1; i ++)
	{
		c[depth] = i;
		dfs(depth + 1);
		if (END) return;
	}

}

int main()
{
	cin >> n;
	for (int i = 0; i < n; i ++)
	{
		cin >> a[i];
		b[a[i]] ++;
		q[a[i]].push_back(i);
	}
	int cur = 2;
	for (int i = 1; i <= 1000; i ++)
	{
		if (b[i] > cur)
		{
			cout << "NO" << endl;
			return 0;
		}
		cur -= b[i];
		cur *= 2;
		if (cur > inf)
			break;
	}
	all = n;
	dfs(0);
	printf("YES\n");
	for (int i = 0; i < n; i ++)
		printf("%s\n",ans[i]);
	return 0;
}