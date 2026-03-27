package main

import (
	"os"
	"path/filepath"
	"strings"

	"musicle-cli/state"
	"musicle-cli/ui"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SettingsPage manages the Profile and Playlist settings tabs
type SettingsPage struct {
	app   *tview.Application
	pages *tview.Pages
	root  *tview.Flex

	// Profile tab fields
	profileDrop   *tview.DropDown
	avatarInput   *tview.InputField
	nameInput     *tview.InputField
	bioInput      *tview.InputField
	langDrop      *tview.DropDown
	profileStatus *tview.TextView

	// Playlist tab fields
	playlistDrop   *tview.DropDown
	artInput       *tview.InputField
	plNameInput    *tview.InputField
	plBioInput     *tview.InputField
	playlistStatus *tview.TextView

	// Active tab: "profile" or "playlist"
	activeTab string
	tabBar    *tview.TextView
	tabPages  *tview.Pages
}

// NewSettingsPage constructs the settings page
func NewSettingsPage(app *tview.Application, pages *tview.Pages) *SettingsPage {
	s := &SettingsPage{
		app:       app,
		pages:     pages,
		activeTab: "profile",
	}
	s.build()
	return s
}

// Root returns the embeddable primitive
func (s *SettingsPage) Root() tview.Primitive { return s.root }

// ── Build ─────────────────────────────────────────────────────────────────────

func (s *SettingsPage) build() {
	// Header
	header := tview.NewTextView()
	header.SetDynamicColors(true)
	header.SetBackgroundColor(ui.ColorBackground)
	header.SetText("[white::b]Music[#1DB954::b]Le[-::-]    [#1DB954::r] Home [-::-]  [#1DB954::r] Settings [-::-]   [#B3B3B3][Esc] Back  [Tab] Switch Tab[-]")

	// Back hint
	backHint := tview.NewTextView()
	backHint.SetDynamicColors(true)
	backHint.SetBackgroundColor(ui.ColorBackground)
	backHint.SetText("  [#B3B3B3][Esc] Back to Home  [Ctrl+I] Save[-]")

	// Tab bar
	s.tabBar = tview.NewTextView()
	s.tabBar.SetDynamicColors(true)
	s.tabBar.SetBackgroundColor(ui.ColorBackground)
	s.refreshTabBar()

	// Tab content pages
	s.tabPages = tview.NewPages()
	s.tabPages.SetBackgroundColor(ui.ColorBackground)

	profileTab := s.buildProfileTab()
	playlistTab := s.buildPlaylistTab()

	s.tabPages.AddPage("profile", profileTab, true, true)
	s.tabPages.AddPage("playlist", playlistTab, true, false)

	// Root
	s.root = tview.NewFlex().SetDirection(tview.FlexRow)
	s.root.SetBackgroundColor(ui.ColorBackground)
	s.root.AddItem(header, 1, 0, false)
	s.root.AddItem(s.tabBar, 1, 0, false)
	s.root.AddItem(s.tabPages, 0, 1, true)
	s.root.AddItem(backHint, 1, 0, false)

	s.root.SetInputCapture(s.handleKeys)
}

func (s *SettingsPage) refreshTabBar() {
	profileStyle := "[#1DB954::r] Profile [-::-]"
	playlistStyle := "[#B3B3B3]  Playlist [-]"
	if s.activeTab == "playlist" {
		profileStyle = "[#B3B3B3]  Profile [-]"
		playlistStyle = "[#1DB954::r] Playlist [-::-]"
	}
	s.tabBar.SetText("  " + profileStyle + "  " + playlistStyle)
}

// ── Profile Tab ───────────────────────────────────────────────────────────────

func (s *SettingsPage) buildProfileTab() *tview.Flex {
	langT := func(en, tr string) string { return state.T(state.Current.Language, en, tr) }

	s.profileDrop = tview.NewDropDown()
	s.profileDrop.SetLabel("  " + langT("Profile", "Profil") + ":  ")
	s.profileDrop.SetLabelColor(ui.ColorAccent)
	s.profileDrop.SetFieldBackgroundColor(tcell.NewRGBColor(40, 40, 40))
	s.profileDrop.SetFieldTextColor(ui.ColorPrimary)
	s.profileDrop.SetPrefixTextColor(ui.ColorAccent)
	s.profileDrop.SetBackgroundColor(ui.ColorBackground)
	s.refreshProfileDrop()

	s.avatarInput = makeInput("  "+langT("Avatar Path", "Avatar Yolu")+":  ", langT("optional", "isteğe bağlı"))
	s.nameInput = makeInput("  "+langT("Display Name", "Görünen Ad")+":  ", "")
	s.bioInput = makeInput("  "+langT("Bio", "Biyografi")+":           ", "")

	s.langDrop = tview.NewDropDown()
	s.langDrop.SetLabel("  " + langT("Language", "Dil") + ":           ")
	s.langDrop.SetOptions([]string{"English", "Türkçe"}, nil)
	s.langDrop.SetLabelColor(ui.ColorAccent)
	s.langDrop.SetFieldBackgroundColor(tcell.NewRGBColor(40, 40, 40))
	s.langDrop.SetFieldTextColor(ui.ColorPrimary)
	s.langDrop.SetPrefixTextColor(ui.ColorAccent)
	s.langDrop.SetBackgroundColor(ui.ColorBackground)

	// Pre-fill from current profile
	s.fillProfileFields()

	s.profileStatus = tview.NewTextView()
	s.profileStatus.SetDynamicColors(true)
	s.profileStatus.SetBackgroundColor(ui.ColorBackground)

	saveBtn := tview.NewButton(langT("  Save Profile  ", "  Profili Kaydet  "))
	saveBtn.SetBackgroundColor(ui.ColorAccent)
	saveBtn.SetLabelColor(tcell.ColorBlack)
	saveBtn.SetActivatedStyle(tcell.StyleDefault.Background(ui.ColorAccent).Foreground(tcell.ColorBlack).Bold(true))
	saveBtn.SetSelectedFunc(func() { s.saveProfile() })

	tab := tview.NewFlex().SetDirection(tview.FlexRow)
	tab.SetBackgroundColor(ui.ColorBackground)
	tab.SetBorder(true)
	tab.SetBorderColor(ui.ColorBorder)
	tab.SetTitle(" " + langT("Profile Settings", "Profil Ayarları") + " ")
	tab.SetTitleColor(ui.ColorPrimary)

	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	tab.AddItem(s.profileDrop, 1, 0, true)
	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	tab.AddItem(s.avatarInput, 1, 0, false)
	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	tab.AddItem(s.nameInput, 1, 0, false)
	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	tab.AddItem(s.bioInput, 1, 0, false)
	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	tab.AddItem(s.langDrop, 1, 0, false)
	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	tab.AddItem(s.profileStatus, 1, 0, false)
	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 0, 1, false)
	tab.AddItem(saveBtn, 1, 0, false)
	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)

	return tab
}

