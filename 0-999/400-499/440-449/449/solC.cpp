#include <cstdio>
#include <utility>

using namespace std;

const int MAXN = 100000;
pair <int, int> grp [MAXN >> 1];
int n, m, p, i, j;
bool g [MAXN + 1] = {false};

int main () {
	m = p = 0;
	scanf ("%d", & n);
	for (i = 3; i <= n >> 1; i += 2)
		if (! g [i]) {
			p = i;
			for (j = 3 * i; j <= n; j += i)
				if (! g [j]) {
					if (p) {
						grp [m ++] = make_pair (p, j);
						g [p] = g [j] = true;
						p = 0;
					}
					else
						p = j;
				}
			if (p) {
				grp [m ++] = make_pair (p, i << 1);
				g [p] = g [i << 1] = true;
			}
		}
	p = 0;
	for (i = 2; i <= n; i += 2)
		if (! g [i]) {
			if (p) {
				grp [m ++] = make_pair (p, i);
				g [p] = g [i] = true;
				p = 0;
			}
			else
				p = i;
		}
	printf ("%d\n", m);
	for (i = 0; i < m; i ++)
		printf ("%d %d\n", grp [i].first, grp [i].second);
	return 0;
}