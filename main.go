package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Microsoft/go-winio"
)

func runWithPrivileges() ([]string, error) {
	privs, err := winio.GetCurrentThreadPrivileges()
	if err != nil {
		return nil, fmt.Errorf("failed to get privileges: %w", err)
	}
	return privs, nil
}

func main() {
	privs, err := winio.GetEnabledPrivileges()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Process privileges in main before enable:\t\t%s\n", strings.Join(privs, ", "))

	if err := winio.EnableProcessPrivileges([]string{winio.SeSecurityPrivilege}); err != nil {
		log.Fatalf("failed to enable SeSecurityPrivilege: %s", err)
	}
	defer winio.DisableProcessPrivileges([]string{winio.SeSecurityPrivilege})

	var runPrivs []string

	if err := winio.RunWithPrivileges([]string{winio.SeBackupPrivilege}, func() error {
		var err error
		runPrivs, err = runWithPrivileges()
		return err
	}); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("RunWithPrivileges privileges:\t\t\t\t%s\n", strings.Join(runPrivs, ", "))

	privs, err = winio.GetEnabledPrivileges()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Process privileges in main after RunWithPrivileges:\t%s\n", strings.Join(privs, ", "))
}
