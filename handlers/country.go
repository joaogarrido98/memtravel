package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"memtravel/cache"
	"memtravel/db"
	"memtravel/language"
	"memtravel/middleware"
)

type Country struct {
	ID     int    `json:"id"`
	ISO    string `json:"iso"`
	Name   string `json:"name,omitempty"`
	NamePT string `json:"namept,omitempty"`
	NameFR string `json:"namefr,omitempty"`
	NameES string `json:"namees,omitempty"`
}

var countryCache = cache.NewSimpleCache()

func (handler *Handler) GetAllCountries(w http.ResponseWriter, r *http.Request) {
	var deferredErr error
	defer func() {
		if deferredErr != nil {
			log.Printf("Error: [%s], context_id: [%s], user_id: [%s]",
				deferredErr.Error(),
				r.Context().Value(middleware.RequestContextID),
				r.Context().Value(middleware.AuthUserID),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	languageID := r.URL.Query().Get(languageParamID)
	if !language.SupportedLanguage(languageID) {
		deferredErr = errorLanguageID
		return
	}

	if cachedCountries, ok := countryCache.Get(languageID); ok {
		writeServerResponse(w, true, cachedCountries)
		return
	}

	var languageRow string
	switch languageID {
	case language.PortugueseID:
		languageRow = "namept"
	case language.FrenchID:
		languageRow = "namefr"
	case language.SpanishID:
		languageRow = "namees"
	case language.EnglishID:
		languageRow = "name"
	default:
		deferredErr = fmt.Errorf("%s not a valid language id", languageID)
		return
	}

	var countries []Country

	rows, deferredErr := handler.database.Query(fmt.Sprintf(db.GetAllCountries, languageRow))
	if deferredErr != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var country Country

		switch languageRow {
		case "namept":
			deferredErr = rows.Scan(&country.ID, &country.ISO, &country.NamePT)
		case "namefr":
			deferredErr = rows.Scan(&country.ID, &country.ISO, &country.NameFR)
		case "namees":
			deferredErr = rows.Scan(&country.ID, &country.ISO, &country.NameES)
		default:
			deferredErr = rows.Scan(&country.ID, &country.ISO, &country.Name)
		}

		if deferredErr != nil {
			return
		}

		countries = append(countries, country)
	}

	deferredErr = rows.Err()
	if deferredErr != nil {
		return
	}

	countryCache.Set(languageID, countries, 2*time.Hour)

	deferredErr = writeServerResponse(w, true, countries)
}
