package domain

import (
	"time"
)

func ValidarDataNoPassado(data time.Time) error {

	// Normaliza a data recebida para UTC (somente DATA)
	dataUTC := data.UTC()
	dataRecebida := time.Date(
		dataUTC.Year(),
		dataUTC.Month(),
		dataUTC.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	// Hoje em UTC (somente DATA)
	agoraUTC := time.Now().UTC()
	hoje := time.Date(
		agoraUTC.Year(),
		agoraUTC.Month(),
		agoraUTC.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	if dataRecebida.Before(hoje) {
		return ErrDataEstaNoPassado
	}

	return nil
}

