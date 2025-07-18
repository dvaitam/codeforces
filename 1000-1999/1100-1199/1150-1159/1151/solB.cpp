#include<bits/stdc++.h>
using namespace std;
int a[505][505],pos[505][2];
int read() {
	char c=getchar(); int x=0,f=1;
	while (c<'0' || c>'9') {
		if (c=='-') f=-1;
		c=getchar();
	}
	while (c>='0' && c<='9') {
		x=x*10+c-'0';
		c=getchar();
	}
	return x*f;
}
int main() {
	int n=read(),m=read();
	for (int i=1; i<=n; i++) {
		for (int j=1; j<=m; j++) {
			a[i][j]=read();
		}
	}
	for (int k=0; k<=9; k++) {
		int s0=0,s1=0,s01=0;
		for (int i=1; i<=n; i++) {
			pos[i][0]=pos[i][1]=0;
			for (int j=1; j<=m; j++) {
				if ((a[i][j]>>k)&1) pos[i][1]=j;
				else pos[i][0]=j;
			}
			if (pos[i][0]==0) ++s1;
			else if (pos[i][1]==0) ++s0;
			else ++s01;
		}
		if (s1&1) {
			printf("TAK\n");
			for (int i=1; i<=n; i++) {
				if (pos[i][0]==0) printf("%d ", pos[i][1]);
				else printf("%d ", pos[i][0]);
			}
			exit(0);
		}
		else if (s01>0) {
			printf("TAK\n");
			bool flag=true;
			for (int i=1; i<=n; i++) {
				if (pos[i][0]==0) printf("%d ", pos[i][1]);
				else if (pos[i][1]==0) printf("%d ", pos[i][0]);
				else if (flag) printf("%d ", pos[i][1]),flag=false;
				else printf("%d ", pos[i][0]);
			}
			exit(0);
		}
	}
	printf("NIE\n");
	return 0;
}