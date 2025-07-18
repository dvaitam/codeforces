#include <bits/stdc++.h>

using namespace std;

int main(){
  ios_base::sync_with_stdio(0);
  cin.tie(0);
  int N;
  cin >> N;
  for (int i = 0; i < N; i++){
    int a; 
    cin >> a;
    int thou, mil, one;
    mil = a/ 1122457;//1125841;
    a %= 1122457;
    thou = (a / 1000) ;
    one = a % 1000;
    //cout << mil << " " << thou << one << endl;
    cout << 133;
    for (int j = 0; j < one; j++) cout << 7;
    if (thou > 0){
      for (int j = 0; j < 994; j++) cout << 1;
      cout << 33;
      for (int j = 0; j < thou; j++) cout << 7;
    }
    if (mil > 0){
      for (int j = 0; j < 46; j++) cout << 3;
      for (int j = 0; j< mil; j++) cout << 7;
    }
    cout << "\n";
  }  
}