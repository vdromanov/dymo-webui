package main

import (
	"fmt"
	"os/exec"
)

var ScriptLocation string

func rawifyString(input string) (output string) {
	line := fmt.Sprintf("%q", input)
	return line[1 : len(line)-1]
}

func PrintLabel(barcodeAlign, barcodeContents, captionContents string) (ret string, err error) {
	cmd := exec.Command("python", ScriptLocation, "-target", barcodeAlign, "-barcode", rawifyString(barcodeContents), "-caption", rawifyString(captionContents), "-subcaption", "''")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return (fmt.Sprint(err) + ": " + string(output)), err
	}
	return string(output), nil
}
