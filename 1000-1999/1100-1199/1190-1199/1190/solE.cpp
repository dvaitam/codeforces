#include <bits/stdc++.h>
using namespace std;

#define ll long long
#define ar array
#define ld long double

const int mxN=1e5;
const ld PI=acos(-1);
int n, m, st[mxN][17];
ar<int, 2> a[mxN];
ar<ld, 2> b[mxN];
ld dist[mxN], th[mxN];
vector<ar<ld, 2>> c;

bool chk(ld r) {
	//cout << "chk " << r << endl;
	//create ranges
	for(int i=0; i<n; ++i) {
		ld ac=acos(r/dist[i]);
		b[i]={th[i]-ac, th[i]+ac};
		if(b[i][0]<0) {
			b[i][0]+=2*PI;
			b[i][1]+=2*PI;
		}
		//cout << b[i][0] << " " << b[i][1] << endl;
	}
	//sort intervals, remove unnecessary ones
	sort(b, b+n);
	c.clear();
	for(int i=0; i<n; ++i) {
		//while last one covers current
		while(c.size()&&c.back()[1]>=b[i][1])
			c.pop_back();
		//while(c.size()&&b[i][1]>=2*PI&&c.front()[1]<=b[i][1]-2*PI)
			//c.pop();
		if(!c.size()||c[0][1]>b[i][1]-2*PI)
			c.push_back(b[i]);
	}
	//while(qu.size()) {
		//c.push_back(qu.front());
		//qu.pop_front();
	//}
	//vector<int> n;
	int aa=c.size();
	for(int i=0; i<aa; ++i) {
		c.push_back(c[i]);
		c.back()[0]+=2*PI;
		c.back()[1]+=2*PI;
	}
	//set up base case for sparse table
	for(int i=0, j=0; i<aa; ++i) {
		while(j<i+aa&&c[j][0]<=c[i][1])
			++j;
		st[i][0]=j-i;
	}
	//set up powers of 2 for sparse table
	for(int k=1; k<17; ++k) {
		for(int i=0; i<aa; ++i) {
			st[i][k]=st[i][k-1]+st[(i+st[i][k-1])%aa][k-1];
		}
	}
	//find if it is actually possible
	for(int s=0; s<aa; ++s) {
		int used=0, covered=0;
		int i=s;
		for(int k=16; ~k; --k) {
			if(used+(1<<k)<=m) {
				covered+=st[i][k];
				i=(i+st[i][k])%aa;
				used+=1<<k;
			}
		}
		if(covered>=aa)
			return 1;
	}
	return 0;
}

int main() {
	ios::sync_with_stdio(0);
	cin.tie(0);
	
	cout << fixed << setprecision(9);
	cin >> n >> m;
	for(int i=0; i<n; ++i) {
		cin >> a[i][0] >> a[i][1];
		if(!a[i][0]&&!a[i][1]) {
			cout << 0;
			return 0;
		}
	}
	sort(a, a+n);
	n=unique(a, a+n)-a;
	ld lb=0, rb=1e6;
	for(int i=0; i<n; ++i) {
		dist[i]=hypot(a[i][0], a[i][1]);
		rb=min(dist[i], rb);
		th[i]=atan2(a[i][1], a[i][0]);
	}
	//cout << chk(1.000001);
	//return 0;
	while(rb-lb>1e-7) {
		ld mb=(lb+rb)/2;
		if(chk(mb))
			lb=mb;
		else
			rb=mb;
	}
	cout << (lb+rb)/2;
}