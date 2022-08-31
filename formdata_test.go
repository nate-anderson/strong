package strong

import (
	"net/url"
	"testing"
)

func TestUnmarshalForm(t *testing.T) {
	var (
		nameFirst = "Jane"
		nameLast  = "Doe"
		email     = "admin@example.com"
		birthYear = "1977"
	)
	fd := url.Values{}
	fd.Set("name_first", nameFirst)
	fd.Set("name_last", nameLast)
	fd.Set("email", email)
	fd.Set("birth_year", birthYear)

	t.Run("unmarshal to struct", func(t *testing.T) {
		t.Skip()
		type Form struct {
			NameFirst string `form:"name_first"`
			NameLast  string `form:"name_last"`
			Email     string `form:"email"`
			BirthYear string `form:"birth_year"`
		}
		f := Form{}
		err := unmarshalForm(fd, &f)
		if err != nil {
			t.Error(err)
		}

		if f.BirthYear != birthYear {
			t.Errorf("mismatched birth_year: expect %s, got %s", birthYear, f.BirthYear)
		}
	})

	t.Run("unmarshal to map", func(t *testing.T) {
		m := map[string]string{}
		err := unmarshalForm(fd, m)
		if err != nil {
			t.Error(err)
		}

		if m["birth_year"] != birthYear {
			t.Errorf("mismatched birth_year: expect %s, got %s", birthYear, m["birth_year"])
		}
	})

}
