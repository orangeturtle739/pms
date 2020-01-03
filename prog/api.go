package prog

import (
	"fmt"
	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/db"
	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/list"
	"github.com/ambientsound/pms/log"
	pms_mpd "github.com/ambientsound/pms/mpd"
	"github.com/ambientsound/pms/multibar"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/spotify/tracklist"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/topbar"
	"github.com/spf13/viper"
	"github.com/zmb3/spotify"
	"time"
)

func (v *Visp) Authenticate() error {
	err := v.setupAuthenticator()
	if err != nil {
		return fmt.Errorf("cannot authenticate with Spotify: %s", err)
	}
	url := v.Auth.AuthURL()
	log.Infof("Please authenticate with Spotify at: %s", url)

	return nil
}

func (v *Visp) Clipboard() songlist.Songlist {
	return v.clipboard
}

func (v *Visp) Db() *db.Instance {
	return nil // FIXME
}

func (v *Visp) Exec(command string) error {
	log.Debugf("Run command: %s", command)
	return v.interpreter.Exec(command)
}

func (v *Visp) Library() *songlist.Library {
	return nil // FIXME
}

func (v *Visp) List() list.List {
	return v.Termui.TableWidget().List()
}

func (v *Visp) ListChanged() {
	// FIXME
}

func (v *Visp) Message(fmt string, a ...interface{}) {
	log.Infof(fmt, a...)
	log.Debugf("Using obsolete Message() for previous message")
}

func (v *Visp) MpdClient() *mpd.Client {
	log.Debugf("nil mpd client; might break")
	return nil // FIXME
}

func (v *Visp) OptionChanged(key string) {
	switch key {
	case options.LogFile:
		logFile := v.Options().GetString(options.LogFile)
		overwrite := v.Options().GetBool(options.LogOverwrite)
		if len(logFile) == 0 {
			break
		}
		err := log.Configure(logFile, overwrite)
		if err != nil {
			log.Errorf("log configuration: %s", err)
			break
		}
		log.Infof("Note: log file will be backfilled with existing log")
		log.Infof("Writing debug log to %s", logFile)

	case options.PollInterval:
		interval := v.Options().GetInt(options.PollInterval)
		v.ticker = time.NewTicker(time.Duration(interval) * time.Second)

	case options.Topbar:
		config := v.Options().GetString(options.Topbar)
		matrix, err := topbar.Parse(v, config)
		if err == nil {
			_ = matrix
			// FIXME
			// v.Termui.Widgets.Topbar.SetMatrix(matrix)
		} else {
			log.Errorf("topbar configuration: %s", err)
		}
	}
}

func (v *Visp) Options() api.Options {
	return viper.GetViper()
}

func (v *Visp) PlayerStatus() (p pms_mpd.PlayerStatus) {
	return // FIXME
}

func (v *Visp) Queue() *songlist.Queue {
	log.Debugf("nil queue; might break")
	return nil // FIXME
}

func (v *Visp) Quit() {
	v.quit <- new(interface{})
}

func (v *Visp) Sequencer() *keys.Sequencer {
	return v.sequencer
}

func (v *Visp) Multibar() *multibar.Multibar {
	return v.multibar
}

func (v *Visp) Spotify() (*spotify.Client, error) {
	if v.client == nil {
		return nil, fmt.Errorf("please run `auth` to authenticate with Spotify")
	}
	err := v.setupAuthenticator()
	if err != nil {
		return nil, fmt.Errorf("unable to obtain Spotify client: %s", err.Error())
	}
	token, err := v.client.Token()
	if err != nil {
		return nil, fmt.Errorf("unable to refresh Spotify token: %s", err)
	}
	_ = v.Tokencache.Write(*token)
	return v.client, nil
}

func (v *Visp) Song() *song.Song {
	log.Debugf("nil song; might break")
	return nil
}

func (v *Visp) Songlist() songlist.Songlist {
	log.Debugf("nil songlist; might break")
	return nil
}

func (v *Visp) Songlists() []songlist.Songlist {
	log.Debugf("nil songlists; might break")
	return nil // FIXME
}

func (v *Visp) Styles() style.Stylesheet {
	return v.stylesheet
}

func (v *Visp) Tracklist() *spotify_tracklist.List {
	switch v := v.UI().TableWidget().List().(type) {
	case *spotify_tracklist.List:
		return v
	default:
		return nil
	}
}

func (v *Visp) UI() api.UI {
	return v.Termui
}
