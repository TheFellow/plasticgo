# Go parser for Plastic SCM

Plastic SCM and Semantic Merge provide a way to run a syntactic diff (and merge) based on different programming languages AST.

In order for this to work properly, a [special format](http://codicesoftware.blogspot.com/2015/09/custom-languages-in-semantic-version.html)
is required. This project will use **golang** internal parser to create the required YAML.

## Install on windows
Be sure that you have a complete installation of [golang on windows](
https://golang.org/doc/install#windows_env
)

###Install golang yaml
```
cd %GOPATH%
go get gopkg.in/yaml.v2
``` 

###Install goparser
```
mkdir src\espeleta.info
cd src\espeleta.info
git clone https://github.com/tespeleta/plasticgo.git
go install espeleta.info/goparser
cd %GOPATH%
```

From now, we can compile and test with:

```
go build espeleta.info/goparser
go test espeleta.info/goparser
```



## Use
Please refer to Plastic SCM detailed [instructions](http://codicesoftware.blogspot.com/2015/09/custom-languages-in-semantic-version.html)
. You can try go parser in the following way:

```
goparser.exe shell flag.txt
``` 

`flag.txt` is the filename used to create an empty file that plastic SCM will understand as an indication that all initializations by the parser are finished. In `golang` case, the initialization is immediate.

Once started, the parser will wait reading from standard input the filenames to parse and the output files:

```
goparser.exe shell flag.txt
\path\to\filename.go
\path\to\outputfile.yaml
OK
\path\to\filename2.go
\path\to\outputfile2.yaml
OK
``` 
The parser will indicate with `KO` if there were any error during the processing, or `OK` otherwise.
