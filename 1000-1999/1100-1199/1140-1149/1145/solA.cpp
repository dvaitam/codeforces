#include <iostream>
using namespace std;
int a[20], n;
int check(int l, int r) {
	bool q = true;
	for(int i = l; i <= r; i++){
		if(a[i] < a[i-1]) {
			q = false;
		}
	}
	return q;
}
int fun(int l, int r) {
	if(l == r) {
		return 1;
	}
	if(check(l, r)) {
		return r + 1 - l;
	}
	int m = (l+r)/2;
	int left, right;
	left = fun(l, m);
	right = fun(m+1, r);
	if(left > right) return left;
	else return right;
}
int main () {
	cin >> n;
	for(int i = 0; i < n; i++){
		cin >> a[i];
	}
	cout << fun(0, n-1);
}