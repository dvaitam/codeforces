#include<bits/stdc++.h>

#define start ios_base::sync_with_stdio(false);cin.tie(0);cout.tie(0)

#define ll long long

#define pb push_back

#define po pop_back

#define mp make_pair

#define fi first

#define se second

using namespace std;

vector <int> prime;

int T, N, K, a, b, c;



int main(){

	start;

	cin >> T;

	while(T--){

		cin >> N >> K;

		if(N % 2 == 1) cout << "1 " << (N-1)/2 << " " << (N-1)/2;

		else{

			if((N/2) % 2 == 0){

				cout << N/2 << " ";

				N /= 2;

				cout << N/2 << " " << N/2;

			}else{

				cout << (N/2)-1 << " " << (N/2)-1 << " " << 2;

			}

		}

		cout << "\n";

	}

}