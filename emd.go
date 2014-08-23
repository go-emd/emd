/*
	emd - Contains utilities that allow users to create new emd 
	projects, compile, distribute, start, and monitor on multiple machines.
	
	The commands available through this executable are:
	
	emd new: Performs a git clone of a boilerplate repository that is 
	already set up to run on a distribution therefore allowing quick 
	and painless starting of projects.
	
	emd compile --path <path to folder containing distribution>: Will 
	compile a distribution project by first parsing the config.json 
	file creating node leader go files then building them with 
	"go build".
	
	emd distribute --path <path to folder containing distribution>: Distributes 
	the distribution using the "rsync" command into the tmp directory of the machine.
	
	emd start --path <path to folder containing distribution>: Starts the 
	distribution by ssh'ing to each individual node in the distribution 
	and starting is node leader in the background.  NOTE: running the compile and 
	distribute commands will need to be done before this one.
	
	emd stop --path <path to folder containing distribution>: Stops the distribution 
	by sending REST calls to the endpoints of each node leader in the distribution 
	therefore killing each process and affectively stopping it.
	
	emd status --path <path to folder containing distribution>: Returns the status 
	of the distribution when running in json format.  It will either be "Healthy", 
	"Unhealthy" or "Unknown"
	
	emd metrics --path <path to folder containing distribution>: Returns the metrics 
	created by the distribution that is running in json format.
*/
package main

import (
	"github.com/go-emd/emd/config"
	"github.com/go-emd/emd/log"
	"code.google.com/p/go.crypto/ssh"
	"github.com/howeyc/gopass"
	"fmt"
	"strings"
	"bytes"
	"runtime"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"net/http"
	"text/template"
)

// eq: Is used to allow boolean logic between string using a 
// template leader file.
func eq(s1, s2 string) bool {
	return s1 == s2
}

// Contains currently used port when assigning them during 
// the compilations of the template leader file.
var externalPorts map[string]int
var currentExtPort int

// getPort: Is used to lookup if a port is already been 
// assigned for a particular connection if so then use 
// that port, if not then assign a new one.
func getPort(alias string) int {
	for k, v := range externalPorts {
		if k == alias {
			return v
		}
	}

	externalPorts[alias] = currentExtPort
	currentExtPort += 1
	return externalPorts[alias]
}

// parseProject: Takes a file path string and searches 
// it backwards to reveal the project name being dealt 
// with.  It then re-reverses the project name to be 
// in the correct order.
func ParseProject(path string) string {
	name := ""

	if path[len(path)-1] == os.PathSeparator {
		path = path[:len(path)-1]
	}

	i := len(path)-1
	for {
		if path[i] == os.PathSeparator { return ReverseString(name) }
		if i < 0 { return "" }

		name = name + string(path[i])
		i = i-1
	}
}

// reverseString: Reverses a string' characters and returns it.
func ReverseString(name string) string {
	reverse := ""
	
	for i := len(name)-1; i >= 0; i-- {
		reverse = reverse + string(name[i])
	}
	
	return reverse
}

// createLeader: Creates node specific leader files by building them 
// using the leader.template file.
func CreateLeader(lPath string, node config.NodeConfig, guiPort, cPath string) error {
	type tType struct {
		Node    config.NodeConfig
		GuiPort string
		ConfigPath string
	}

	tmpl := template.New("leader.template")
	tmpl.Funcs(template.FuncMap{"eq": eq, "getPort": getPort})
	tmpl, err := tmpl.ParseFiles(filepath.Join(lPath, "leader.template"))
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(lPath, node.Hostname + ".go"))
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, tType{node, guiPort, cPath})
	if err != nil {
		return err
	}

	return nil
}

// buildLeader: Runs "go build" on each node leader to get the executable.
func BuildLeader(path, hostname string) (string, error) {
	out, err := exec.Command("go", "build", "-o", filepath.Join(path, "bin", hostname), filepath.Join(path, hostname+".go")).CombinedOutput()

	return string(out), err
}

/*
 *
 * Compiles and builds the distribution leader files.
 *
 */
