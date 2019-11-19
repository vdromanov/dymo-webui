package main

import (
	"flag"
	"fmt"
	"net/http"
)

func runPrinting(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	out, err := PrintLabel(r.FormValue("barcodeAlign"), r.FormValue("captionContents"), r.FormValue("barcodeContents"))
	if err != nil {
		http.Error(w, out, http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func main() {
	var port int

	flag.StringVar(&ScriptLocation, "exec", "labels.py", "Specify printing python script")
	flag.IntVar(&port, "port", 4080, "Port, to run on")
	flag.Parse()

	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	http.Handle("/", http.FileServer(AssetFile()))
	http.HandleFunc("/print", runPrinting)
	server.ListenAndServe()
}
