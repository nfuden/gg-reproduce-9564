//go:build mage

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
)

// var Default = FullInstall

// Readme prints out a quick readme
func Readme() {
	fmt.Println("This project is run on mage. To see the available commands, run `mage -l`")
	fmt.Println("To run the full install, run `mage fullinstall`")
	fmt.Println("To run some requests and see the output run `mage demo:run2000`")
	fmt.Println("Each resource folder can be applied either via kubectl or `mage resources:apply<folder_name>`")
}

// FullInstall is the default install for this repo
// Run by default
func FullInstall() {
	clean := AmIClean()
	if !clean {
		fmt.Println("Not clean")
		return

	}
	InstallGloo(false)
	Resources{}.Apply()
}

func smashEnv() (license string, hasBoth bool) {
	license, hasLic := os.LookupEnv("GLOO_LICENSE_KEY")
	otherLic, hasOtherLic := os.LookupEnv("GLOO_EDGE_LICENSE_KEY")
	if hasLic && hasOtherLic {
		return "", true
	}
	if hasOtherLic {
		return otherLic, false
	}
	return license, false
}

func InstallGlooPOC() {
	InstallGloo(true)
}
func InstallGlooPublished() {
	InstallGloo(false)
}

// InstallGloo will install edge with given values
func InstallGloo(shouldPOC bool) {

	toInstall := "install.sh"
	if shouldPOC {
		toInstall = "install_poc.sh"
	}

	err := simpleRun(func(cmd *exec.Cmd) {

		cmd.Stderr = os.Stdout
		cmd.Stdout = os.Stdout
		cmd.Env = os.Environ()
		newLic, hasBoth := smashEnv()
		if hasBoth {
			fmt.Println("Both GLOO_LICENSE_KEY and GLOO_EDGE_LICENSE_KEY are set. Please only set one")
		} else {
			cmd.Env = append(cmd.Env, fmt.Sprintf("GLOO_EDGE_LICENSE_KEY=%s", newLic))
		}
	}, filepath.Join(".", "install", toInstall))
	if err != nil {
		fmt.Println(err)
	}
}

type Resources mg.Namespace

// Apply applies all sub resources
func (r Resources) Apply() {
	fmt.Println("Applying resources")
	r.ApplyServices()
	r.ApplyUpstreams()
	r.ApplyGateways()
	r.ApplyRouteTables()
	r.ApplyVirtualServices()

}

// ApplyServices applies services
func (Resources) ApplyServices() {
	fmt.Println("Applying services")
	applyResource("services")
}

// ApplyGateways applies gateways
func (Resources) ApplyGateways() {
	fmt.Println("Applying services")
	applyResource("gateways")

}

// ApplyUpstreams applies upstreams
func (Resources) ApplyUpstreams() {
	fmt.Println("Applying upstreams")
	applyResource("upstreams")
}

func (Resources) ApplyRouteTables() {
	fmt.Println("Applying route tables")
	applyResource("routetables")
}

func (Resources) ApplyVirtualServices() {
	fmt.Println("Applying virtual services")
	applyResource("virtualservices")
}

func applyResource(toApply string) error {
	files, err := os.ReadDir(filepath.Join(".", toApply))
	if err != nil {
		fmt.Println("failed reading directory:", err)
	}

	for _, file := range files {
		fmt.Printf("applying %s via file: %s\n", toApply, file.Name())
		err = simpleRun(func(cmd *exec.Cmd) {}, "kubectl", "apply", "-f", filepath.Join(".", toApply, file.Name()))
		if err != nil {
			fmt.Printf("failed to apply: %s, %v\n", file.Name(), err)
		}

	}
	return nil
}

// TODO: AmIClean will eventually check to see if other resources exist that might clash
func AmIClean() (cleaninstall bool) {
	fmt.Println("TODO: we need to determine what it means to be clean")
	return true
}

func simpleRun(extraFunc func(*exec.Cmd), toRun ...string) error {
	cmd := exec.CommandContext(context.TODO(), toRun[0], toRun[1:]...)

	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout

	extraFunc(cmd)

	err := cmd.Run()
	return err

}

type Demo mg.Namespace

func (Demo) Run2000() {
	err := simpleRun(func(cmd *exec.Cmd) {}, strings.Split("hey -n 2000 -c 45 -H test:true http://localhost:8080/response-headers", " ")...)
	if err != nil {
		fmt.Println(err)
		return
	}

	// err = simpleRun(func(cmd *exec.Cmd) {}, "curl", "http://localhost:19000")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}
