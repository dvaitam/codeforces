#include <cstdio>

using namespace std;

int main() {
	int n;
	scanf("%d",&n);
	for(int i=0;i<n;i++) putc("abcd"[i%4],stdout);
	printf("\n");
	return 0;
}