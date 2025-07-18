# include <bits/stdc++.h>
# define ll long long
# define fi first
# define se second
# define pii pair<int, int>
using namespace std;

int main() {
	int cases;
	scanf("%d", &cases);
	while(cases--) {
		int N;
		scanf("%d", &N);
		vector<int> arr(N);
		
		for(int i=0;i<N;i++) {
			scanf("%d", &arr[i]);
		}
		
		if(arr[0] == arr[N - 1]) printf("NO\n");
		else {
			printf("YES\n");
			for(int i=0;i<N;i++) {
				if(i == 1) printf("R");
				else printf("B");
			}
			printf("\n");
		}
	}
	
}