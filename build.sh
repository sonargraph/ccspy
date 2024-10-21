rm -r bin/*
./buildForTarget.sh linux amd64 linux-x86_64
./buildForTarget.sh linux arm64 linux-aarch64
./buildForTarget.sh darwin amd64 mac-x86_64
./buildForTarget.sh darwin arm64 mac-aarch64
./buildForTarget.sh windows amd64 windows-x86_64
./buildForTarget.sh windows arm64 windows-aarch64
mv bin/windows-x86_64/ccspy bin/windows-x86_64/ccspy.exe
mv bin/windows-aarch64/ccspy bin/windows-aarch64/ccspy.exe

