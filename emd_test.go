package main

import (
	"testing"
	//"os"
	"os/exec"
	"io/ioutil"
)

func TestNewProject(t *testing.T) {
	answers := []string{
		".git",
		"README.md",
		"config.json",
		"leaders",
		"workers",
	}

	exec.Command("go", "run", "emd.go", "new").CombinedOutput()
	files, err := ioutil.ReadDir("./boilerplate")
	if err != nil {
		t.Fatal(err)
	} else {
		for i, file := range files {
			if file.Name() != answers[i] {
				t.Log(file.Name() + "!=" + answers[i])
				t.Fail()
			}
		}
	}
}

func TestCompile(t *testing.T){
	answer := "MD5 (boilerplate/leaders/bin/localhost) = dc287bfea93976a1d00b4e96c133e3b4\n"

	output, err := exec.Command("go", "run", "emd.go", "compile", "--path", "boilerplate").CombinedOutput()
	if err != nil {
		t.Log(string(output))
		t.Fatal(err)
		t.Fail()
	} else {
		output, _ = exec.Command("md5", "boilerplate/leaders/bin/localhost").Output()
		if string(output) != answer {
			t.Log(string(output) + "!=" + answer)
			t.Fail()
		}
	}
}

func TestDistribute(t *testing.T){
	_, err := exec.Command("go", "run", "emd.go", "distribute", "--path", "boilerplate").CombinedOutput()
	if err != nil {
		t.Fatal(err)
		t.Fail()
	} else {
		// TODO
	}
}

func TestStart(t *testing.T){
	_, err := exec.Command("go", "run", "emd.go", "start", "--path", "boilerplate").CombinedOutput()
	if err != nil {
		t.Fatal(err)
		t.Fail()
	} else {
		// TODO
	}
}

func TestStatus(t *testing.T){
	_, err := exec.Command("go", "run", "emd.go", "status", "--path", "boilerplate").CombinedOutput()
	if err != nil {
		t.Fatal(err)
		t.Fail()
	} else {
		// TODO
	}
}

func TestMetrics(t *testing.T){
	_, err := exec.Command("go", "run", "emd.go", "metrics", "--path", "boilerplate").CombinedOutput()
	if err != nil {
		t.Fatal(err)
		t.Fail()
	} else {
		// TODO
	}
}

func TestStop(t *testing.T){
	_, err := exec.Command("go", "run", "emd.go", "stop", "--path", "boilerplate").CombinedOutput()
	if err != nil {
		t.Fatal(err)
		t.Fail()
	} else {
		// TODO
	}
}

func TestClean(t *testing.T){
	_, err := exec.Command("go", "run", "emd.go", "clean", "--path", "boilerplate").CombinedOutput()
	if err != nil {
		t.Fatal(err)
		t.Fail()
	} else {
		// TODO
	}

	_, _ = exec.Command("rm", "-rf", "boilerplate").Output()
}
