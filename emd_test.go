package main

import (
	"io/ioutil"
	"os/exec"
	"testing"
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
	files, err := ioutil.ReadDir("boilerplate")
	if err != nil {
		t.Fatal(err)
	} else {
		for i, file := range files {
			if file.Name() != answers[i] {
				t.Fatal(file.Name() + "!=" + answers[i])
			}
		}
	}
}

func TestCompile(t *testing.T) {
	answer := "MD5 (boilerplate/leaders/bin/localhost) = dc287bfea93976a1d00b4e96c133e3b4\n"

	output, err := exec.Command("go", "run", "emd.go", "compile", "--path", "boilerplate").CombinedOutput()
	if err != nil {
		t.Log(string(output))
		t.Fatal(err)
	} else {
		output, _ = exec.Command("md5", "boilerplate/leaders/bin/localhost").Output()
		if string(output) != answer {
			t.Fatal(string(output) + "!=" + answer)
		}
	}
}

func TestDistribute(t *testing.T) {
	_, err := exec.Command("go", "run", "emd.go", "distribute", "--path", "boilerplate").CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStart(t *testing.T) {
	_, err := exec.Command("go", "run", "emd.go", "start", "--path", "boilerplate").CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStatus(t *testing.T) {
	_, err := exec.Command("go", "run", "emd.go", "status", "--path", "boilerplate").CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMetrics(t *testing.T) {
	_, err := exec.Command("go", "run", "emd.go", "metrics", "--path", "boilerplate").CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStop(t *testing.T) {
	_, err := exec.Command("go", "run", "emd.go", "stop", "--path", "boilerplate").CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClean(t *testing.T) {
	_, err := exec.Command("go", "run", "emd.go", "clean", "--path", "boilerplate").CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	_, _ = exec.Command("rm", "-rf", "boilerplate").Output()
}
