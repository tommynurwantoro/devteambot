package resty

import "github.com/go-resty/resty/v2"

type JadwalSholatOrg struct {
	Client *resty.Client
}

func (m *JadwalSholatOrg) Startup() error {
	m.Client = resty.New().EnableTrace()
	m.Client.SetRetryCount(3)
	m.Client.SetBaseURL("https://raw.githubusercontent.com/lakuapik/jadwalsholatorg/master/adzan/")

	return nil
}

func (m *JadwalSholatOrg) Shutdown() error { return nil }

type GetJadwalSholatResponse struct {
	Tanggal string `json:"tanggal"`
	Imsyak  string `json:"imsyak"`
	Shubuh  string `json:"shubuh"`
	Terbit  string `json:"terbit"`
	Dhuha   string `json:"dhuha"`
	Dzuhur  string `json:"dzuhur"`
	Ashar   string `json:"ashr"`
	Maghrib string `json:"magrib"`
	Isya    string `json:"isya"`
}
