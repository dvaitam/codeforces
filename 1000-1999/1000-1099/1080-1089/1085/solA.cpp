#include <bits/stdc++.h>

using namespace std;
int main () {
    string s, n; cin >> s;
    int length = s.size();
    if(length % 2 == 0){
    while(length){
        if(length % 2 == 0){
            n += s[length / 2 - 1];
            s.erase((length - 1) / 2, 1);
            length--;
        }
        else{
            n += s[length / 2];
            s.erase((length - 1) / 2, 1);
            length--;
        }
    }
    }
    else {
        while(length){
            n += s[length / 2];
            s.erase(length / 2, 1);
            length--;
        }
    }
    cout << n;
}