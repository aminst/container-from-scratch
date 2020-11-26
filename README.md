# simple-container-with-go
<b>A simple container implementation using go lang.</b>

It is based on this [nice talk](https://youtube.com/watch?v=8fi7uSYlOdc  ) by Liz Rice.  
I Also found lots of useful information about namespaces reading [these articles](https://medium.com/@teddyking/linux-namespaces-850489d3ccf) by Ed King.  
You may also want to use cgroups to limit resources and other stuff.

### Usage
You will need to prepare a temporary file system within a directory called `newroot`.  
This is used by the container for the mount namespace.  
Run the container using `go run main.go <command> <arguments>`  
For example `go run main.go bash`
