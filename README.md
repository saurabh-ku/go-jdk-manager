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
![Screenshot 2024-01-22 at 10.19.27 PM.png](..%2F..%2F..%2F..%2F..%2Fvar%2Ffolders%2F9_%2Ft2yhfnvx22sbqn5wwqw1ksx00000gn%2FT%2FTemporaryItems%2FNSIRD_screencaptureui_P5areR%2FScreenshot%202024-01-22%20at%2010.19.27%E2%80%AFPM.png)

### Download 
The manager will automatically download the selected JDK version and set up a symlink it to the jdk link.
Setting the `JAVA_HOME` to jdk in the bash profile will then allow the manager to change the underlying jdk version without the user having to change the bash profile. 
![Screenshot 2024-01-22 at 10.20.15 PM.png](..%2F..%2F..%2F..%2F..%2Fvar%2Ffolders%2F9_%2Ft2yhfnvx22sbqn5wwqw1ksx00000gn%2FT%2FTemporaryItems%2FNSIRD_screencaptureui_nYbueI%2FScreenshot%202024-01-22%20at%2010.20.15%E2%80%AFPM.png)

### Switch
JDK versions which have already been downloaded can we switched to seamlessly. The manager will simply change the symlink to the version selected. 
![Screenshot 2024-01-22 at 10.23.31 PM.png](..%2F..%2F..%2F..%2F..%2Fvar%2Ffolders%2F9_%2Ft2yhfnvx22sbqn5wwqw1ksx00000gn%2FT%2FTemporaryItems%2FNSIRD_screencaptureui_L1xCVm%2FScreenshot%202024-01-22%20at%2010.23.31%E2%80%AFPM.png)