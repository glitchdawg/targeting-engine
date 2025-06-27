package storage

import (
	"database/sql"
	"strings"

	"github.com/glitchdawg/campaign-targeting-engine/internal/models"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresStore{DB: db}, nil
}



func (s *PostgresStore) GetCampaigns() (map[string]models.Campaign, error) {
	rows, err := s.DB.Query(`SELECT id, name, img, cta, status FROM campaigns`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	campaigns := make(map[string]models.Campaign)
	for rows.Next() {
		var c models.Campaign
		if err := rows.Scan(&c.ID, &c.Name, &c.Image, &c.CTA, &c.Status); err != nil {
			return nil, err
		}
		campaigns[c.ID] = c
	}
	return campaigns, nil
}

func (s *PostgresStore) GetCountryStates()(map[string][]string,error){
	
	row, err:=s.DB.Query(`SELECT country, states FROM country_states`)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	countryStates := make(map[string][]string)
	for row.Next() {
		var country, states string
		if err:=row.Scan(&country, &states); err != nil {
			return nil, err
		}
		stateList:=[]string{}
		if states!=""{
			for _, state:= range strings.Split(states,","){
				stateList=append(stateList, strings.TrimSpace(state))
			}
		}
		countryStates[country]=stateList
	}
	return countryStates, nil

}

func (s *PostgresStore) GetTargetingRules() (map[string]models.TargetingRule, error) {
	rows, err := s.DB.Query(`
        SELECT campaign_id, include_country, exclude_country, include_os, exclude_os, include_app, exclude_app, include_state, exclude_state,
        FROM targeting_rules`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rules := make(map[string]models.TargetingRule)
	for rows.Next() {
		var r models.TargetingRule
		var includeCountry, excludeCountry, includeOS, excludeOS, includeApp, excludeApp, includeState, excludeState sql.NullString
		if err := rows.Scan(
			&r.CampaignID, &includeCountry, &excludeCountry, &includeOS, &excludeOS, &includeApp, &excludeApp, &includeState, &excludeState,
		); err != nil {
			return nil, err
		}
		r.IncludeCountry = splitNullString(includeCountry)
		r.ExcludeCountry = splitNullString(excludeCountry)
		r.IncludeOS = splitNullString(includeOS)
		r.ExcludeOS = splitNullString(excludeOS)
		r.IncludeApp = splitNullString(includeApp)
		r.ExcludeApp = splitNullString(excludeApp)
		r.IncludeState = splitNullString(includeState)
		r.ExcludeState = splitNullString(excludeState)
		rules[r.CampaignID] = r
	}
	return rules, nil
}

func splitNullString(ns sql.NullString) []string {
	if !ns.Valid || ns.String == "" {
		return nil
	}
	parts := strings.Split(ns.String, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
