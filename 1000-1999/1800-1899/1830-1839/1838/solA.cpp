#include <bits/stdc++.h>
#define MAXN 105
using namespace std;
int n,arr[MAXN];
int main() {
	int T;
	cin>>T;
	while(T--) {
		cin>>n;
		int pos=-1;
		for (int i=1;i<=n;i++) {
			cin>>arr[i];
			if (arr[i]<0) pos=i;
		}
		if (pos!=-1) {
			cout<<arr[pos]<<'\n';
			continue;
		}
		int maxi=1;
		for (int i=1;i<=n;i++)
			if (arr[i]>arr[maxi])
				maxi=i;
		cout<<arr[maxi]<<'\n';
	}
	return 0;
}