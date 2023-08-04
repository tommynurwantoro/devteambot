package resty

import "github.com/go-resty/resty/v2"

type MyQuran struct {
	Client *resty.Client
}

func (m *MyQuran) Startup() error {
	m.Client = resty.New().EnableTrace()
	m.Client.SetRetryCount(3)
	m.Client.SetBaseURL("https://api.myquran.com/v1")

	return nil
}

func (m *MyQuran) Shutdown() error { return nil }

type GetSholatResponse struct {
	Data struct {
		Jadwal struct {
			Tanggal string `json:"tanggal"`
			Dzuhur  string `json:"dzuhur"`
			Ashar   string `json:"ashar"`
		} `json:"jadwal"`
	} `json:"data"`
}
