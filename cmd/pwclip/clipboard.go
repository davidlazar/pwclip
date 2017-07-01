package main

import (
	"bytes"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func getClipboard() ([]byte, error) {
	return exec.Command(clipboardGetCmd[0], clipboardGetCmd[1:]...).Output()
}

func setClipboard(data []byte) error {
	cmd := exec.Command(clipboardSetCmd[0], clipboardSetCmd[1:]...)
	cmd.Stdin = bytes.NewReader(data)
	return cmd.Run()
}

func setClipboardTemporarily(data []byte, d time.Duration) error {
	prev, err := getClipboard()
	if err != nil {
		return err
	}

	sigchan := make(chan os.Signal, 1)
	go func() {
		for _ = range sigchan {
			setClipboard(prev)
			os.Exit(0)
		}
	}()
	signal.Notify(sigchan, os.Interrupt, os.Kill)

	if err := setClipboard(data); err != nil {
		return err
	}
	time.Sleep(d)
	if err := setClipboard(prev); err != nil {
		return err
	}
	return nil
}
