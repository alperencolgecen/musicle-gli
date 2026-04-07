package main

import (
	"os"
	"path/filepath"

	"musicle-cli/bridge"
	"musicle-cli/state"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	// Resolve config directory
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		cfgDir = os.TempDir()
	}
	state.Current.ConfigDir = filepath.Join(cfgDir, "musicle")

	// Determine project dir for Python engine
	exe, err := os.Executable()
	if err != nil {
		exe = "."
	}
	projectDir := filepath.Dir(exe)
	// In dev mode the binary is in the project root already
	bridge.Init(projectDir)

	// Attempt to load existing config
	if err := state.Current.LoadConfig(); err != nil {
		// No config yet: first launch
		state.Current.IsFirstLaunch = true
		state.Current.Language = state.LangEnglish
	} else {
		// Config loaded: scan profiles
		if scanErr := state.Current.ScanProfiles(); scanErr != nil || len(state.Current.Profiles) == 0 {
			state.Current.IsFirstLaunch = true
		} else {
			state.Current.CurrentProfile = &state.Current.Profiles[0]
			if len(state.Current.CurrentProfile.Playlists) > 0 {
				state.Current.CurrentPlaylist = &state.Current.CurrentProfile.Playlists[0]
			}
		}
	}

	app := tview.NewApplication()
	app.EnableMouse(true)

	pages := tview.NewPages()

	// Check if first launch
	if state.Current.IsFirstLaunch {
		// Show setup wizard
		setupPage := NewSetupWizard(app, pages)
		pages.AddPage("setup", setupPage.Root(), true, true)
		pages.AddPage("home", NewHomePage(app, pages).Root(), true, false)
		pages.AddPage("settings", NewSettingsPage(app, pages).Root(), true, false)

		// Block Ctrl+C globally during setup too
		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			// F1: Home page (if available)
			if event.Key() == tcell.KeyF1 {
				pages.SwitchToPage("home")
				return nil
			}
			// F2: Settings page (if available)
			if event.Key() == tcell.KeyF2 {
				pages.SwitchToPage("settings")
				return nil
			}
			if event.Key() == tcell.KeyCtrlC {
				return nil // Prevent app termination
			}
			return event
		})

		app.SetRoot(pages, true).Run()
	} else {
		// Show home page directly
		homePage := NewHomePage(app, pages)
		settingsPage := NewSettingsPage(app, pages)
		pages.AddPage("setup", NewSetupWizard(app, pages).Root(), true, false)
		pages.AddPage("home", homePage.Root(), true, true)
		pages.AddPage("settings", settingsPage.Root(), true, false)

		// Global navigation keys and Esc key
		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			// F1: Home page
			if event.Key() == tcell.KeyF1 {
				pages.SwitchToPage("home")
				return nil
			}
			// F2: Settings page
			if event.Key() == tcell.KeyF2 {
				pages.SwitchToPage("settings")
				return nil
			}
			if event.Key() == tcell.KeyEsc {
				name, _ := pages.GetFrontPage()
				if name == "home" || name == "settings" {
					showExitDialog(app, pages)
					return nil
				}
			}
			// Block Ctrl+C globally
			if event.Key() == tcell.KeyCtrlC {
				return nil // Prevent app termination
			}
			return event
		})

		if err := app.SetRoot(pages, true).Run(); err != nil {
			panic(err)
		}
	}

	bridge.StopAll()
}

// showExitDialog shows a confirmation modal to quit the app
func showExitDialog(app *tview.Application, pages *tview.Pages) {
	const dialogName = "exitDialog"
	if pages.HasPage(dialogName) {
		return
	}

	modal := tview.NewModal().
		SetText("Exit MusicLe?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(_ int, label string) {
			pages.RemovePage(dialogName)
			if label == "Quit" {
				app.Stop()
			}
		})

	modal.SetBackgroundColor(tcell.NewRGBColor(30, 30, 30))
	modal.SetTextColor(tcell.ColorWhite)
	modal.SetBorderColor(tcell.NewRGBColor(29, 185, 84))

	pages.AddPage(dialogName, modal, false, true)
	app.SetFocus(modal)
}
