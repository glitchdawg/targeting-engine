package service

import (
	"testing"

	"github.com/glitchdawg/campaign-targeting-engine/internal/models"
)

func TestGetCampaigns(t *testing.T) {
	campaigns := map[string]models.Campaign{
		"spotify":      {ID: "spotify", Name: "Spotify", Image: "https://somelink", CTA: "Download", Status: "ACTIVE"},
		"duolingo":     {ID: "duolingo", Name: "Duolingo", Image: "https://somelink2", CTA: "Install", Status: "ACTIVE"},
		"subwaysurfer": {ID: "subwaysurfer", Name: "Subway Surfer", Image: "https://somelink3", CTA: "Play", Status: "ACTIVE"},
		"inactive":     {ID: "inactive", Name: "Inactive", Image: "https://inactive", CTA: "None", Status: "INACTIVE"},
	}
	rules := map[string]models.TargetingRule{
		"spotify": {
			CampaignID:     "spotify",
			IncludeCountry: []string{"US", "Canada"},
			IncludeOS:      []string{"Android", "iOS"}, // Only match Android/iOS
		},
		"duolingo": {
			CampaignID:     "duolingo",
			IncludeOS:      []string{"Android", "iOS"},
			ExcludeCountry: []string{"US"},
		},
		"subwaysurfer": {
			CampaignID: "subwaysurfer",
			IncludeOS:  []string{"Android"},
			IncludeApp: []string{"com.gametion.ludokinggame"},
		},
		"inactive": {
			CampaignID:     "inactive",
			IncludeCountry: []string{"US"},
		},
	}

	svc := NewDeliveryServiceImpl(campaigns, rules)

	tests := []struct {
		name    string
		app     string
		country string
		os      string
		wantIDs []string
		wantErr string
	}{
		{
			name:    "missing app",
			app:     "",
			country: "germany",
			os:      "android",
			wantErr: "missing app param",
		},
		{
			name:    "missing country",
			app:     "com.abc.xyz",
			country: "",
			os:      "android",
			wantErr: "missing country param",
		},
		{
			name:    "missing os",
			app:     "com.abc.xyz",
			country: "germany",
			os:      "",
			wantErr: "missing os param",
		},
		{
			name:    "no matching campaigns",
			app:     "com.abc.xyz",
			country: "us",
			os:      "web",
			wantIDs: []string{},
		},
		{
			name:    "single match duolingo",
			app:     "com.abc.xyz",
			country: "germany",
			os:      "android",
			wantIDs: []string{"duolingo"},
		},
		{
			name:    "multiple matches",
			app:     "com.gametion.ludokinggame",
			country: "us",
			os:      "android",
			wantIDs: []string{"spotify", "subwaysurfer"},
		},
		{
			name:    "inactive campaign not returned",
			app:     "com.abc.xyz",
			country: "us",
			os:      "android",
			wantIDs: []string{"spotify"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.GetCampaigns(tt.app, tt.country, tt.os)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Errorf("expected error %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.wantIDs) {
				t.Errorf("expected %d campaigns, got %d", len(tt.wantIDs), len(got))
			}
			gotIDs := map[string]bool{}
			for _, c := range got {
				gotIDs[c.ID] = true
			}
			for _, id := range tt.wantIDs {
				if !gotIDs[id] {
					t.Errorf("expected campaign %q in result", id)
				}
			}
		})
	}
}
