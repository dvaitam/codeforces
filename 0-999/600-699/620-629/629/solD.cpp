#include <stdio.h>
#include <algorithm>

using namespace std;

const double pi = 3.14159265358979323846;

struct cake
{
	int idx;
	double r2h;

	bool operator <(const cake &rhs) const
	{
		return r2h > rhs.r2h;
	}

}in[100005];

int n;
double tree[100005], v[100005];

void update(int idx, double value)
{
	for (;idx <= n;idx += idx & -idx)
		if (tree[idx] < value)
			tree[idx] = value;
}

double query(int idx)
{
	double ret = 0;
	for (;idx;idx -= idx & -idx)
		if (ret < tree[idx])
			ret = tree[idx];

	return ret;
}

int main()
{
	scanf("%d", &n);

	for (int i = 1;i <= n;++i)
	{
		int r, h; scanf("%d %d", &r, &h);

		in[i].idx = n - i + 1;
		in[i].r2h = r * r;
		in[i].r2h *= h;
	}

	sort(in + 1, in + n + 1);

	double ans = 0;
	int last = 1;
	for (int i = 1;i <= n;++i)
	{
		v[i] = in[i].r2h + query(in[i].idx - 1);
		if (i < n && in[i].r2h == in[i + 1].r2h);
		else for (;last <= i;++last) update(in[last].idx, v[last]);

		ans = max(ans, v[i]);
	}
	printf("%.9lf", ans * pi);

	return 0;
}