func Compile() {
	path := ""

	if len(os.Args) == 0 || len(os.Args) > 2 {
		log.ERROR.Println("Usage: emd compile --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd compile --help")
		os.Exit(1)
	} else if os.Args[0] == "--help" || os.Args[0] == "-h" || os.Args[0] == "help" {
		log.INFO.Println("Usage: emd compile --path <path to dir containing config.json>")
		log.INFO.Println("Usage: emd compile --help")
		os.Exit(1)
	} else if os.Args[0] == "--path" || os.Args[0] == "-p" && len(os.Args) == 2 {
		path = os.Args[1]
	} else {
		log.ERROR.Println("Usage: emd compile --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd compile --help")
		os.Exit(1)
	}

	var cfg config.Config
	currentExtPort = 40000
	externalPorts = make(map[string]int)
	config.Process(filepath.Join(path, "config.json"), &cfg)

	// Loop through all nodes in config and create
	//   leader files for each, then build them 
	//   placing them into the /leaders/bin dir.
	for _, n := range cfg.Nodes {
		err := CreateLeader(filepath.Join(path, "leaders"), n, cfg.GUI_port, filepath.Join(path, "config.json"))
		if err != nil {
			log.ERROR.Println(err)
			os.Exit(1)
		}
		
		log.INFO.Println("Leader "+n.Hostname+" compiled successfully")
		
		out, err := BuildLeader(filepath.Join(path, "leaders"), n.Hostname)
		if err != nil {
			log.ERROR.Println(out)
			log.ERROR.Println(err)
			os.Exit(1)
		}
		if out != "" { log.INFO.Println(out) }
		log.INFO.Println("Leader "+n.Hostname+" built successfully")
	}
	
	log.INFO.Println("Compile successful")
}

/*
 *
 * Cleans and removes leader files and executables.
 *
 */
func Clean() {
	path := ""

	if len(os.Args) == 0 || len(os.Args) > 2 {
		log.ERROR.Println("Usage: emd clean --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd clean --help")
		os.Exit(1)
	} else if os.Args[0] == "--help" || os.Args[0] == "-h" || os.Args[0] == "help" {
		log.INFO.Println("Usage: emd clean --path <path to dir containing config.json>")
		log.INFO.Println("Usage: emd clean --help")
		os.Exit(1)
	} else if os.Args[0] == "--path" || os.Args[0] == "-p" && len(os.Args) == 2 {
		path = os.Args[1]
	} else {
		log.ERROR.Println("Usage: emd clean --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd clean --help")
		os.Exit(1)
	}

	var cfg config.Config
	config.Process(filepath.Join(path, "config.json"), &cfg)

	for _, n := range cfg.Nodes {
		log.INFO.Println("Removing "+filepath.Join(path, "leaders", n.Hostname+".go"))
		err := os.Remove(filepath.Join(path, "leaders", n.Hostname+".go"))
		if err != nil {
			log.ERROR.Println(err)
			os.Exit(1)
		}

		log.INFO.Println("Removing "+filepath.Join(path, "leaders", "bin", n.Hostname))
		err = os.Remove(filepath.Join(path, "leaders", "bin", n.Hostname))
		if err != nil {
			log.ERROR.Println(err)
			os.Exit(1)
		}
	}
	
	log.INFO.Println("Clean successful")
}

/*
 *
 * Rsync must be installed in order to distribute the distribution.
 *
 */
func Distribute() {
	path := ""

	if len(os.Args) == 0 || len(os.Args) > 2 {
		log.ERROR.Println("Usage: emd distribute --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd distribute --help")
		os.Exit(1)
	} else if os.Args[0] == "--help" || os.Args[0] == "-h" || os.Args[0] == "help" {
		log.INFO.Println("Usage: emd distribute --path <path to dir containing config.json>")
		log.INFO.Println("Usage: emd distribute --help")
		os.Exit(1)
	} else if os.Args[0] == "--path" || os.Args[0] == "-p" && len(os.Args) == 2 {
		path = os.Args[1]
		
		// Need to make sure trailing slash is removed.
		if path[len(path)-1] == os.PathSeparator {
			path = path[:len(path)-1]
		}
	} else {
		log.ERROR.Println("Usage: emd distribute --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd distribute --help")
		os.Exit(1)
	}

	user, err := user.Current()
	if err != nil {
		log.ERROR.Println(err)
		os.Exit(1)
	}

	var cfg config.Config
	config.Process(filepath.Join(path, "config.json"), &cfg)

	for _, n := range cfg.Nodes{
		log.INFO.Println("Distributing to "+n.Hostname)
		_, err := exec.Command("rsync", "-a", "-z", path, user.Username+"@"+n.Hostname+":"+os.TempDir()).Output()
		if err != nil {
			log.ERROR.Println(err)
			os.Exit(1)
		}
	}

	log.INFO.Println("Distribute successful")
}

/*
 *
 * Start the distribution given the path to it on each node.
 *
 */
