#include <bits/stdc++.h>
using namespace std;
int main(){
	ios::sync_with_stdio(0);
	cin.tie(0); cout.tie(0);
	int T; cin >> T;
	while(T--){
		int N, M, K; cin >> N >> M >> K;
		char C[40][40];
		for(int i = 0; i <= (N-1) * 2; i++){
			for(int j = 0; j <= (M-1) * 2; j++){
				C[i][j] = '#';
			}
		}
		if(K < N + M - 2 || K % 2 != (N + M - 2) % 2){
			cout << "NO\n"; continue;
		}else if(K % 4 == (N + M - 2) % 4){
			int now = 0;
			for(int j = 1; j < (M-1)*2; j+=2){
				if(now == 0) C[0][j] = 'R';
				else C[0][j] = 'B';
				now = 1 - now;
			}
			for(int i = 1; i < (N-1)*2; i+=2){
				if(now == 0) C[i][(M-1)*2] = 'R';
				else C[i][(M-1)*2] = 'B';
				now = 1 - now;
			}
			if(now == 0){
				C[N*2-3][M*2-2] = 'B';
				C[N*2-3][M*2-4] = 'B';
				C[N*2-4][M*2-3] = 'R';
				C[N*2-2][M*2-3] = 'R';
			}else{
				C[N*2-3][M*2-2] = 'R';
				C[N*2-3][M*2-4] = 'R';
				C[N*2-4][M*2-3] = 'B';
				C[N*2-2][M*2-3] = 'B';
			}
		}else{
			int now = 0;
			for(int j = 1; j < (M-1)*2; j+=2){
				if(now == 0) C[0][j] = 'R';
				else C[0][j] = 'B';
				now = 1 - now;
			}
			for(int i = 1; i < (N-1)*2; i+=2){
				if(now == 0) C[i][(M-1)*2] = 'R';
				else C[i][(M-1)*2] = 'B';
				now = 1 - now;
			}
			if(now == 1){
				C[N*2-3][M*2-2] = 'B';
				C[N*2-3][M*2-4] = 'B';
				C[N*2-4][M*2-3] = 'R';
				C[N*2-2][M*2-3] = 'R';
			}else{
				C[N*2-3][M*2-2] = 'R';
				C[N*2-3][M*2-4] = 'R';
				C[N*2-4][M*2-3] = 'B';
				C[N*2-2][M*2-3] = 'B';
			}
		}
		cout << "YES\n";
		for(int i = 0; i <= N*2-2; i+=2){
			for(int j = 1; j <= M*2-2; j+=2){
				if(C[i][j] == 'B') cout << "B ";
				else cout << "R ";
			}
			cout << "\n";
		}
		for(int i = 1; i <= N*2-2; i+=2){
			for(int j = 0; j <= M*2-2; j+=2){
				if(C[i][j] == 'B') cout << "B ";
				else cout << "R ";
			}
			cout << "\n";
		}
	}
}