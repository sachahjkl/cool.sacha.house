package common

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

type Me struct {
	Nom             string    `json:"nom"`
	Prenom          string    `json:"prenom"`
	Username        string    `json:"username"`
	Gender          string    `json:"gender"`
	DateNaissance   time.Time `json:"dateNaissance"`
	PlaceOfLiving   string    `json:"placeOfLiving"`
	Github          string    `json:"github"`
	Linkedin        string    `json:"linkedin"`
	Mail            string    `json:"mail"`
	CurriculumVitae string    `json:"curriculumVitae"`
	EthAddress      string    `json:"ethAddress"`
	MoneroAdress    string    `json:"moneroAdress"`
	Links           struct {
		Dotfiles string `json:"dotfiles"`
		Hayekfr  string `json:"hayekfr"`
	} `json:"links"`
}

var meCache *Me

// TODO(sacha): shortcuts to access preformatted variables
// export const PRETTY_PRENOM = capitalize(MOI.prenom.toLowerCase());
// export const PRETTY_NOM = MOI.nom.toUpperCase();
// export const SITE_TITLE = `${PRETTY_PRENOM} ${PRETTY_NOM}`;

func DataAboutMe() (*Me, error) {

	// use cache if present
	if meCache != nil {
		return meCache, nil
	}

	jsonFile, err := os.Open("./assets/me.json")
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	meData := Me{}

	json.Unmarshal(bytes, &meData)

	meCache = &meData

	return meCache, nil
}
