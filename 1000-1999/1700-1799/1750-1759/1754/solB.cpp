#include <iostream>

using namespace std;



int main(){

    int t;

    cin >> t;

    for (int k=0; k<t; k++){

        int n;

        cin >> n;

        if (n==2){

            cout << 1 << " " << 2;

        }

        else if (n==3){

            cout << 1 << " " << 2 << " " << 3;

        }

        else {

            if (n%2==0){

                for (int i=0; i<n; i++){

                    if (i==0){

                        cout << (n/2)+1;

                    }

                    else if (i%2==0){

                        cout << " " << ((n+i)/2)+1;

                    }

                    else {

                        cout << " " << (i/2)+1;

                    }

                }

            }

            else {

                cout << n;

                n--;

                for (int i=0; i<n; i++){

                    if (i%2==0){

                        cout << " " << ((n+i)/2)+1;

                    }

                    else {

                        cout << " " << (i/2)+1;

                    }

                }

            }

        }

        cout << endl;

    }

}