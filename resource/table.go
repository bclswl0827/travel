package resource

import (
	"embed"
	"path/filepath"
)

type ak135TableId string

const (
	AK135_P_SHALLOW ak135TableId = "ak135_P_shallow.txt"
	AK135_P_DEEP    ak135TableId = "ak135_P_deep.txt"
	AK135_S_SHALLOW ak135TableId = "ak135_S_shallow.txt"
	AK135_S_DEEP    ak135TableId = "ak135_S_deep.txt"
)

const (
	AK135_PcP ak135TableId = "ak135_PcP.txt"
	AK135_ScS ak135TableId = "ak135_ScS.txt"
	AK135_ScP ak135TableId = "ak135_ScP.txt"
)

const (
	AK135_PKPab ak135TableId = "ak135_PKPab.txt"
	AK135_PKPbc ak135TableId = "ak135_PKPbc.txt"
	AK135_PKPdf ak135TableId = "ak135_PKPdf.txt"
)

const (
	AK135_SKSac ak135TableId = "ak135_SKSac.txt"
	AK135_SKSdf ak135TableId = "ak135_SKSdf.txt"
	AK135_SKP   ak135TableId = "ak135_SKP.txt"
)

//go:embed data
var ak135Dir embed.FS

func GetAK135Table(tableId ak135TableId) string {
	dataBytes, err := ak135Dir.ReadFile(filepath.Join("data", string(tableId)))
	if err != nil {
		return ""
	}
	return string(dataBytes)
}
