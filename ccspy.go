package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var extensions = []string{".c", ".cpp", ".C", ".cc", ".CPP", ".c++", ".cp", ".cxx"}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func writeLine(f *os.File, line string) error {
	_, err := f.WriteString(line + "\n")
	return err
}

func writeCommandData(targetDir string, cwd string, sourceFileName string, args []string) {
	if !filepath.IsAbs(sourceFileName) {
		sourceFileName = filepath.Join(cwd, sourceFileName)
	}
	sourceFileName = filepath.Clean(sourceFileName)
	
	var fileName = getMD5Hash(sourceFileName) + ".txt"
	var filePath = filepath.Join(targetDir, fileName)

	// create file
	f, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		return
	}

	// remember to close the file
	defer f.Close()

	err = writeLine(f, cwd)
	if err != nil {
		log.Println(err)
		return
	}
	err = writeLine(f, sourceFileName)
	if err != nil {
		log.Println(err)
		return
	}
	for _, opt := range args {
		err = writeLine(f, opt)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func main() {
	var defaultCCompiler = os.Getenv("CCSPY_CC")
	var defaultCppCompiler = os.Getenv("CCSPY_CXX")
	var defaultTargetDir = os.Getenv("CCSPY_TARGET_DIR")
	var args = os.Args[1:]
	var targetDirectory = defaultTargetDir
	var compilerCommand string

	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Look for ccspy specific options - they must appear before all other options
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
	if len(compilerCommand) == 0 && len(defaultCCompiler) == 0 && len(defaultCppCompiler) == 0 {
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

	// Make sure the target directory exists
	_, err = os.Stat(targetDirectory)
	if os.IsNotExist(err) {
		err = os.Mkdir(targetDirectory, 0o755)
		if err != nil {
			log.Fatal("Cannot create directory: " + targetDirectory)
		}
	}

	// Separate options from compilation units (sources)
	var sources = make([]string, 0, 1)
	var argsWithoutSources = make([]string, 0, len(args))
	var cppCount = 0
	var cCount = 0

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			argsWithoutSources = append(argsWithoutSources, arg)
			continue
		}
		isSource := false
		for _, ext := range extensions {
			if strings.HasSuffix(arg, ext) {
				sources = append(sources, arg)
				if ext == ".c" {
					cCount++
				} else {
					cppCount++
				}
				isSource = true
				break
			}
		}
		if !isSource {
			argsWithoutSources = append(argsWithoutSources, arg)
		}
	}

	var wg sync.WaitGroup

	// Now log the options for each compilation unit in a separate file in a parallel running go-routine
	wg.Add(len(sources))
	for _, src := range sources {
		go func(src string) {
			writeCommandData(targetDirectory, cwd, src, argsWithoutSources)
			wg.Done()
		}(src)
	}

	// Decide if to use the C or the C++ compiler
	if len(compilerCommand) == 0 {
		if cppCount > 0 {
			compilerCommand = defaultCppCompiler
		} else if cCount > 0 {
			compilerCommand = defaultCCompiler
		} else {
			// If there is no source file default to the C++ compiler
			compilerCommand = defaultCppCompiler
		}
	}

	// Now call the real compiler
	var cmd = exec.Command(compilerCommand, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	wg.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			var exitCode = exitError.ExitCode()
			os.Exit(exitCode)
		}
	}
}
