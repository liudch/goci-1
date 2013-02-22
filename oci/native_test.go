package oci_test

import (
	"fmt"
	"github.com/egravert/goci/oci"
	"testing"
)

func TestCreateEnvironment(t *testing.T) {
	_, err := oci.CreateEnvironment()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateAndFreeErrorHandle(t *testing.T) {
	env, _ := oci.CreateEnvironment()
	errHandle, err := oci.AllocateErrorHandle(env)
	if err != nil {
		t.Fatal(err)
	}

	err = oci.FreeErrorHandle(errHandle)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSuccessfulBasicLogin(t *testing.T) {
	env, _ := oci.CreateEnvironment()
	_, err := oci.BasicLogin(env, "hr", "oracle", "192.168.69.131/ORCL")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSFailedBasicLogin(t *testing.T) {
	env, _ := oci.CreateEnvironment()
	_, err := oci.BasicLogin(env, "boom", "fail", "192.168.69.131/ORCL")
	if err == nil {
		t.Fatal("Expected error, but received nil")
	}
}

func ExampleBasicLogin() {
	var env oci.EnvHandle
	var err error

	if env, err = oci.CreateEnvironment(); err != nil {
		fmt.Println(err)
		return
	}

	if _, err := oci.BasicLogin(env, "scott", "tiger", "127.0.0.1/ORCL"); err != nil {
		fmt.Println(err)
		return
	}
}

func ExamplePing() {
	env, _ := oci.CreateEnvironment()
	svr, err := oci.BasicLogin(env, "hr", "oracle", "192.168.69.131/ORCL")
	if err != nil {
		fmt.Println(err)
	}

	err = oci.Ping(env, svr)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Success!")
	}
  // Output: Success!
}
