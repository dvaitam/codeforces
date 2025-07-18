#include <bits/stdc++.h>
using namespace std;

typedef long long ll;

const int mxN=1e6;
int n, k, a[mxN], l[mxN];

struct dsu {
	int p[mxN], ld[mxN], lb[mxN];
	void init() {
		iota(p, p+n, 0);
		iota(lb, lb+n, 0);
	}
	int find(int x) {
		return x==p[x]?x:(p[x]=find(p[x]));
	}
	void join(int x, int y) {
		x=find(x), y=find(y); //necessary?
		p[x]=y;
		lb[y]=lb[x];
		ld[y]+=ld[x];
	}
	void upd(int i, int j) {
		i=find(i);
		++ld[i];
		if(j<n)
			--ld[j];
		while(ld[i]>=0&&lb[i])
			join(lb[i]-1, i);
	}
	void upd2(int i) {
		if(i<k)
			return;
		int j=find(i-k+1);
		while(lb[j]>i-k)
			join(lb[j]-1, j);
	}
	int qry() {
		return ld[find(0)];
	}
} d;

int main() {
	ios::sync_with_stdio(0); cin.tie(0);
	cin >> n >> k;
	for(int i=0; i<n; ++i)
		cin >> a[i];
	d.init();
	for(int i=0; i<n; ++i) {
		d.upd2(i);
		l[i]=i-1;
		while(l[i]>=0&&a[l[i]]<a[i]) {
			l[i]=l[l[i]];
		}
		d.upd(l[i]+1, i+1);
		if(i>=k-1)
		cout << d.qry() << " ";
	}
	return 0;
}