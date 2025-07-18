#include <stdio.h>
#include <algorithm>

const int N=200010;

int n, a[N], s;

long long ans;

inline void read(int &x) {
	char c=getchar(); x=0; while (c<'0'||c>'9') c=getchar();
	while (c>='0'&&c<='9') x=x*10+c-'0', c=getchar(); return;
}

inline int abs(int x) {
	return x<0?-x:x;
}

int main() {
	read(n); read(s);
	register int i;
	int pos, l, r;
	for (i=1; i<=n; i++) read(a[i]);
	std::sort(a+1,a+n+1);
	pos=std::upper_bound(a+1,a+n+1,s)-a-1;
	l=1+n>>1; r=pos;
	if (l>r) {
		std::swap(l,r);
		l++;
	}
	for (i=l; i<=r; i++) {
		ans+=abs(a[i]-s);
	}
	printf("%lld", ans);
	return 0;
}