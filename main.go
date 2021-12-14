package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-clix/cli"
	neat "github.com/zzehring/kubectl-neat/v2/cmd"
)

func main() {
	log.SetFlags(0)

	cmd := cli.Command{
		Use:   "kubectl-neat-diff [file1] [file2]",
		Short: "Remove fields from kubectl diff that carry low / no information",
		Args:  cli.ArgsExact(2),
	}

	removeLinesRegEx := cmd.Flags().StringP("remove-matching-lines", "R", "",
		"Remove lines matching RegEx from 'kubectl get' outputs prior to diff.")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		if err := neatifyDir(args[0], *removeLinesRegEx); err != nil {
			return err
		}
		if err := neatifyDir(args[1], *removeLinesRegEx); err != nil {
			return err
		}

		c := exec.Command("diff", "-uN", args[0], args[1])
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	}

	if err := cmd.Execute(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func neatifyDir(dir, removeLinesRe string) error {
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

		if removeLinesRe != "" {
			regExp, err := regexp.Compile(removeLinesRe)
			if err != nil {
				return err
			}
			n = []byte(removeMatchingLines(string(n), regExp))
		}
		if err := ioutil.WriteFile(filename, []byte(n), fi.Mode()); err != nil {
			return err
		}
	}

	return nil
}

func removeMatchingLines(source string, re *regexp.Regexp) string {
	lines := strings.Split(source, "\n")
	var nonMatchingLines []string
	for _, line := range lines {
		if !re.MatchString(line) {
			nonMatchingLines = append(nonMatchingLines, line)
		}
	}
	return strings.Join(nonMatchingLines, "\n")
}