func (s *SettingsPage) refreshProfileDrop() {
	var opts []string
	for _, p := range state.Current.Profiles {
		opts = append(opts, p.DisplayName+" ("+p.FolderName+")")
	}
	if len(opts) == 0 {
		opts = []string{"(no profiles)"}
	}
	s.profileDrop.SetOptions(opts, func(_ string, idx int) {
		if idx < len(state.Current.Profiles) {
			state.Current.CurrentProfile = &state.Current.Profiles[idx]
			s.fillProfileFields()
		}
	})
	s.profileDrop.SetCurrentOption(0)
}

func (s *SettingsPage) fillProfileFields() {
	if state.Current.CurrentProfile == nil || s.avatarInput == nil || s.nameInput == nil || s.bioInput == nil {
		return
	}
	p := state.Current.CurrentProfile
	s.avatarInput.SetText(p.AvatarPath)
	s.nameInput.SetText(p.DisplayName)
	s.bioInput.SetText(p.Bio)
	if p.Language == state.LangTurkish {
		s.langDrop.SetCurrentOption(1)
	} else {
		s.langDrop.SetCurrentOption(0)
	}
}

func (s *SettingsPage) saveProfile() {
	langT := func(en, tr string) string { return state.T(state.Current.Language, en, tr) }
	if state.Current.CurrentProfile == nil {
		s.profileStatus.SetText("[#FF4444]  No profile selected[-]")
		return
	}
	name := strings.TrimSpace(s.nameInput.GetText())
	bio := strings.TrimSpace(s.bioInput.GetText())
	avatar := strings.TrimSpace(s.avatarInput.GetText())

	langIdx, _ := s.langDrop.GetCurrentOption()
	newLang := state.LangEnglish
	if langIdx == 1 {
		newLang = state.LangTurkish
	}

	state.Current.Language = newLang // Update global language too
	state.Current.CurrentProfile.Language = newLang

	if err := state.Current.SaveProfileMeta(state.Current.CurrentProfile.FolderName, name, bio); err != nil {
		s.profileStatus.SetText("[#FF4444]  ✗ " + err.Error() + "[-]")
		return
	}
	// Also save lang.txt
	_ = os.WriteFile(filepath.Join(state.Current.ProfilesDir(), state.Current.CurrentProfile.FolderName, "lang.txt"), []byte(string(newLang)), 0644)

	if avatar != "" && avatar != state.Current.CurrentProfile.AvatarPath {
		_ = state.CopyFile(avatar, state.Current.PlaylistDir(state.Current.CurrentProfile.FolderName, "avatar"))
	}
	state.Current.CurrentProfile.DisplayName = name
	state.Current.CurrentProfile.Bio = bio
	s.profileStatus.SetText("[#1DB954]  ✓ " + langT("Saved!", "Kaydedildi!") + "[-]")
	_ = state.Current.SaveConfig()
}

