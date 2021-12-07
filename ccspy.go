package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func getExtensions() []string {
	return []string{".c", ".cpp", ".C", ".cc", ".CPP", ".c++", ".cp", ".cxx"}
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func writeLine(f *os.File, line string) {
	_, err := f.WriteString(line + "\n")
	if err != nil {
		log.Fatal(err)
	}
}

func writeCommandData(targetDir string, cwd string, sourceFileName string, args []string) {

	if !path.IsAbs(sourceFileName) {
		sourceFileName = path.Join(cwd, sourceFileName)
	}

	var fileName = getMD5Hash(sourceFileName) + ".txt"
	var filePath = path.Join(targetDir, fileName)

	// create file
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file
	defer f.Close()

	writeLine(f, cwd)
	writeLine(f, sourceFileName)
	writeLine(f, strings.Join(args, "$"))
}

func main() {
	var defaultCompiler = os.Getenv("CCSPY_COMPILER")
	var defaultTargetDir = os.Getenv("CCSPY_TARGET_DIR")
	var args = os.Args[1:]
	var compilerCommand = defaultCompiler
	var targetDirectory = defaultTargetDir

	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var counter = 0

	for _, v := range args {
		if strings.HasPrefix(v, "-ccspyCompiler=") {
			compilerCommand = strings.Split(v, "=")[1]
		} else if strings.HasPrefix(v, "-ccspyTargetDir=") {
			targetDirectory = strings.Split(v, "=")[1]
		} else {
			break
		}
		counter++
	}
	args = args[counter:]
	if len(compilerCommand) == 0 {
		// Use prefix mode, ccspy is just inserted as the first element of the command line
		if len(args) == 0 {
			log.Fatal("ccspy requires at least one parameter")
		}
		compilerCommand = args[0]
		args = args[1:]
	}
	if len(targetDirectory) == 0 {
		log.Fatal("You must define the target directory either via '-ccspyTargetDir=...' or via environment variable CCSPY_TARGET_DIR")
	}
	_, err = os.Stat(targetDirectory)
	if os.IsNotExist(err) {
		err = os.Mkdir(targetDirectory, 0o755)
		if err != nil {
			log.Fatal("Cannot create directory: " + targetDirectory)
		}
	}

	var sources = make([]string, 0, 1)
	var argsWithoutSources = make([]string, 0, len(args))

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			argsWithoutSources = append(argsWithoutSources, arg)
			continue
		}
		for _, ext := range getExtensions() {
			if strings.HasSuffix(arg, ext) {
				sources = append(sources, arg)
			} else {
				argsWithoutSources = append(argsWithoutSources, arg)
			}
		}
	}
	for _, src := range sources {
		writeCommandData(targetDirectory, cwd, src, argsWithoutSources)
	}

	var cmd = exec.Command(compilerCommand, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			var exitCode = exitError.ExitCode()
			os.Exit(exitCode)
		}
	}
}
