#include <bits/stdc++.h>



using namespace std;



#ifdef LOCAL

#include "debug.h"

#else

#define dbg(...)

#endif



void solve(){

	//3 indicies which are odd

	//2 even, 1 odd

	int n;

	cin >> n;

	vector<int>odd, even;

	for (int i = 1; i <= n; i++) {

		int a;

		cin >> a;

		if (a&1){

			odd.push_back(i);

		}

		else {

			even.push_back(i);

		}

	}

	int ans[3];

	bool ok = false;

	if (odd.size() >= 3){

		ok = true;

		for (int i = 0;i < 3; i++){

			ans[i] = odd[i];

		} 

	}

	else if (odd.size() && even.size() >= 2){

		ok = true;

		ans[0] = odd[0];

		for (int i = 0;i < 2; i++){

			ans[i+1] = even[i];

		}

	}

	cout << (ok ? "YES" : "NO") << '\n';

	if (ok){

		for (int i = 0;i < 3; i++){

			cout << ans[i] << " \n"[i == 2];

		}

	}

}





int main(){

	ios::sync_with_stdio(false);

	cin.tie(0);



    int tc = 1;

    cin >> tc;

    for (int i = 1; i <= tc; i++){

        //cout << "Case #" << i << ": ";

        solve();

    }

    

    return 0;

}