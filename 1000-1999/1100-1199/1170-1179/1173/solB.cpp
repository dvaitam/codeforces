#include<bits/stdc++.h>
using namespace std;

long long n;

int main(){
	cin >> n;
	cout << n/2+1 << "\n";
	long long r=1, c=1;
	for (int i=0;i<n;i++){
		cout << r << " " << c << "\n";
		if (i%2) r++;
		else c++;
	}
}