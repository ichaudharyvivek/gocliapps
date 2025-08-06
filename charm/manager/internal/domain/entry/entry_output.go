package domain

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const divider = "---------"

func FormattedEntires(entries []*Entry) []byte {
	var op string
	for i := len(entries) - 1; i >= 0; i-- {
		op += fmt.Sprintf("ID: %d\nCreated: %s\nMessage:\n\n %s\n %s\n", entries[i].ID, entries[i].CreatedAt.Format(time.RFC1123), entries[i].Message, divider)
	}

	return []byte(op)
}

func FormattedEntry(entry *Entry) string {
	return fmt.Sprintf("**ID:** %d\n\n**Created:** %s\n\n**Message:**\n\n%s\n\n%s\n", entry.ID, entry.CreatedAt.Format(time.RFC1123), entry.Message, divider)
}

func OutputEntriesToMarkdown(entries []*Entry) error {
	f, err := os.OpenFile("./output.md", os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("cannot open file %w", err)
	}
	defer f.Close()

	op := FormattedEntires(entries)

	_, err = f.Write(op)
	if err != nil {
		return fmt.Errorf("cannot save file %w", err)
	}

	return err
}

func OutputEntriesToPDF(entries []*Entry) error {
	op := FormattedEntires(entries)

	pandoc := exec.Command("pandoc", "-s", "-o", "output.pdf")
	wc, wcerr := pandoc.StdinPipe() // io.WriteCloser, err
	if wcerr != nil {
		return fmt.Errorf("cannot stdin to pandoc: %w", wcerr)
	}
	goerr := make(chan error)
	done := make(chan bool)
	go func() {
		var err error
		defer func() {
			err = wc.Close()
		}()
		_, err = wc.Write(op)
		goerr <- err
		close(goerr)
		close(done)
	}()
	if err := <-goerr; err != nil {
		return fmt.Errorf("cannot write file to pandoc: %w", err)
	}
	err := pandoc.Run()
	if err != nil {
		return fmt.Errorf("cannot run pandoc: %w", err)
	}
	return nil
}
