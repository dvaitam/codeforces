#include <bits/stdc++.h>

using namespace std;

bool prime(int n){
  int s = (int)(sqrt(n) + 1);
  for (int i = 2; i < s; i++){
    if (n % i == 0) return 0; 

  }
  return 1;
}

int nextP(int i){
  while (!prime(i)) i++;
  return i;
}

int main(){
  ios_base::sync_with_stdio(0);
  cin.tie(0);
  int N;
  cin >> N;
  int M = nextP(N);
  cout << M << endl;
  for (int i = 0; i < N - 1; i++){
    cout << i + 1 << " " << i+2 << "\n";
  }
  cout << 1 << " " << N << "\n";
  for (int i = 0; i < M - N; i++){
    cout << i + 1 << " " << i + 1 + N /2 << "\n";
  }
  
}