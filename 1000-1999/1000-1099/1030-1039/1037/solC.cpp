#include <cstdio>

int main(){
	int n;
	scanf("%d", &n);
	char a[n+1];
	char b[n+1];
	scanf("%s%s", a, b);
	int ans = 0;
	for(int i=1; i < n; ++i){
		if(a[i-1] == b[i] && a[i] == b[i-1] && a[i-1] != a[i]){
			++ans;
			a[i-1] = b[i-1];
			a[i] = b[i];
		}
	}
	for(int i=0; i < n; ++i){
		if(a[i] != b[i]){
			++ans;
		}
	}
	printf("%d\n", ans);
	return 0;
}