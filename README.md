# go-jdk-manager
A simple JDK manager written in Go

The intent of this project is to be a learning exercise in building CLI applications in Go and not an example of clean or idiomatic Go.

## How to install
1. Have go 1.26.1 or greater installed on your local machine
2. Clone the repo
```bash
git clone git@github.com:saurabh-ku/go-jdk-manager.git
```
3. Get modules
```bash
go mod download
```
3. Build
```bash
go build main.go
```
4. Run
```bash
./main
```

## Features
### List version
![list.png](assets%2Flist.png)

### Download 
The manager will automatically download the selected JDK version and set up a symlink it to the jdk link.
Setting the `JAVA_HOME` to jdk in the bash profile will then allow the manager to change the underlying jdk version without the user having to change the bash profile. 
![download.png](assets%2Fdownload.png)

### Switch
JDK versions which have already been downloaded can we switched to seamlessly. The manager will simply change the symlink to the version selected. 
![switch.png](assets%2Fswitch.png)