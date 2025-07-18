#include <bits/stdc++.h>



using namespace std;



int a[110];



int main(){

	cin.tie(0), cout.tie(0);

	ios::sync_with_stdio(false);

	int t;

	cin >> t;

	while(t --){

		int n;

		cin >> n;

		for(int i = 0; i < n; i ++)

			cin >> a[i];

		for(int i = 0; i < n; i ++){

			int x = 0;

			for(int j = 0; j < n; j ++)

				if(j != i)	x ^= a[j];

			if(x == a[i]){

				cout << x << endl;

				break;

			}

		}

	}

	return 0;

}