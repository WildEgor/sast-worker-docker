package vul_checker

type VulListItem struct {
	Code string
	Msg  string
}

type IVulChecker interface {
	GetCVEByImageNames(names []string) []VulListItem
}

type TrivyImageVulItems struct {
	VulnerabilityID string `json:"VulnerabilityID"`
	Title           string `json:"Title"`
	Desc            string `json:"Description"`
}

type TrivyImageCheckResultItem struct {
	Target          string               `json:"Target"`
	Vulnerabilities []TrivyImageVulItems `json:"Vulnerabilities"`
}

type TrivyImageCheckResult struct {
	Results []TrivyImageCheckResultItem `json:"Results"`
}
