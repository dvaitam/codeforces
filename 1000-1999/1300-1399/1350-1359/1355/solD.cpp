#include <bits/stdc++.h>

 

using namespace std;

 

typedef long long ll;

#define fast ios::sync_with_stdio(0);cin.tie(0);

 

int main() {

	fast



	int N, S;

	cin >> N >> S;

	

	int biggest = S - N + 1;

	if (S - biggest + 1 >= biggest)

		cout << "NO" << endl;

	else {

		cout << "YES" << endl << biggest << " ";

		for (int i = 0; i < N-1; i++)

			cout << 1 << " ";

		cout << endl << S - biggest + 1 << endl;

	}



}