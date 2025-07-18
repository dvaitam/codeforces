#include <iostream>


using namespace std;

int n,k;

int main() 
{

    
    cin>>n>>k;

    int a = 1;
    int b = n;

    k--;

    int md = 0 ;
    for (int i=1;i<=n;i++) {
        if (k==0) {
            if (md==0) {
                if (i%2 == 0) {
                    md = 1;
                } else 
                    md = 2;
            }

            if (md == 1) {
                cout<<b<<' ';
                b--;
            }
            if (md == 2) {
                cout<<a<<' ';
                a++;
            }
        } else {
        
            if (i%2==0) {
                cout<<a<<' ';
                a++;
            } else {
                cout<<b<<' ';
                b--;
            }

            if (i>1) k--;
        }
    }

}