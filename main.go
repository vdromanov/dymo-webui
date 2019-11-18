package main

import (
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
	ScriptLocation = "labels.py"
	server := http.Server{
		Addr: ":4080",
	}
	http.Handle("/", http.FileServer(AssetFile()))
	http.HandleFunc("/print", runPrinting)
	server.ListenAndServe()
}
