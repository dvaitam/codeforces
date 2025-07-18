#include <bits/stdc++.h>
using namespace std;

int T, n, m, sum;
int vis[2005];
struct node {
	int num;
	int wei;
} a[1005];

void csh() {
	vis[1] = 0; vis[2] = 2; vis[3] = 3;
	for(int i = 4;i <= 1000; i++) {
		vis[i] = vis[i-1]+2;
	}
}

bool cmp(node a, node b) {
	if(a.wei == b.wei) return a.num < b.num;
	return a.wei < b.wei;
}

int main() {
	csh();
	scanf("%d", &T);
	while(T--) {
		sum = 0;
		scanf("%d %d", &n, &m);
		for(int i = 1;i <= n; i++) {
			scanf("%d", &a[i].wei);
			a[i].num = i;
			sum += a[i].wei;
		}
		
		if(m < n || n == 2) {
			puts("-1");
			continue;
		}
		
		if(m < vis[n]) {
			int ans = sum * 2;
			sort(a+1, a+1+n, cmp);
			int t = m - n;
			if(t) {
				ans += t*(a[1].wei+a[2].wei);
			}
			printf("%d\n", ans);
			for(int i = 1;i <= m; i++) {
				if(i > n) {
					printf("%d %d\n", a[1].num, a[2].num);
				}
				else {
					if(i != n) printf("%d %d\n", i, i+1);
					else printf("%d 1\n", i);
				}
			}
		}
		else {
			int ans = 0;
			sort(a+1,a+1+n,cmp);
			for(int i = 3;i <= n; i++) {
				ans += (2*a[i].wei+a[1].wei+a[2].wei);
			}
			int t = m-(n-3+1)*2;
			ans += t*(a[1].wei+a[2].wei);
			printf("%d\n", ans);
			for(int i = 3;i <= n; i++) {
				printf("%d %d\n",a[1].num, a[i].num);
				printf("%d %d\n",a[2].num, a[i].num);
			}
			for(int i = 1;i <= t; i++) {
				printf("%d %d\n", a[1].num, a[2].num);
			}
		}
	}
}