func Start() {
	path := ""

	if len(os.Args) == 0 || len(os.Args) > 2 {
		log.ERROR.Println("Usage: emd start --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd start --help")
		os.Exit(1)
	} else if os.Args[0] == "--help" || os.Args[0] == "-h" || os.Args[0] == "help" {
		log.INFO.Println("Usage: emd start --path <path to dir containing config.json>")
		log.INFO.Println("Usage: emd start --help")
		os.Exit(1)
	} else if os.Args[0] == "--path" || os.Args[0] == "-p" && len(os.Args) == 2 {
		path = os.Args[1]
		
		// Need to make sure trailing slash is removed.
		if path[len(path)-1] == os.PathSeparator {
			path = path[:len(path)-1]
		}
	} else {
		log.ERROR.Println("Usage: emd start --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd start --help")
		os.Exit(1)
	}

	user, err := user.Current()
	if err != nil {
		log.ERROR.Println(err)
		os.Exit(1)
	}

	projectName := ParseProject(path)
	if projectName == "" {
		log.ERROR.Println("Unable to parse distribution name from path.")
		log.ERROR.Println("Usage: emd start --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd start --help")
		os.Exit(1)
	}

	var cfg config.Config
	config.Process(filepath.Join(path, "config.json"), &cfg)

	passwdAnswered := false
	useSamePasswd := false
	var password []byte

	for _, n := range cfg.Nodes{
		log.INFO.Println("Starting leader on "+n.Hostname)

		if !useSamePasswd {
			fmt.Printf(user.Username+"@"+n.Hostname+"'s password: ")
			password = gopass.GetPasswd()

			if !passwdAnswered {
				passwdAnswered = true
				var ans string
			
				fmt.Printf("Is this password the same for all nodes (y/n): ")
				fmt.Scanf("%s", &ans)
				if strings.ToUpper(ans) == "Y" {
					useSamePasswd = true
				}
			}
		}

		config := &ssh.ClientConfig{
			User: user.Username,
			Auth: []ssh.AuthMethod{
				ssh.Password(string(password)),
			},
		}
		client, err := ssh.Dial("tcp", n.Hostname+":22", config)
		if err != nil {
			log.ERROR.Println(err)
			os.Exit(1)
		}
		
		session, err := client.NewSession()
		if err != nil {
			log.ERROR.Println(err)
			os.Exit(1)
		}
		defer session.Close()
	
		var b bytes.Buffer
		session.Stdout = &b
		
		var cmd string
		if runtime.GOOS == "windows" {
			cmd = "start /B"+filepath.Join(os.TempDir(), projectName, "leaders", "bin", n.Hostname)+" > "+os.DevNull
		} else {
			cmd = "nohup "+filepath.Join(os.TempDir(), projectName, "leaders", "bin", n.Hostname)+" > "+os.DevNull+" 2>&1 &"
		}

		if err := session.Run(cmd); err != nil {
			log.ERROR.Println(err)
			os.Exit(1)
		}
	}

	log.INFO.Println("Start successful")
}

/*
 *
 * Perform GET request to stop distribution.
 *
 */
func Stop() {
	path := ""

	if len(os.Args) == 0 || len(os.Args) > 2 {
		log.ERROR.Println("Usage: emd stop --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd stop --help")
		os.Exit(1)
	} else if os.Args[0] == "--help" || os.Args[0] == "-h" || os.Args[0] == "help" {
		log.INFO.Println("Usage: emd stop --path <path to dir containing config.json>")
		log.INFO.Println("Usage: emd stop --help")
		os.Exit(1)
	} else if os.Args[0] == "--path" || os.Args[0] == "-p" && len(os.Args) == 2 {
		path = os.Args[1]
		
		// Need to make sure trailing slash is removed.
		if path[len(path)-1] == os.PathSeparator {
			path = path[:len(path)-1]
		}
	} else {
		log.ERROR.Println("Usage: emd stop --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd stop --help")
		os.Exit(1)
	}
	
	var cfg config.Config
	config.Process(filepath.Join(path, "config.json"), &cfg)
	
	log.INFO.Println("Stopping distribution")
	
	for _, n := range cfg.Nodes {
		log.INFO.Println("Stopping node "+n.Hostname)

		// Stop all the workers
		_, err := http.Get("http://"+n.Hostname+":"+cfg.GUI_port+"/stop")
		if err != nil {
			log.ERROR.Println(err)
			os.Exit(1)
		}

		// Stop the leader
		_, err = http.Get("http://"+n.Hostname+":"+cfg.GUI_port+"/stop")
		if err != nil {
			if strings.Contains(err.Error(), "EOF") { continue }
			log.ERROR.Println(err)
			os.Exit(1)
		}
	}
	
	log.INFO.Println("Stop successful")
}

/*
 *
 * Perform GET request to get status of distribution.
 *
 */
