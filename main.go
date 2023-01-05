package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-clix/cli"
	neat "github.com/itaysk/kubectl-neat/cmd"
)

func main() {
	log.SetFlags(0)

	cmd := cli.Command{
		Use:   "kubectl-neat-diff [file1] [file2]",
		Short: "Remove fields from kubectl diff that carry low / no information",
		Args:  cli.ArgsExact(2),
	}

	cmd.Run = func(cmd *cli.Command, args []string) error {
		if err := neatifyDir(args[0]); err != nil {
			return err
		}
		if err := neatifyDir(args[1]); err != nil {
			return err
		}

		c := exec.Command("diff", "-uN", args[0], args[1])
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	}

	err := cmd.Execute()
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			// exit status 1 means there is a diff, so we ignore this
			if e.ExitCode() == 1 {
				return
			}
		}
		// otherwise log all errors
		log.Fatalln("Error:", err)
	}
}

func neatifyDir(dir string) error {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		filename := filepath.Join(dir, fi.Name())
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		n, err := neat.NeatYAMLOrJSON(data, "same")
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(filename, []byte(n), fi.Mode()); err != nil {
			return err
		}
	}

	return nil
}
