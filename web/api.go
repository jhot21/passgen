package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	passgen "codeberg.org/jhot/passgen"
)

type generateRequest struct {
	Passphrase      string `json:"passphrase"`
	Salt            string `json:"salt"`
	Length          int    `json:"length"`
	Version         int    `json:"version"`
	NoUppers        bool   `json:"no_uppers"`
	NoNumbers       bool   `json:"no_numbers"`
	NoSpecials      bool   `json:"no_specials"`
	CustomSpecials  string `json:"custom_specials"`
	NoAmbiguous     bool   `json:"no_ambiguous"`
	CustomAmbiguous string `json:"custom_ambiguous"`
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	var req generateRequest

	switch r.Method {
	case http.MethodGet:
		q := r.URL.Query()
		req.Passphrase = q.Get("passphrase")
		req.Salt = q.Get("salt")
		if l := q.Get("length"); l != "" {
			req.Length, _ = strconv.Atoi(l)
		}
		if v := q.Get("version"); v != "" {
			req.Version, _ = strconv.Atoi(v)
		}
		req.NoUppers = ParseBool(q.Get("no_uppers"))
		req.NoNumbers = ParseBool(q.Get("no_numbers"))
		req.NoSpecials = ParseBool(q.Get("no_specials"))
		req.CustomSpecials = q.Get("custom_specials")
		req.NoAmbiguous = ParseBool(q.Get("no_ambiguous"))
		req.CustomAmbiguous = q.Get("custom_ambiguous")

	case http.MethodPost:
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid JSON body")
			return
		}
	}

	if req.Passphrase == "" {
		respondError(w, r, http.StatusBadRequest, "passphrase is required")
		return
	}

	b := passgen.New(req.Passphrase)
	if req.Salt != "" {
		b = b.WithSalt(req.Salt)
	}
	if req.Length > 0 {
		b = b.WithLength(req.Length)
	}
	if req.Version > 0 {
		b = b.WithVersion(req.Version)
	}
	if req.NoUppers {
		b = b.NoUppers()
	}
	if req.NoNumbers {
		b = b.NoNumbers()
	}
	if req.NoSpecials {
		b = b.NoSpecials()
	}
	if req.CustomSpecials != "" {
		b = b.WithSpecials(req.CustomSpecials)
	}
	if req.CustomAmbiguous != "" {
		b = b.WithAmbiguousChars(req.CustomAmbiguous)
	} else if req.NoAmbiguous {
		b = b.NoAmbiguous()
	}

	pass, err := b.Generate()
	if err != nil {
		respondError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"password": pass})
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(pass))
}

// ParseBool returns true for "true", "yes", "1", "on" (case-insensitive).
func ParseBool(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "yes", "1", "on":
		return true
	}
	return false
}

func wantsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "application/json")
}

func respondError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(map[string]string{"error": msg})
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(msg))
}
