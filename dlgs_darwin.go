package dialog

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/lucidoc/dialog/cocoa"
)

func (b *MsgBuilder) yesNo() bool {
	expression := fmt.Sprintf("tell app \"System Events\" to display dialog \"%s\" buttons {\"No\", \"Yes\"} default button 2 with title \"%s\"", b.Msg, b.Dlg.Title)
	cmd := exec.Command("osascript", "-e", expression)

	var out bytes.Buffer
	cmd.Stdout = &out

	cmd.Start()

	if err := cmd.Wait(); err != nil {
		panic(err)
	} else {
		output := strings.Trim(out.String(), "\n\r ")
		output = strings.ToLower(output)
		regex := regexp.MustCompile("(no|yes)")
		output = regex.FindString(output)
		return output == "yes"
	}
}

func (b *MsgBuilder) info() {
	expression := fmt.Sprintf("tell app \"System Events\" to display dialog \"%s\" buttons {\"Ok\"} default button 1 with title \"%s\"", b.Msg, b.Dlg.Title)
	cmd := exec.Command("osascript", "-e", expression)
	cmd.Start()

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}

func (b *MsgBuilder) error() {
	expression := fmt.Sprintf("tell app \"System Events\" to display dialog \"Error:\n\n%s\" buttons {\"Ok\"} default button 1 with title \"%s\"", b.Msg, b.Dlg.Title)
	cmd := exec.Command("osascript", "-e", expression)
	cmd.Start()

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}

func (b *FileBuilder) load() (string, error) {
	return b.run(false)
}

func (b *FileBuilder) save() (string, error) {
	return b.run(true)
}

func (b *FileBuilder) run(save bool) (string, error) {
	star := false
	var exts []string
	for _, filt := range b.Filters {
		for _, ext := range filt.Extensions {
			if ext == "*" {
				star = true
			} else {
				exts = append(exts, ext)
			}
		}
	}
	if star && save {
		/* OSX doesn't allow the user to switch visible file types/extensions. Also
		** NSSavePanel's allowsOtherFileTypes property has no effect for an open
		** dialog, so if "*" is a possible extension we must always show all files. */
		exts = nil
	}
	f, err := cocoa.FileDlg(save, b.Dlg.Title, exts, star, b.StartDir, b.StartFile)
	if f == "" && err == nil {
		return "", ErrCancelled
	}
	return f, err
}

func (b *DirectoryBuilder) browse() (string, error) {
	f, err := cocoa.DirDlg(b.Dlg.Title, b.StartDir)
	if f == "" && err == nil {
		return "", ErrCancelled
	}
	return f, err
}
