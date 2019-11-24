package main

//go-generate:go-bindata -fs -prefix "static/" static/

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	ScriptLocation string
	CmdArgs        []string
	ExecFname      string
)

func fillRunningCmd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	CmdArgs = []string{ScriptLocation, "-target", RawifyString(r.FormValue("barcodeAlign")), "-barcode", RawifyString(r.FormValue("barcodeContents")), "-caption", RawifyString(r.FormValue("captionContents")), "-subcaption", "''"}
	http.Redirect(w, r, "/log.html", http.StatusSeeOther)
	fmt.Println(CmdArgs)
}

func main() {
	var port int

	flag.StringVar(&ScriptLocation, "script", "", "Specify printing python script, if any")
	flag.IntVar(&port, "port", 4080, "Port, to run on")
	flag.Parse()
	ExecFname = "/usr/bin/python"
	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	http.Handle("/", http.FileServer(AssetFile()))
	http.HandleFunc("/print", fillRunningCmd)
	http.HandleFunc("/ws", ServeWs)
	server.ListenAndServe()
}