// ── Playlist Tab ──────────────────────────────────────────────────────────────

func (s *SettingsPage) buildPlaylistTab() *tview.Flex {
	langT := func(en, tr string) string { return state.T(state.Current.Language, en, tr) }

	s.playlistDrop = tview.NewDropDown()
	s.playlistDrop.SetLabel("  " + langT("Playlist", "Playlist") + ":  ")
	s.playlistDrop.SetLabelColor(ui.ColorAccent)
	s.playlistDrop.SetFieldBackgroundColor(tcell.NewRGBColor(40, 40, 40))
	s.playlistDrop.SetFieldTextColor(ui.ColorPrimary)
	s.playlistDrop.SetPrefixTextColor(ui.ColorAccent)
	s.playlistDrop.SetBackgroundColor(ui.ColorBackground)
	s.refreshPlaylistDrop()

	s.artInput = makeInput("  "+langT("Art Path", "Görsel Yolu")+":        ", langT("optional", "isteğe bağlı"))
	s.plNameInput = makeInput("  "+langT("Playlist Name", "Playlist Adı")+":   ", "")
	s.plBioInput = makeInput("  "+langT("Description", "Açıklama")+":       ", "")
	s.fillPlaylistFields()

	s.playlistStatus = tview.NewTextView()
	s.playlistStatus.SetDynamicColors(true)
	s.playlistStatus.SetBackgroundColor(ui.ColorBackground)

	saveBtn := tview.NewButton(langT("  Save  ", "  Kaydet  "))
	saveBtn.SetBackgroundColor(ui.ColorAccent)
	saveBtn.SetLabelColor(tcell.ColorBlack)
	saveBtn.SetActivatedStyle(tcell.StyleDefault.Background(ui.ColorAccent).Foreground(tcell.ColorBlack).Bold(true))
	saveBtn.SetSelectedFunc(func() { s.savePlaylist() })

	deleteBtn := tview.NewButton(langT("  Delete  ", "  Sil  "))
	deleteBtn.SetBackgroundColor(ui.ColorError)
	deleteBtn.SetLabelColor(tcell.ColorWhite)
	deleteBtn.SetActivatedStyle(tcell.StyleDefault.Background(ui.ColorError).Foreground(tcell.ColorWhite).Bold(true))
	deleteBtn.SetSelectedFunc(func() { s.confirmDeletePlaylist() })

	btnRow := tview.NewFlex()
	btnRow.SetBackgroundColor(ui.ColorBackground)
	btnRow.AddItem(saveBtn, 12, 0, false)
	btnRow.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 2, 0, false)
	btnRow.AddItem(deleteBtn, 12, 0, false)
	btnRow.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 0, 1, false)

	tab := tview.NewFlex().SetDirection(tview.FlexRow)
	tab.SetBackgroundColor(ui.ColorBackground)
	tab.SetBorder(true)
	tab.SetBorderColor(ui.ColorBorder)
	tab.SetTitle(" " + langT("Playlist Settings", "Playlist Ayarları") + " ")
	tab.SetTitleColor(ui.ColorPrimary)

	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	tab.AddItem(s.playlistDrop, 1, 0, true)
	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	tab.AddItem(s.artInput, 1, 0, false)
	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	tab.AddItem(s.plNameInput, 1, 0, false)
	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	tab.AddItem(s.plBioInput, 1, 0, false)
	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	tab.AddItem(s.playlistStatus, 1, 0, false)
	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 0, 1, false)
	tab.AddItem(btnRow, 1, 0, false)
	tab.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)

	return tab
}

func (s *SettingsPage) refreshPlaylistDrop() {
	var opts []string
	if state.Current.CurrentProfile != nil {
		for _, pl := range state.Current.CurrentProfile.Playlists {
			opts = append(opts, pl.Name+" ("+pl.FolderName+")")
		}
	}
	if len(opts) == 0 {
		opts = []string{"(no playlists)"}
	}
	s.playlistDrop.SetOptions(opts, func(_ string, idx int) {
		if state.Current.CurrentProfile != nil && idx < len(state.Current.CurrentProfile.Playlists) {
			state.Current.CurrentPlaylist = &state.Current.CurrentProfile.Playlists[idx]
			s.fillPlaylistFields()
		}
	})
	s.playlistDrop.SetCurrentOption(0)
}

