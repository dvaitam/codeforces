#include<stdlib.h>

#include<time.h>

#include<cmath>

#include<cstring>

#include<cstdio>

#include<queue>

#include<vector>

#include<algorithm>

using namespace std;

typedef long long LL;

typedef unsigned long long UL;

typedef vector<int> vi;

typedef pair<int, int> pii;

#define sz(x) ((int)(x.size()))

#define sqr(x) ((x)*(x))

#define pb push_back

#define mp make_pair

#define fi first

#define se second

const UL SEED = 3e2 + 57;

const UL BASE = 1e4 + 317;

const LL INF = 2e7 + 7;

const int N = 1e5 + 7;

const int MOD = 1e9 + 7;

const double Pi = acos(-1.);

const double EPS = 1e-7;

int n, m;

char str[100];

struct Node {

	char host[30], path[30];

	bool operator<(const Node &p) const {

		return strcmp(host, p.host) < 0

				|| (strcmp(host, p.host) == 0 && strcmp(path, p.path) < 0);

	}

	bool operator==(const Node &p) const {

		return strcmp(host, p.host) == 0 && strcmp(path, p.path) == 0;

	}

} a[N];

struct HASH {

	int id;

	UL val;

	bool operator<(const HASH &p) const {

		return val < p.val;

	}

} b[N];

void read(int x) {

	scanf(" %s", str);

	// host

	int i = 7, j = 0;

	while (('a' <= str[i] && str[i] <= 'z') || str[i] == '.') {

		a[x].host[j++] = str[i++];

	}

	a[x].host[j] = '\0';

	// path

	j = 0;

	while (('a' <= str[i] && str[i] <= 'z') || str[i] == '.' || str[i] == '/') {

		a[x].path[j++] = str[i++];

	}

	if (!j)

		a[x].path[j++] = '$';

	a[x].path[j] = '\0';

}

UL hsv(char s[]) {

	UL val = 0;

	for (int i = 0; s[i]; ++i)

		val = val * SEED + s[i];

	return val;

}

int main() {

	scanf("%d", &n);

	for (int i = 0; i < n; ++i) {

		read(i);

//		printf("%s, %s\n", a[i].host, a[i].path);

	}

	sort(a, a + n);

	n = unique(a, a + n) - a;

	for (int i = 0; i < n;) {

		int j = i + 1;

		while (j < n && strcmp(a[i].host, a[j].host) == 0)

			++j;

		b[m].id = i;

		for (; i < j; ++i)

			b[m].val = b[m].val * BASE + hsv(a[i].path);

		++m;

	}

	sort(b, b + m);

	int k = 0;

	for (int i = 0; i < m;) {

		int j = i + 1;

		while (j < m && b[i].val == b[j].val)

			++j;

		if (j - i > 1)

			++k;

		i = j;

	}

	printf("%d\n", k);

	for (int i = 0; i < m;) {

		int j = i + 1;

		while (j < m && b[i].val == b[j].val)

			++j;

		if (j - i > 1) {

			printf("http://%s", a[b[i].id].host);

			for (k = i + 1; k < j; ++k)

				printf(" http://%s", a[b[k].id].host);

			puts("");

		}

		i = j;

	}

	return 0;

}