#include <bits/stdc++.h>
#define int long long

using namespace std;

signed main() {
    int n=0;
    int m=0;
    int a=0;
    int b=0;
    int i=0;
    int j=0;
    string s;
    
    cin>>n>>m;
    vector<string> vec(n);
    vector<vector<bool>> vis(n, vector<bool>(m, false));
    
    for(i=0;i<n;i++) {
        cin>>s;
        vec.at(i) = s;
        
        
        if(a == 0) {
            for(j=0;j<m;j++) {
                if(s[j] == 'S') {
                    a = i;
                    b = j;
                }
            }
        }
    }
    
    i = a;
    j = b;
    
    while(true) {
        vis.at(i).at(j) = true;
        
        if(i < (n-1) && vec.at(i+1).at(j) == '*' && !vis.at(i+1).at(j)) {
            i++;
            cout<<'D';
        }
        
        else if(i > 0 && vec.at(i-1).at(j) == '*' && !vis.at(i-1).at(j)) {
            i--;
            cout<<'U';
        }
        
        else if(j < (m-1) && vec.at(i).at(j+1) == '*' && !vis.at(i).at(j+1)) {
            j++;
            cout<<'R';
        }
        
        else if(j > 0 && vec.at(i).at(j-1) == '*' && !vis.at(i).at(j-1)) {
            j--;
            cout<<'L';
        }
        
        else if(abs(i-a) + abs(j-b) == 1) {
            if(a-i == 1)
            cout<<'D';
            
            else if(a-i == -1)
            cout<<'U';
            
            else if(b-j == 1)
            cout<<'R';
            
            else
            cout<<'L';
            
            break;
        }
    }
    
    cout<<endl;
    
    return 0;
}