func Status() {
	path := ""

	if len(os.Args) == 0 || len(os.Args) > 2 {
		log.ERROR.Println("Usage: emd status --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd status --help")
		os.Exit(1)
	} else if os.Args[0] == "--help" || os.Args[0] == "-h" || os.Args[0] == "help" {
		log.INFO.Println("Usage: emd status --path <path to dir containing config.json>")
		log.INFO.Println("Usage: emd status --help")
		os.Exit(1)
	} else if os.Args[0] == "--path" || os.Args[0] == "-p" && len(os.Args) == 2 {
		path = os.Args[1]
		
		// Need to make sure trailing slash is removed.
		if path[len(path)-1] == os.PathSeparator {
			path = path[:len(path)-1]
		}
	} else {
		log.ERROR.Println("Usage: emd status --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd status --help")
		os.Exit(1)
	}
	
	var cfg config.Config
	config.Process(filepath.Join(path, "config.json"), &cfg)
	
	for _, n := range cfg.Nodes {
		log.INFO.Println("Obtaining status of node "+n.Hostname)

		resp, err := http.Get("http://"+n.Hostname+":"+cfg.GUI_port+"/status")
		if err != nil {
			log.ERROR.Println(err)
			os.Exit(1)
		}
		
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.ERROR.Println(err)
			os.Exit(1)
		}
		log.INFO.Println(string(content))
		resp.Body.Close()
	}

	log.INFO.Println("Status successful")
}

/*
 *
 * Perform GET request to get metrics of distribution.
 *
 */
func Metrics() {
	path := ""

	if len(os.Args) == 0 || len(os.Args) > 2 {
		log.ERROR.Println("Usage: emd metrics --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd metrics --help")
		os.Exit(1)
	} else if os.Args[0] == "--help" || os.Args[0] == "-h" || os.Args[0] == "help" {
		log.INFO.Println("Usage: emd metrics --path <path to dir containing config.json>")
		log.INFO.Println("Usage: emd metrics --help")
		os.Exit(1)
	} else if os.Args[0] == "--path" || os.Args[0] == "-p" && len(os.Args) == 2 {
		path = os.Args[1]
		
		// Need to make sure trailing slash is removed.
		if path[len(path)-1] == os.PathSeparator {
			path = path[:len(path)-1]
		}
	} else {
		log.ERROR.Println("Usage: emd metrics --path <path to dir containing config.json>")
		log.ERROR.Println("Usage: emd metrics --help")
		os.Exit(1)
	}
	
	var cfg config.Config
	config.Process(filepath.Join(path, "config.json"), &cfg)
	
	for _, n := range cfg.Nodes {
		log.INFO.Println("Obtaining metrics of node "+n.Hostname)

		resp, err := http.Get("http://"+n.Hostname+":"+cfg.GUI_port+"/metrics")
		if err != nil {
			log.ERROR.Println(err)
			os.Exit(1)
		}

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.ERROR.Println(err)
			os.Exit(1)
		}
		log.INFO.Println(string(content))
		resp.Body.Close()
	}

	log.INFO.Println("Metrics successful")
}

/*
 *
 * Create directory structure for new distribution (copy boilerplate 
 *    distribution contents into projectName directory.
 *
 */
func NewProject() {
	log.INFO.Println("Creating new distribution")

	// Copy boilerplate/* stuff into the path/name given.
	_, err := exec.Command("git", "clone", "https://github.com/go-emd/boilerplate.git").Output()
	if err != nil {
		log.ERROR.Println(err)
		os.Exit(1)
	}

	log.INFO.Println("New successful")
}

func main() {
	log.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	// Get the command line arguments
	var action string
	if len(os.Args) >= 2 {
		action = os.Args[1]
	} else {
		log.ERROR.Println("Usage: emd <action> args {new|compile|clean|distribute|start|stop|status|metrics}")
		os.Exit(1)
	}

	os.Args = os.Args[2:]

	switch action {
	case "new":
		NewProject()
	case "compile":
		log.INFO.Println("Performing: compile")
		Compile()
	case "clean":
		log.INFO.Println("Performing: clean")
		Clean()
	case "distribute":
		log.INFO.Println("Performing: distribute")
		Distribute()
	case "start":
		log.INFO.Println("Performing: start")
		Start()
	case "stop":
		log.INFO.Println("Performing: stop")
		Stop()
	case "status":
		log.INFO.Println("Performing: status")
		Status()
	case "metrics":
		log.INFO.Println("Performing: metrics")
		Metrics()
	default:
		log.ERROR.Println("Invalid action.")
		log.ERROR.Println("Usage: emd <action> args {new|compile|clean|distribute|start|stop|status|metrics}")
		os.Exit(1)
	}
}
