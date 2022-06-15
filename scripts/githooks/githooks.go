package main

import (
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

var scriptTemplate = `#!/bin/bash
echo "Formatting go imports..."
files=$(git diff --name-only --cached --diff-filter=ARM | grep "\.go$" | xargs gosimports -l -w -local "github.com/cszczepaniak/git-tools")
if [ -z "$files" ] 
then
	echo "No changes!"
else
	echo "$files" | xargs git add
fi
`

func main() {
	_, err := os.Stat(`.git`)
	if os.IsNotExist(err) {
		log.Fatal(`.git directory not found. Are you in the root of the repo?`)
	} else if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(`.git/hooks/pre-commit`)
	if !os.IsNotExist(err) {
		if err == nil {
			var overwrite bool
			serr := survey.AskOne(&survey.Confirm{
				Message: `Pre-commit hook already exists. Would you like to overwrite it?`,
			}, &overwrite)
			if serr != nil {
				log.Fatal(err)
			}
			if !overwrite {
				return
			}
		} else {
			log.Fatal(err)
		}
	}

	err = os.WriteFile(`.git/hooks/pre-commit`, []byte(scriptTemplate), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
