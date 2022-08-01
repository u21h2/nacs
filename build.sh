rm -rf build/*
mkdir -p build/pocs
mkdir build/nacs_linux_amd64
mkdir build/nacs_linux_arm64
mkdir build/nacs_darwin_amd64
mkdir build/nacs_darwin_arm64
mkdir build/nacs_win_amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w " -trimpath -o build/nacs_linux_amd64/nacs
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w " -trimpath -o build/nacs_linux_arm64/nacs
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w " -trimpath -o build/nacs_win_amd64/nacs.exe
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w " -trimpath -o build/nacs_darwin_amd64/nacs
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w " -trimpath -o build/nacs_darwin_arm64/nacs
cp -r web/pocv1/pocs build/pocs/xrayv1
cp -r web/poc/pocs/nuclei build/pocs/nuclei
cd build/
cp -r pocs nacs_linux_amd64
cp -r pocs nacs_linux_arm64
cp -r pocs nacs_darwin_amd64
cp -r pocs nacs_darwin_arm64
cp -r pocs nacs_win_amd64
COPYFILE_DISABLE=true tar -zcvf nacs_darwin_amd64.tar.gz --exclude="*.DS_Store"  nacs_darwin_amd64
COPYFILE_DISABLE=true tar -zcvf nacs_linux_arm64.tar.gz --exclude="*.DS_Store"  nacs_linux_arm64
COPYFILE_DISABLE=true tar -zcvf nacs_darwin_arm64.tar.gz --exclude="*.DS_Store"  nacs_darwin_arm64
COPYFILE_DISABLE=true tar -zcvf nacs_linux_amd64.tar.gz --exclude="*.DS_Store" nacs_linux_amd64
zip -q -r nacs_win_amd64.zip nacs_win_amd64