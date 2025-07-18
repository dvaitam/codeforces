#include <cstdio>

#include <algorithm>

#include <queue>

#include <utility>

using namespace std;



const int maxn=100000+5, maxm=100000+5;

int n, m;

struct seg {

	int idx;

	int l, r;

	inline bool operator < (const seg &b) const { return l<b.l; }

};

seg s[maxm];



int main() {

	scanf("%d%d", &n, &m);

	for(int i=1; i<=m; ++i) {

		int r, c; scanf("%d%d", &r, &c);

		//r=n-r+1, c=n-c+1;

		s[i].idx=i;

		s[i].l=n-r+1, s[i].r=c;

	}

	sort(s+1, s+m+1);

	queue<int> ans; 

	priority_queue<pair<int, int>, vector<pair<int, int> >,

	 greater<pair<int, int> > > pq;

	for(int l=1, i=1; l<=n; ++l) {

		while(s[i].l==l) {

			pq.push(make_pair(s[i].r, s[i].idx));

			++i;

		}

		while(!pq.empty() && pq.top().first<l) pq.pop();

		if(!pq.empty()) {

			ans.push(pq.top().second);

			pq.pop();

		}

	}

	printf("%d\n", (int)ans.size());

	while(!ans.empty()) { printf("%d ", ans.front()); ans.pop(); }

	putchar('\n');

	return 0;

}