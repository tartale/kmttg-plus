package model

func (t *Tivo) ServerName() string {
	return t.Tsn[0:3] + "-" + t.Tsn[3:7] + "-" + t.Tsn[7:11] + "-" + t.Tsn[11:15]
}
