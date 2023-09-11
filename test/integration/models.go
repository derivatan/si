//go:build integration

package integration

import (
	"github.com/derivatan/si"
	"github.com/gofrs/uuid"
	"time"
)

type contact struct {
	si.Model

	Email          string
	Phone          int
	RadioFrequency float64
	LastCall       time.Time
	OnSocialMedia  bool
	ArtistID       uuid.UUID

	artist si.RelationData[artist]
}

func (c contact) GetModel() si.Model {
	return c.Model
}

func (c contact) GetTable() string {
	return "contacts"
}

func (c contact) Artist() *si.Relation[contact, artist] {
	return si.BelongsTo[contact, artist](c, "artist", func(c *contact) *si.RelationData[artist] {
		return &c.artist
	})
}

type artist struct {
	si.Model

	Name string

	contact si.RelationData[contact]
	albums  si.RelationData[album]
}

func (a artist) GetModel() si.Model {
	return a.Model
}

func (a artist) GetTable() string {
	return "artists"
}

func (a artist) Contact() *si.Relation[artist, contact] {
	return si.HasOne[artist, contact](a, "contact", func(a *artist) *si.RelationData[contact] {
		return &a.contact
	})
}

func (a artist) Albums() *si.Relation[artist, album] {
	return si.HasMany[artist, album](a, "albums", func(a *artist) *si.RelationData[album] {
		return &a.albums
	})
}

type album struct {
	si.Model

	Name     string
	Year     int
	ArtistID uuid.UUID

	artist si.RelationData[artist]
}

func (a album) GetModel() si.Model {
	return a.Model
}

func (a album) GetTable() string {
	return "albums"
}

func (a album) Artist() *si.Relation[album, artist] {
	return si.BelongsTo[album, artist](a, "artist", func(a *album) *si.RelationData[artist] {
		return &a.artist
	})
}
