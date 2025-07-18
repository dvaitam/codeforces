#include <iostream>

#include <cmath>

using std::cin, std::cout;

 

int main(){

    int t,n,n1,m,a,b,x,y,c;

    cin >> t;

    while(t--){

        cin>> n;

        n1 = n;

        b=0,a=0,x=0,y=0,m=3;

        if(n%2 == 0){

            a = 2;

            while(n%2 == 0){

                n /= 2;

                x++;

            }

        }

        else{

            while(m*m < n1){

                while(n%m == 0){

                    n /= m;

                    x++;

                }

                if(x != 0){a = m; break;}

                m += 2;

                

            }

        }

        if((n == 1 && x < 6) || m*m >= n1){cout<<"NO"<<'\n';}

        else if(n == 1){

            cout<<"YES"<<'\n';

            b = a*a;

            c = pow(a,x-3);

            cout<< a <<' '<< b <<' '<< c <<'\n';}

        else if(x > 2){

            cout<<"YES"<<'\n';

            b = pow(a,x-1);

            c = n1/pow(a,x);

            cout<< a <<' '<< b <<' '<< c <<'\n';}

        else{

            while(m*m < n1){

                while(n%m == 0){

                    n /= m;

                    y++;

                }

                if(y != 0){b = m; break;}

                m += 2;

            }

            if((n == 1 && x+y < 4) || m*m >= n1){cout<<"NO"<<'\n';}

        else if(n == 1 && x==1){

            cout<<"YES"<<'\n';

            c = pow(b,y-1);

            cout<< a <<' '<< b <<' '<< c <<'\n';}

        else if(n == 1 && y==1){

            cout<<"YES"<<'\n';

            c = pow(a,x-1);

            cout<< a <<' '<< b <<' '<< c <<'\n';}

        else if(n == 1){

            cout<<"YES"<<'\n';

            c = a*b;

            a = pow(a,x-1);

            b = pow(b,y-1);

            cout<< a <<' '<< b <<' '<< c <<'\n';}

        else{

            cout<<"YES"<<'\n';

            a = pow(a,x);

            b = pow(b,y);

            cout<< a <<' '<< b <<' '<< n <<'\n';

            }

        }  

    }

}