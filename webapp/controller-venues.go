package main

import (
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/eventhunt-org/webapp/framework"
	"github.com/eventhunt-org/webapp/webapp/db"
)

/*
 * Page to add a venue to an event.
 */
func (a *app) venueNew(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "login")

	// middlewareLIO ensures we have a user
	u := r.Context().Value("user").(*db.User)
	// middlewareEvent ensures we have an Event
	e := r.Context().Value("event").(*db.Event)

	// Only the owner of the group can add/edit an event
	if !e.TheGroup.HasCreate(u.ID) {

		session.AddFlash(framework.Flash{
			framework.FlashFail,
			"You don't have permission to add a venue to this event.",
		})

		session.Save(r, w)
		http.Redirect(w, r, "/events/"+e.IDString(), http.StatusFound)
		return
	}

	cities, err := db.GetCitiesByAll(a.DB)
	if err != nil {

		slog.Error("Failed to get all cities.", "err", err)
		session.AddFlash(framework.Flash{
			framework.FlashFail,
			"Cities list failed to load.",
		})

		session.Save(r, w)
		http.Redirect(w, r, "/events/"+e.IDString(), http.StatusFound)
		return
	}

	renderPage(a, "venues/new", w, r, map[string]interface{}{
		"User":   u,
		"Event":  e,
		"Cities": cities,
	})
}

/*
 * Processing of the new venue event page.
 */
func (a *app) venueNewPost(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "login")

	// middlewareLIO ensures we have a user
	u := r.Context().Value("user").(*db.User)
	// middlewareEvent ensures we have an Event
	e := r.Context().Value("event").(*db.Event)

	// Only the owner of the group can add/edit an event
	if !e.TheGroup.HasCreate(u.ID) {

		session.AddFlash(framework.Flash{
			framework.FlashFail,
			"You don't have permission to add a venue to this event.",
		})

		session.Save(r, w)
		http.Redirect(w, r, "/events/"+e.IDString(), http.StatusFound)
		return
	}

	r.ParseForm()
	defer r.Body.Close()

	name := r.Form.Get("venue-name")
	address := r.Form.Get("street-address")
	city := r.Form.Get("city")

	cityID, err := strconv.ParseUint(city, 10, 64)
	if err != nil {

		slog.Error("City ID is not valid.", "id", city)
		session.AddFlash(framework.Flash{
			framework.FlashFail,
			"City was invalid.",
		})

		session.Save(r, w)
		http.Redirect(w, r, "/events/"+e.IDString(), http.StatusFound)
		return
	}

	v, err := db.NewVenue(a.DB, name, address, cityID)
	if err != nil {

		slog.Error("Failed to create venue.", "err", err)
		session.AddFlash(framework.Flash{
			framework.FlashFail,
			"Failed to create venue.",
		})

		session.Save(r, w)
		http.Redirect(w, r, "/events/"+e.IDString(), http.StatusFound)
		return
	}

	e.Venue = v
	e.VenueID = &v.ID
	e.Save()

	http.Redirect(w, r, "/events/"+e.IDString(), http.StatusFound)
	return
}

/*
 * Processing of the new venue event page, for www.
 */
func (a *app) venueWWWPost(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "login")

	// middlewareLIO ensures we have a user
	u := r.Context().Value("user").(*db.User)
	// middlewareEvent ensures we have an Event
	e := r.Context().Value("event").(*db.Event)

	// Only the owner of the group can add/edit an event
	if !e.TheGroup.HasCreate(u.ID) {

		session.AddFlash(framework.Flash{
			framework.FlashFail,
			"You don't have permission to add a venue to this event.",
		})

		session.Save(r, w)
		http.Redirect(w, r, "/events/"+e.IDString(), http.StatusFound)
		return
	}

	r.ParseForm()
	defer r.Body.Close()

	urlWWW, err := url.Parse(r.Form.Get("url"))
	if err != nil {

		slog.Error("URL is not valid.", "url", urlWWW)
		session.AddFlash(framework.Flash{
			framework.FlashFail,
			"URL was invalid.",
		})

		session.Save(r, w)
		http.Redirect(w, r, "/events/"+e.IDString()+"/new-venue#www", http.StatusFound)
		return
	}

	e.LocationURL = urlWWW.String()
	e.Save()

	http.Redirect(w, r, "/events/"+e.IDString(), http.StatusFound)
	return
}
