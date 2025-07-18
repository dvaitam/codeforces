#include <iostream>
using namespace std;


int main() {
	int t;
	cin>>t;
	while(t--){
	    int n;
	    cin>>n;
	    if(n==1) cout<<1<<endl;
	    else if(n==2) cout<<1<<" "<<2<<endl;
        else{
            cout<<2<<" ";
            for(int i=1;i<=(n-3)/2;i++){
                cout<<i+3<<" ";
            }
            cout<<1<<" ";
            for(int i=1;i<=(n-2)/2;i++){
                cout<<i+3+(n-3)/2<<" ";
            }
            cout<<3<<endl;
        }
	}
	return 0;
}