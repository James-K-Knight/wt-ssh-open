package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/James-K-Knight/wt-ssh-open/sshOpt"
)

// Parse the ssh connection parameters based off SSH URI draft
// https://datatracker.ietf.org/doc/html/draft-ietf-secsh-scp-sftp-ssh-uri
func ParseParams(params string) string {
	var parameters string

	if strings.Contains(params, ",") {
		for _, cparams := range strings.Split(params, ",") {
			if strings.Contains(cparams, "=") && sshOpt.Validate(strings.Split(params, "=")[0]) {
				parameters += "-o " + cparams
				log.Print("Option Added: " + cparams)
			}
		}
	} else {
		if strings.Contains(params, "=") && sshOpt.Validate(strings.Split(params, "=")[0]) {
			parameters += "-o " + params
			log.Print("Option Added: " + params)
		}
	}

	return parameters
}

func SplitUri(input string) (string, string) {
	var sshOptions string

	if strings.Contains(input, ";") {
		sp1 := strings.SplitN(input, ";", 2)
		sshOptions = ParseParams(sp1[1])
		return sp1[0], sshOptions
	}
	return input, sshOptions
}

// Check if Windows Terminal, openSSH and URI are in order
func SanityCheck(wtpath string, parsed *url.URL) {
	var err error

	if _, err = os.Stat(wtpath); os.IsNotExist(err) {
		log.Fatal("[ERROR] windows terminal not found ", err)
	}
	if _, err = exec.LookPath("ssh.exe"); err != nil {
		log.Fatal("[ERROR] windows terminal not found ", err)
	}

	if parsed.Scheme != "ssh" {
		log.Fatal("[ERROR] URI scheme is invalid")
	}
	if _, userpassword := parsed.User.Password(); userpassword {
		log.Fatal("[ERROR] Insecure URI contains password")
	}
}

// Bypass WindowsAPP 0kb bug with os.exec source:https://stackoverflow.com/a/71874620
func RunTerminal(exepath string, arguments []string) {
	var err error

	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{nil, nil, nil}
	/*
		To redirect IO, pass in stdin, stdout, stderr as required
		procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	*/

	// The argv slice will become os.Args in the new process,
	// so it normally starts with the program name
	if _, err = os.StartProcess(exepath, append([]string{exepath}, arguments...), procAttr); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error
	var verbose bool
	var parsed *url.URL
	var username string
	var sshOptions string
	var sshArguments []string

	flag.BoolVar(&verbose, "v", false, "Run SSH.exe in Windows terminal verbosly")
	flag.Parse()

	values := flag.Args()

	if len(values) == 0 {
		log.Fatal(values)
		fmt.Println("Usage: wtssh.exe [-v] ssh://[<user>[;ConnectTimeout=<Timeout>]@]<host>[:<port>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	parsed, err = url.Parse(values[0])

	if err != nil {
		log.Fatal(err)
	}

	username, sshOptions = SplitUri(parsed.User.Username())
	wtPath := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft\\WindowsApps\\wt.exe")

	SanityCheck(wtPath, parsed)

	sshArguments = append(sshArguments, "new-tab")
	sshArguments = append(sshArguments, "ssh.exe")

	if verbose {
		sshArguments = append(sshArguments, "-v")
	}

	sshArguments = append(sshArguments, parsed.Hostname())
	log.Print("Hostname: " + parsed.Hostname())

	if parsed.Port() != "" {
		sshArguments = append(sshArguments, "-p")
		sshArguments = append(sshArguments, parsed.Port())
		log.Print("Using port: " + parsed.Port())
	} else {
		log.Print("Using standard port")
	}

	if username != "" {
		sshArguments = append(sshArguments, "-l")
		sshArguments = append(sshArguments, username)
		log.Print("Username: " + username)
	}

	// Often / is appended to path, if so ignore
	if parsed.Path != "" && parsed.Path != "/" {
		sshArguments = append(sshArguments, parsed.Path)
		log.Print("Command: " + parsed.Path)
	}

	if sshOptions != "" {
		sshArguments = append(sshArguments, sshOptions)
	}

	RunTerminal(wtPath, sshArguments)
}
