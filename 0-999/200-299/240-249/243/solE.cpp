#include <cstdio>
#include <algorithm>
#include <vector>
#include <cstdlib>

using namespace std;

const int N = 550;

int n, root, bo[N], ans[N], a[N * 3];
char ss[N][N];

struct node
{
	vector<int> ch;
	int tag, col, leaf;
} d[N * 3];
int dn = 0;

bool chk(node &it)
{
	int o = 0;
	for (int j = 0; j < it.ch.size(); ++j)
	{
		a[j] = d[it.ch[j]].col;
		if (a[j] > a[o]) o = j;
	}
	for (int j = 0; j < o; ++j)
		if (a[j] > a[j + 1]) return 0;
	for (int j = o; j < it.ch.size() - 1; ++j)
		if (a[j] < a[j + 1]) return 0;
	return 1;
}
bool srt(node &it)
{
	int o = 0, z0 = 1, z1 = 1;
	for (int j = 0; j < it.ch.size(); ++j)
	{
		a[j] = d[it.ch[j]].col;
		if (j && a[j] < a[j - 1]) z0 = 0;
		if (j && a[j] > a[j - 1]) z1 = 0;
	}
	if (z0) return 1;
	if (z1)
	{
		reverse(it.ch.begin(), it.ch.end());
		return 1;
	}
	return 0;
}
void gameover()
{
	puts("NO");
	exit(0);
}
void link(int u, int v)
{
	if (v)
		d[u].ch.push_back(v);
}
void link(int u, vector<int> &v)
{
	for (int j = 0; j < v.size(); ++j)
		link(u, v[j]);
}
void merge(vector<int> &u, vector<int> v)
{
	for (int j = 0; j < v.size(); ++j)
		if (v[j]) u.push_back(v[j]);
}
void build()
{
	root = ++dn;
	for (int i = 1; i <= n; ++i)
	{
		int nd = ++dn;
		link(root, nd);
		d[nd].leaf = i;
	}
}
void paint(int u)
{
	if (d[u].leaf)
	{
		int id = d[u].leaf;
		if (bo[id])
			d[u].col = 2;
		else
			d[u].col = 0;
		return;
	}
	int cnt[3] = {0};
	for (int j = 0; j < d[u].ch.size(); ++j)
	{
		int v = d[u].ch[j];
		paint(v);
		++cnt[d[v].col];
	}
	d[u].col = 1;
	if (!cnt[0] && !cnt[1]) d[u].col = 2;
	if (!cnt[1] && !cnt[2]) d[u].col = 0;
}
int bud(vector<int> &v)
{
	if (v.size() == 0) return 0;
	if (v.size() == 1) return v[0];
	int nd = ++dn;
	link(nd, v);
	return nd;
}
vector<int> pushrig(int u)
{
	vector<int> res;
	if (!u) return res;
	if (d[u].tag && !srt(d[u])) gameover();
	int cnt[3] = {0};
	vector<int> ls[3];
	for (int j = 0; j < d[u].ch.size(); ++j)
	{
		int v = d[u].ch[j];
		int c = d[v].col;
		++cnt[c];
		ls[c].push_back(v);
	}
	if (cnt[1] > 1) gameover();

	if (!d[u].tag)
	{
		res.push_back(bud(ls[0]));
		if (cnt[1])
			merge(res, pushrig(ls[1][0]));
		res.push_back(bud(ls[2]));
	}
	else
	{
		merge(res, ls[0]);
		if (cnt[1])
			merge(res, pushrig(ls[1][0]));
		merge(res, ls[2]);
	}

	return res;
}
void setcond(int u)
{
	if (d[u].col != 1) return;
	int sc = 0, scn = 0;
	int cnt[3] = {0};
	for (int j = 0; j < d[u].ch.size(); ++j)
	{
		int v = d[u].ch[j];
		int c = d[v].col;
		++cnt[c];
		if (c) sc = v, ++scn;
	}
	if (scn == 1) return setcond(sc);

	if (cnt[1] > 2) gameover();
	vector<int> ls[3];
	for (int j = 0; j < d[u].ch.size(); ++j)
	{
		int v = d[u].ch[j];
		ls[d[v].col].push_back(v);
	}
	ls[1].push_back(0);
	ls[1].push_back(0);

	if (d[u].tag)
	{
		if (!chk(d[u])) gameover();
	}
	else
	{
		d[u].ch.clear();
		int a = bud(ls[2]);
		d[a].col = 2;
		if (ls[1].size() > 2)
		{
			int b = ++dn;
			link(b, ls[1][0]);
			link(b, a);
			link(b, ls[1][1]);
			d[b].col = 1;
			d[b].tag = 1;
	
			link(u, b);
			link(u, ls[0]);
			return setcond(b);
		}
		else
		{
			link(u, a);
			link(u, ls[0]);
			return;
		}
	}

	// now u's tag must be '1', which means fixed
	vector<int> newch;
	int flag = 0;
	for (int j = 0; j < d[u].ch.size(); ++j)
	{
		int v = d[u].ch[j];
		int c = d[v].col;
		if (c == 1)
		{
			vector<int> res = pushrig(v);
			if (flag)
				reverse(res.begin(), res.end());
			merge(newch, res);
			flag = 1;
		}
		else
		{
			flag |= (c == 2);
			newch.push_back(v);
		}
	}
	d[u].ch = newch;
}
void workans(int u)
{
	if (d[u].leaf)
	{
		ans[++ans[0]] = d[u].leaf;
		return;
	}
	for (int j = 0; j < d[u].ch.size(); ++j)
		workans(d[u].ch[j]);
}
void showtree(int u)
{
	printf("%d  tag=%d  ", u, d[u].tag);
	putchar('(');
	for (int j = 0; j < d[u].ch.size(); ++j)
		printf("%d,", d[u].ch[j]);
	printf(")\n");

	for (int j = 0; j < d[u].ch.size(); ++j)
		showtree(d[u].ch[j]);
}
int main()
{
	scanf("%d", &n);
	for (int i = 1; i <= n; ++i)
		scanf("%s", ss[i] + 1);

	build();
	for (int i = 1; i <= n; ++i)
	{
		for (int j = 1; j <= n; ++j)
			bo[j] = (ss[i][j] == '1');
		paint(root);
		setcond(root);

		//showtree(root),puts(""),puts("");
	}

	workans(root);
	puts("YES");
	for (int i = 1; i <= n; ++i)
	{
		for (int j = 1; j <= n; ++j)
			putchar(ss[i][ans[j]]);
		puts("");
	}
}