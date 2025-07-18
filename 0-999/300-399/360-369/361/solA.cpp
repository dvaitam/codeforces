#include <iostream>
#include <vector>

using namespace std;

int main()
{
    int n , k;
    cin >> n >> k;

    vector<vector<int> > V(n , vector<int> (n , 1));
    int Put = k - n + 1 ;
    for(int i = 0 ; i < n ; ++i){
        V[i][i] = Put;
    }

    for(int i = 0 ; i < n ; ++i){
        for(int j = 0 ; j < n ; ++j){
            cout << V[i][j] << ' ';
        }
        cout << endl;
    }
    return 0;
}