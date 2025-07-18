#include <bits/stdc++.h>

using namespace std;
using LL = long long;

const int MAXN = 1000 + 3;

int n, k;
int a[MAXN][MAXN];

bool work(){
  cin >> n >> k;
  for(int i = 1; i <= n; i++){
    for(int j = 1; j <= n; j++) a[i][j] = 0;
  }
  if(k % 2 == 1) return 0;
  if(k % 4 == 0){
    cout << "Yes\n";
    for(int i = 1; i + 1 <= n && k > 0; i += 2){
      for(int j = 1; j + 1 <= n && k > 0; j += 2){
        for(int x = i; x < i + 2; x++){
          for(int y = j; y < j + 2; y++){
            a[x][y] = 1; 
          }
        }
        k -= 4;
      }
    }
    for(int i = 1; i <= n; i++){
      for(int j = 1; j <= n; j++) cout << a[i][j] << " ";
      cout << "\n";
    }
    return 1;
  }
  if(6 <= k && k <= n * n - 10 && k % 4 == 2){
    cout << "Yes\n";
    k -= 6;
    for(int i = 1; i + 1 <= n && k > 0; i += 2){
      for(int j = (i <= 3 ? 4 : 1); j + 1 <= n && k > 0; j += 2){
        for(int x = i; x < i + 2; x++){
          for(int y = j; y < j + 2; y++){
            a[x][y] = 1; 
          }
        }
        k -= 4;
      }
    }
    a[1][1] = a[1][2] = a[2][1] = a[2][3] = a[3][2] = a[3][3] = 1;
    for(int i = 1; i <= n; i++){
      for(int j = 1; j <= n; j++) cout << a[i][j] << " ";
      cout << "\n";
    }
    return 1;
  }
  if(k == n){
    cout << "Yes\n";
    for(int i = 1; i <= n; i++){
      for(int j = 1; j <= n; j++){
        cout << (j == 1 ? 1 : 0) << " ";
      }
      cout << "\n";
    }
    return 1;
  }  
  if(k == n * n - 6){
    cout << "Yes\n";
    a[1][1] = a[1][2] = a[2][1] = a[2][3] = a[3][1] = a[3][2] = a[3][3] = a[3][4] = a[4][1] = a[4][4] = 1;
    for(int i = 1; i <= n; i++){
      for(int j = 1; j <= n; j++) cout << (i > 4 || j > 4 ? 1 : a[i][j]) << " ";
      cout << "\n";
    }
    return 1;
  }
  return 0; 
}

int main(){
  ios::sync_with_stdio(0), cin.tie(0);
  int T;
  cin >> T;
  while(T--){
    cout << (work() == 0 ? "No\n" : "");
  }
  return 0;
} 
// 6 14