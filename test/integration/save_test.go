package integration

import (
	"github.com/derivatan/si"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	db := DB(t)
	ids := Seed(db, []artist{
		{Name: "Roger Waters"},
	})
	c := contact{
		Email:          "roger@waters.com",
		Phone:          1357924680,
		RadioFrequency: 21.789,
		LastCall:       time.Date(1943, time.September, 6, 5, 4, 3, 0, time.Local),
		OnSocialMedia:  true,
		ArtistID:       ids[0],
	}
	err := si.Save(db, &c)
	c2 := si.Query[contact]().MustFind(db)

	assert.NoError(t, err)
	assert.NotNil(t, c2)
	assert.Equal(t, c.Email, c2.Email)
	assert.Equal(t, c.Phone, c2.Phone)
	assert.Equal(t, c.RadioFrequency, c2.RadioFrequency)
	assert.Equal(t, c.LastCall, c2.LastCall.Local())
	assert.Equal(t, c.OnSocialMedia, c2.OnSocialMedia)
	assert.Equal(t, c.ArtistID, c2.ArtistID)
}

func TestUpdate(t *testing.T) {
	db := DB(t)
	ids := Seed(db, []artist{
		{Name: "Timbuktu"},
	})
	Seed(db, []contact{
		{

			// TODO: Nullable values

			Email:          "Ett brev",
			Phone:          192837465,
			RadioFrequency: 98.54,
			LastCall:       time.Now(),
			OnSocialMedia:  false,
			ArtistID:       ids[0],
		},
	})

	email := "Det Löser sig"
	phone := 7592836
	radio := 73.11
	lastCall := time.Date(1234, 5, 6, 7, 8, 9, 0, time.Local)
	onSM := true
	c := si.Query[contact]().MustFind(db)
	c.Email = email
	c.Phone = phone
	c.RadioFrequency = radio
	c.LastCall = lastCall
	c.OnSocialMedia = onSM
	err := si.Save(db, c)

	c2 := si.Query[contact]().MustFind(db)

	assert.NoError(t, err)
	assert.Equal(t, email, c2.Email)
	assert.Equal(t, phone, c2.Phone)
	assert.Equal(t, radio, c2.RadioFrequency)
	assert.Equal(t, lastCall, c2.LastCall.Local())
	assert.Equal(t, onSM, c2.OnSocialMedia)
}
