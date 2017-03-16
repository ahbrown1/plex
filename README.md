# plex

This is a general purpose GO program that takes a list of programs or
scripts and runs them concurrently. The maximum number allowed to run
concurrently can be throttled by setting the max-ops option. Once the limit
is reached, the program simply blocks until one of the program exits. If no
limit is given, all are started simultaneously.

Each script/program instance runs in a separate temporary directory.

The list of progams is supplied as a file from the command line or may be
piped to stdin.

Example :

       go build 
       echo 'pwd' >  test.sh
       echo 'pwd' >> test.sh
       echo 'pwd' >> test.sh

         ./plex test.sh

        echo '     -OR-    '

        cat test.sh | ./plex


