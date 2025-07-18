#include <cstdio>

#include <algorithm>

#include <cstring>

using namespace std;



struct Cube {

	int id;

	int color, size;

	inline bool operator < (const Cube &b) const {

		if(color!=b.color) return color<b.color;

		return size>b.size;

	}

};

const int N=100000+5;

int n;

Cube cube[N];

int start[N];



long long mx[N];

int mxidx[N];

inline void update_mx(int len, long long v, int idx) {

	if(mx[len]<v) {

		mx[len]=v;

		mxidx[len]=idx;

	}

}

long long ans=0;

int ansidx1, ansidx2;

inline void update_ans(long long v, int idx1, int idx2) {

	if(ans<v) {

		ans=v;

		ansidx1=idx1, ansidx2=idx2;

	}

}



int main() {

	scanf("%d", &n);

	for(int i=1; i<=n; ++i) {

		cube[i].id=i;

		scanf("%d%d", &cube[i].color, &cube[i].size);

	}

	sort(cube+1, cube+n+1);

	memset(mx, 0xc0, sizeof(mx));

	for(int r=1; r<=n; ) {

		int l=r;

		while(r<=n && cube[l].color==cube[r].color) ++r;

		long long sum;

		sum=0;

		for(int i=l; i<r; ++i) {

			start[i]=l;

			sum+=cube[i].size;

			int len=i-l+1;

			update_ans(mx[len-1]+sum, mxidx[len-1], i);

			update_ans(mx[len]+sum, mxidx[len], i);

			update_ans(mx[len+1]+sum, mxidx[len+1], i);

		}

		sum=0;

		for(int i=l; i<r; ++i) {

			sum+=cube[i].size;

			int len=i-l+1;

			update_mx(len, sum, i);

		}

	}

	printf("%I64d\n", ans);

	int len1=ansidx1-start[ansidx1]+1, len2=ansidx2-start[ansidx2]+1;

	printf("%d\n", len1+len2);

	if(len1<len2) swap(len1, len2), swap(ansidx1, ansidx2);

	while(len1 || len2) {

		if(len1>len2) {

			printf("%d ", cube[ansidx1--].id);

			--len1;

		} else {

			printf("%d ", cube[ansidx2--].id);

			--len2;

		}

	}

	putchar('\n');

	return 0;

}