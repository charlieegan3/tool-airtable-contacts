package tool

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/charlieegan3/tool-airtable-contacts/pkg/tool/jobs"

	"github.com/Jeffail/gabs/v2"
	"github.com/charlieegan3/toolbelt/pkg/apis"
	"github.com/gorilla/mux"
)

// AirtableContacts is a tool that wraps a number of jobs to manage an airtable contacts database
type AirtableContacts struct {
	config *gabs.Container
	db     *sql.DB
}

func (a *AirtableContacts) Name() string {
	return "airtable-contacts"
}

func (a *AirtableContacts) FeatureSet() apis.FeatureSet {
	return apis.FeatureSet{
		Config: true,
		Jobs:   true,
	}
}

func (a *AirtableContacts) SetConfig(config map[string]any) error {
	a.config = gabs.Wrap(config)

	return nil
}
func (a *AirtableContacts) Jobs() ([]apis.Job, error) {
	var j []apis.Job
	var path string
	var ok bool

	path = "endpoint"
	endpoint, ok := a.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	// load airtable config
	path = "airtable.key"
	airtableKey, ok := a.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "airtable.base"
	airtableBase, ok := a.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "airtable.table"
	airtableTable, ok := a.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "airtable.view"
	airtableView, ok := a.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	// load CardDAV config
	path = "carddav.server_url"
	cardDAVServer, ok := a.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "carddav.user"
	cardDAVUser, ok := a.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "carddav.password"
	cardDAVPassword, ok := a.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	// load VCard Config
	path = "vcard.use_v3"
	vCardV3, ok := a.config.Path(path).Data().(bool)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "vcard.photo.size"
	vCardPhotoSize, ok := a.config.Path(path).Data().(int)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	// load schedules
	path = "jobs.day.schedule"
	daySchedule, ok := a.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "jobs.week.schedule"
	weekSchedule, ok := a.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "jobs.sync.schedule"
	syncSchedule, ok := a.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	return []apis.Job{
		&jobs.Day{
			ScheduleOverride: daySchedule,
			Endpoint:         endpoint,
			AirtableKey:      airtableKey,
			AirtableBase:     airtableBase,
			AirtableTable:    airtableTable,
			AirtableView:     airtableView,
		},
		&jobs.Week{
			ScheduleOverride: weekSchedule,
			Endpoint:         endpoint,
			AirtableKey:      airtableKey,
			AirtableBase:     airtableBase,
			AirtableTable:    airtableTable,
			AirtableView:     airtableView,
		},
		&jobs.Sync{
			ScheduleOverride: syncSchedule,
			Endpoint:         endpoint,
			AirtableKey:      airtableKey,
			AirtableBase:     airtableBase,
			AirtableTable:    airtableTable,
			AirtableView:     airtableView,
			CardDAVServer:    cardDAVServer,
			CardDAVUser:      cardDAVUser,
			CardDAVPassword:  cardDAVPassword,
			VCardV3:          vCardV3,
			VCardPhotoSize:   vCardPhotoSize,
		},
	}, nil
}

func (a *AirtableContacts) ExternalJobsFuncSet(f func(job apis.ExternalJob) error) {}

func (a *AirtableContacts) DatabaseMigrations() (*embed.FS, string, error) {
	return &embed.FS{}, "migrations", nil
}
func (a *AirtableContacts) DatabaseSet(db *sql.DB)              {}
func (a *AirtableContacts) HTTPPath() string                    { return "" }
func (a *AirtableContacts) HTTPAttach(router *mux.Router) error { return nil }
func (a *AirtableContacts) HTTPHost() string                    { return "" }