func (s *SettingsPage) fillPlaylistFields() {
	if state.Current.CurrentPlaylist == nil || s.artInput == nil || s.plNameInput == nil || s.plBioInput == nil {
		return
	}
	pl := state.Current.CurrentPlaylist
	s.artInput.SetText(pl.ArtPath)
	s.plNameInput.SetText(pl.Name)
	s.plBioInput.SetText(pl.Bio)
}

func (s *SettingsPage) savePlaylist() {
	langT := func(en, tr string) string { return state.T(state.Current.Language, en, tr) }
	if state.Current.CurrentProfile == nil || state.Current.CurrentPlaylist == nil {
		s.playlistStatus.SetText("[#FF4444]  No playlist selected[-]")
		return
	}
	name := strings.TrimSpace(s.plNameInput.GetText())
	bio := strings.TrimSpace(s.plBioInput.GetText())
	if err := state.Current.SavePlaylistMeta(
		state.Current.CurrentProfile.FolderName,
		state.Current.CurrentPlaylist.FolderName,
		name, bio,
	); err != nil {
		s.playlistStatus.SetText("[#FF4444]  ✗ " + err.Error() + "[-]")
		return
	}
	state.Current.CurrentPlaylist.Name = name
	state.Current.CurrentPlaylist.Bio = bio
	s.playlistStatus.SetText("[#1DB954]  ✓ " + langT("Saved!", "Kaydedildi!") + "[-]")
}

func (s *SettingsPage) confirmDeletePlaylist() {
	langT := func(en, tr string) string { return state.T(state.Current.Language, en, tr) }
	const dialogName = "deleteConfirm"
	if state.Current.CurrentPlaylist == nil {
		return
	}
	plName := state.Current.CurrentPlaylist.Name
	modal := tview.NewModal().
		SetText(langT("Delete playlist '"+plName+"'? This cannot be undone.",
			"'"+plName+"' silinsin mi? Bu geri alınamaz.")).
		AddButtons([]string{langT("Delete", "Sil"), langT("Cancel", "İptal")}).
		SetDoneFunc(func(_ int, label string) {
			s.pages.RemovePage(dialogName)
			if label == langT("Delete", "Sil") {
				s.deletePlaylist()
			}
		})
	modal.SetBackgroundColor(tcell.NewRGBColor(30, 30, 30))
	modal.SetBorderColor(ui.ColorError)
	s.pages.AddPage(dialogName, modal, false, true)
	s.app.SetFocus(modal)
}

func (s *SettingsPage) deletePlaylist() {
	langT := func(en, tr string) string { return state.T(state.Current.Language, en, tr) }
	if state.Current.CurrentProfile == nil || state.Current.CurrentPlaylist == nil {
		return
	}
	if err := state.Current.DeletePlaylist(
		state.Current.CurrentProfile.FolderName,
		state.Current.CurrentPlaylist.FolderName,
	); err != nil {
		s.playlistStatus.SetText("[#FF4444]  ✗ " + err.Error() + "[-]")
		return
	}
	state.Current.CurrentPlaylist = nil
	_ = state.Current.ScanProfiles()
	if len(state.Current.Profiles) > 0 {
		state.Current.CurrentProfile = &state.Current.Profiles[0]
		if len(state.Current.CurrentProfile.Playlists) > 0 {
			state.Current.CurrentPlaylist = &state.Current.CurrentProfile.Playlists[0]
		}
	}
	s.refreshPlaylistDrop()
	s.fillPlaylistFields()
	s.playlistStatus.SetText("[#1DB954]  ✓ " + langT("Deleted.", "Silindi.") + "[-]")
}

// ── Key handler ───────────────────────────────────────────────────────────────

func (s *SettingsPage) handleKeys(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEsc:
		s.pages.SwitchToPage("home")
		return nil
	case tcell.KeyTab:
		if event.Modifiers()&tcell.ModCtrl == 0 { // Real Tab, not Ctrl+I
			if s.activeTab == "profile" {
				s.activeTab = "playlist"
				s.tabPages.SwitchToPage("playlist")
			} else {
				s.activeTab = "profile"
				s.tabPages.SwitchToPage("profile")
			}
			s.refreshTabBar()
			return nil
		}
	case tcell.KeyEnter:
		s.saveProfile()
		s.savePlaylist()
		return nil
	}
	return event
}
