#include <iostream>

#include <string>



using namespace std;



int main()

{

    ios_base::sync_with_stdio(false); cin.tie(NULL); cout.tie(NULL);

    int T; cin>>T;



    while (T--){

        int N; cin>>N;

        string str; cin>>str;

        int odd = 0, even = 0;

        for (int i=0; i<N; i++)

            if (str[i] == '1') odd += 1;

        if (odd == 0 || odd & 1){

            cout<<"NO\n";

            continue;

        }

        cout<<"YES\n";



        int dx = 0;

        if (str[N-1] != '1'){

            while (str[dx] == '0') dx += 1;

            dx += 1;

        }



        for (int i=0; i<N; i++){

            int j = i;

            if (i) cout<<dx + 1<<" "<<(i + dx)%N + 1<<"\n";

            while (str[(j + dx)%N] == '0'){

                cout<<(j + dx)%N + 1<<" "<<(j +dx + 1)%N + 1<<"\n";

                j += 1;

            }

            i = j;

        }



    }

    return 0;

}