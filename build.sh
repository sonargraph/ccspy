rm -r bin/*
./buildForTarget.sh linux amd64
./buildForTarget.sh linux arm64
./buildForTarget.sh darwin amd64
./buildForTarget.sh darwin arm64
./buildForTarget.sh windows amd64
./buildForTarget.sh windows arm64
mv bin/windows_amd64/ccspy bin/windows_amd64/ccspy.exe
mv bin/windows_arm64/ccspy bin/windows_arm64/ccspy.exe

