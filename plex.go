/*
This is a general purpose GO program that takes a list of
programs or scripts and runs them concurrently. The maximum
number allowed to run concurrently can be throttled by setting the
max-ops option. Once the limit is reached, the program
simply blocks until one of the program exits.
If no limit is given, all are started simultaneously.

Each script/program instance runs in a separate temporary directory.

The list of progams is supplied as a file from the command line
or may be piped to stdin.

Example :
   go build
   echo 'pwd' >  test.sh 
   echo 'pwd' >> test.sh 
   echo 'pwd' >> test.sh 

	 ./plex test.sh

	 echo '     -OR-    '

	 cat test.sh | ./plex

*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
)

var maxOps int
var tmpRoot string
var prefix string
var args []string

func init() {
	flag.IntVar(&maxOps, "max-ops", 10, "maximum concurrent operations")
	flag.StringVar(&tmpRoot, "workdir", "", "temporary working dir root")
	flag.StringVar(&prefix, "prefix", "plex_", "temporary working dir prefix")
	flag.Parse()
	args = flag.Args()
}

func doit(wdir string, cStr string) {
	var cmdStr string
	cmdStr = fmt.Sprintf("( cd /%s;  %s)", wdir, cStr)
	fmt.Println(cmdStr)
	cmd := exec.Command("sh", "-c", cmdStr)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}

func main() {
	var err error
	var f *os.File

	switch flag.NArg() {
	case 0:
		f = os.Stdin
		break

	case 1:
		f, err = os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		break
	default:
		fmt.Printf("input must be from stdin or file\n")
		os.Exit(1)
	}
	var wg sync.WaitGroup
	const total = 120
	worker := make(chan string)
	limiter := make(chan int, maxOps)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		limiter <- 0 // block here if necessary

		go func() {

			// create temporary working directory
			tmpdir, err := ioutil.TempDir(tmpRoot, prefix) // create working directory
			if err != nil {
				log.Fatal(err)
			}

			//------------------------------------------------------
			defer func() {
				os.RemoveAll(tmpdir) // remove temp working directory
				<-limiter            // release buffered channel slot
				wg.Done()            // signal routine done
			}()
			//------------------------------------------------------

			cmdStr := <-worker // get data

			doit(tmpdir, cmdStr)
		}()

		wg.Add(1)
		worker <- scanner.Text() // write program name/path to channel
	}

	fmt.Printf("Waiting for final routines ...\n")
	wg.Wait() // wait for all routines to complete
	fmt.Printf("done.\n")

}
