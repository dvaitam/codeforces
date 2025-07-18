#include <iostream>
using namespace std;

int main() {
	int w, h, k;
	cin >> w >> h >> k;
	int res = 0;
	while (k != 0) {
		int sum = 2 * h + 2 * w - 4;
		h -= 4;
		w -= 4;
		k--;
		res += sum;
	}
	cout << res;
}