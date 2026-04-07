package main

import (
	"strings"
	"time"

	"musicle-cli/state"
	"musicle-cli/ui"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sqweek/dialog"
)

// SetupWizard manages the 4-step first-launch setup flow
type SetupWizard struct {
	app      *tview.Application
	pages    *tview.Pages
	root     *tview.Pages // internal step pages
	step     int
	lang     state.Language
	rootDir  string
	profile  struct{ folder, avatar, name, bio string }
	playlist struct{ folder, art, name, bio string }
}

// NewSetupWizard creates and returns a new SetupWizard
func NewSetupWizard(app *tview.Application, mainPages *tview.Pages) *SetupWizard {
	sw := &SetupWizard{
		app:   app,
		pages: mainPages,
		root:  tview.NewPages(),
		lang:  state.LangEnglish,
	}
	sw.buildStep1()
	return sw
}

// Root returns the root primitive for embedding in the main Pages
func (sw *SetupWizard) Root() tview.Primitive { return sw.root }

// ── helpers ──────────────────────────────────────────────────────────────────

func (sw *SetupWizard) styledBox(title string) *tview.Flex {
	box := tview.NewFlex().SetDirection(tview.FlexRow)
	box.SetBackgroundColor(ui.ColorBackground)
	_ = title
	return box
}

func makeInput(label, placeholder string) *tview.InputField {
	f := tview.NewInputField()
	f.SetLabel(label)
	f.SetPlaceholder(placeholder)
	f.SetLabelColor(ui.ColorAccent)
	f.SetFieldBackgroundColor(tcell.NewRGBColor(40, 40, 40))
	f.SetFieldTextColor(ui.ColorPrimary)
	f.SetPlaceholderStyle(tcell.StyleDefault.Foreground(ui.ColorSecondary))
	f.SetBackgroundColor(ui.ColorBackground)

	// Re-enable Ctrl+C and Ctrl+V functions
	lastInputTime := time.Time{}
	f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Block F1/F2 in input fields to prevent accidental navigation
		if event.Key() == tcell.KeyF1 || event.Key() == tcell.KeyF2 {
			return nil // Block navigation keys while typing
		}

		// Allow copy/paste shortcuts
		if event.Modifiers()&tcell.ModCtrl != 0 {
			switch event.Key() {
			case tcell.KeyCtrlC, tcell.KeyCtrlV, tcell.KeyCtrlX:
				return event // Let terminal handle copy/paste
			}
		}

		// Time-based limiting for regular typing only
		now := time.Now()
		if now.Sub(lastInputTime) < 200*time.Millisecond && event.Modifiers() == 0 {
			return nil // Block rapid character repeats
		}
		lastInputTime = now
		return event
	})

	return f
}

func makeErrorView() *tview.TextView {
	tv := tview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetBackgroundColor(ui.ColorBackground)
	tv.SetTextColor(ui.ColorError)
	return tv
}

func hintBar(text string) *tview.TextView {
	tv := tview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetText("[#B3B3B3]" + text + "[-]")
	tv.SetBackgroundColor(tcell.NewRGBColor(18, 18, 18))
	return tv
}

func makeNextButton(onClick func()) *tview.Button {
	btn := tview.NewButton(" -> ")
	btn.SetBackgroundColor(ui.ColorAccent)
	btn.SetLabelColor(tcell.ColorBlack)
	btn.SetActivatedStyle(tcell.StyleDefault.Background(tcell.NewRGBColor(30, 215, 96)).Foreground(tcell.ColorBlack).Bold(true))
	btn.SetSelectedFunc(onClick)
	return btn
}

func logoText(lang state.Language) string {
	_ = lang
	return "[white::b]Music[#1DB954::b]Le[-::-]"
}

func centeredFlex(inner tview.Primitive, width, height int) *tview.Flex {
	// Horizontal centering
	hFlex := tview.NewFlex()
	hFlex.SetBackgroundColor(ui.ColorBackground)
	hFlex.AddItem(nil, 0, 1, false)      // Left spacer
	hFlex.AddItem(inner, width, 0, true) // Content with fixed width
	hFlex.AddItem(nil, 0, 1, false)      // Right spacer

	// Vertical centering
	vFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	vFlex.SetBackgroundColor(ui.ColorBackground)
	vFlex.AddItem(nil, 0, 1, false)       // Top spacer
	vFlex.AddItem(hFlex, height, 0, true) // Content with fixed height
	vFlex.AddItem(nil, 0, 1, false)       // Bottom spacer

	return vFlex
}

// ── Step 1 — Welcome & Language ───────────────────────────────────────────────

func (sw *SetupWizard) buildStep1() {
	logo := tview.NewTextView()
	logo.SetDynamicColors(true)
	logo.SetText("\n\n" + logoText(sw.lang) + "\n\n  [#B3B3B3]A Spotify-inspired CLI music player[-]")
	logo.SetBackgroundColor(ui.ColorBackground)
	logo.SetTextAlign(tview.AlignCenter)

	// Language dropdown
	dropdown := tview.NewDropDown()
	dropdown.SetLabel("  Language / Dil:  ")
	dropdown.SetOptions([]string{"English", "Türkçe"}, func(text string, _ int) {
		if text == "Türkçe" {
			sw.lang = state.LangTurkish
		} else {
			sw.lang = state.LangEnglish
		}
	})
	dropdown.SetCurrentOption(0)
	dropdown.SetLabelColor(ui.ColorAccent)
	dropdown.SetFieldBackgroundColor(tcell.NewRGBColor(40, 40, 40))
	dropdown.SetFieldTextColor(ui.ColorPrimary)
	dropdown.SetPrefixTextColor(ui.ColorAccent)
	dropdown.SetBackgroundColor(ui.ColorBackground)

	hint := hintBar("  [Enter] Next  |  [Esc] Quit")
	nextBtn := makeNextButton(func() {
		state.Current.Language = sw.lang
		sw.buildStep2()
	})

	controls := tview.NewFlex().SetDirection(tview.FlexColumn)
	controls.SetBackgroundColor(ui.ColorBackground)
	controls.AddItem(hint, 0, 1, false)
	controls.AddItem(nextBtn, 6, 0, true)

	inner := tview.NewFlex().SetDirection(tview.FlexRow)
	inner.SetBackgroundColor(ui.ColorBackground)
	inner.SetBorder(true)
	inner.SetBorderColor(ui.ColorAccent)
	inner.SetTitle(" [white::b]Music[#1DB954::b]Le[-::-]  Setup ")
	inner.SetTitleColor(ui.ColorPrimary)
	inner.AddItem(logo, 7, 0, false)
	inner.AddItem(dropdown, 1, 0, true)
	inner.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 0, 1, false)
	inner.AddItem(controls, 1, 0, true)

	page := centeredFlex(inner, 60, 18)
	page.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			state.Current.Language = sw.lang
			sw.buildStep2()
			return nil
		}
		return event
	})
	sw.root.AddPage("step1", page, true, true)
}

// ── Step 2 — Directory Selection ─────────────────────────────────────────────

func (sw *SetupWizard) buildStep2() {
	langT := func(en, tr string) string { return state.T(sw.lang, en, tr) }

	dirInput := makeInput("  "+langT("Music Directory", "Müzik Dizini")+":  ", langT("e.g. C:\\Music", "Örn: C:\\Müzik"))
	errView := makeErrorView()
	hint := hintBar("  [Enter] " + langT("Next", "İleri") + "  |  [F2] " + langT("Browse", "Gözat") + "  |  [Esc] " + langT("Back", "Geri"))

	// Folder browser function
	browseFolder := func() {
		// Open folder selection dialog
		selectedPath, err := dialog.Directory().Title(langT("Select Music Directory", "Müzik Dizini Seç")).Browse()
		if err != nil {
			errView.SetText("[#FF4444]  ✗ " + langT("Failed to open folder browser", "Klasör gezgini açılamadı") + "[-]")
			return
		}
		if selectedPath != "" {
			dirInput.SetText(selectedPath)
			errView.SetText("")
		}
	}

	onNext := func() {
		val := strings.TrimSpace(dirInput.GetText())
		if val == "" {
			errView.SetText("[#FF4444]  ✗ " + langT("Please select a directory", "Lütfen bir dizin seçin") + "[-]")
			return
		}
		_ = state.Current.InitializeBaseDirs(val) // Early initialization
		sw.rootDir = val
		sw.buildStep3()
	}
	nextBtn := makeNextButton(onNext)

	// Browse button
	browseBtn := tview.NewButton(langT("Browse", "Gözat"))
	browseBtn.SetBackgroundColor(tcell.NewRGBColor(100, 100, 100))
	browseBtn.SetLabelColor(tcell.ColorWhite)
	browseBtn.SetSelectedFunc(browseFolder)

	controls := tview.NewFlex().SetDirection(tview.FlexColumn)
	controls.SetBackgroundColor(ui.ColorBackground)
	controls.AddItem(hint, 0, 1, false)
	controls.AddItem(browseBtn, 10, 0, false)
	controls.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	controls.AddItem(nextBtn, 6, 0, true)

	inner := tview.NewFlex().SetDirection(tview.FlexRow)
	inner.SetBackgroundColor(ui.ColorBackground)
	inner.SetBorder(true)
	inner.SetBorderColor(ui.ColorAccent)
	inner.SetTitle(" " + langT("Directory Setup", "Dizin Seçimi") + " ")
	inner.SetTitleColor(ui.ColorPrimary)

	titleBox := tview.NewTextView()
	titleBox.SetDynamicColors(true)
	titleBox.SetText("\n  " + logoText(sw.lang) + "\n\n  [#B3B3B3]" + langT("Where should MusicLe store your music?", "MusicLe müziği nereye kaydetsin?") + "[-]")
	titleBox.SetBackgroundColor(ui.ColorBackground)

	inner.AddItem(titleBox, 6, 0, false)
	inner.AddItem(dirInput, 1, 0, true)
	inner.AddItem(errView, 1, 0, false)
	inner.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 0, 1, false)
	inner.AddItem(controls, 1, 0, true)

	page := centeredFlex(inner, 70, 15)
	errView.SetText("")

	page.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			onNext()
			return nil
		case tcell.KeyF2:
			browseFolder()
			return nil
		case tcell.KeyEsc:
			sw.root.SwitchToPage("step1")
			sw.app.SetFocus(sw.root)
			return nil
		}
		return event
	})

	sw.root.AddPage("step2", page, true, true)
	sw.root.SwitchToPage("step2")
	sw.app.SetFocus(dirInput)
}

// ── Step 3 — Profile Setup ────────────────────────────────────────────────────

func (sw *SetupWizard) buildStep3() {
	langT := func(en, tr string) string { return state.T(sw.lang, en, tr) }

	folderInput := makeInput("  "+langT("Folder Name", "Klasör Adı")+":       ", "myprofile")
	avatarInput := makeInput("  "+langT("Profile Picture path", "Profil Fotoğrafı")+": ", langT("optional", "isteğe bağlı"))
	nameInput := makeInput("  "+langT("Display Name", "Görünen Ad")+":       ", "MusicLe User")
	bioInput := makeInput("  "+langT("Bio", "Biyografi")+":              ", langT("Music lover", "Müzik sever"))
	errView := makeErrorView()
	hint := hintBar("  [Enter] " + langT("Next", "İleri") + "  |  [Esc] " + langT("Back", "Geri"))

	onNext := func() {
		folder := strings.TrimSpace(folderInput.GetText())
		if folder == "" {
			errView.SetText("[#FF4444]  ✗ " + langT("Folder name required", "Klasör adı gerekli") + "[-]")
			return
		}
		sw.profile.folder = folder
		sw.profile.avatar = strings.TrimSpace(avatarInput.GetText())
		sw.profile.name = strings.TrimSpace(nameInput.GetText())
		if sw.profile.name == "" {
			sw.profile.name = folder
		}
		sw.profile.bio = strings.TrimSpace(bioInput.GetText())
		sw.buildStep4()
	}
	nextBtn := makeNextButton(onNext)

	controls := tview.NewFlex().SetDirection(tview.FlexColumn)
	controls.SetBackgroundColor(ui.ColorBackground)
	controls.AddItem(hint, 0, 1, false)
	controls.AddItem(nextBtn, 6, 0, true)

	inner := tview.NewFlex().SetDirection(tview.FlexRow)
	inner.SetBackgroundColor(ui.ColorBackground)
	inner.SetBorder(true)
	inner.SetBorderColor(ui.ColorAccent)
	inner.SetTitle(" " + langT("Profile Setup", "Profil Oluştur") + " ")
	inner.SetTitleColor(ui.ColorPrimary)

	titleBox := tview.NewTextView()
	titleBox.SetDynamicColors(true)
	titleBox.SetText("\n  " + logoText(sw.lang) + "  [#B3B3B3]— " + langT("Step 3: Create your profile", "Adım 3: Profilinizi oluşturun") + "[-]")
	titleBox.SetBackgroundColor(ui.ColorBackground)

	inner.AddItem(titleBox, 3, 0, false)
	inner.AddItem(folderInput, 1, 0, true)
	inner.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	inner.AddItem(avatarInput, 1, 0, false)
	inner.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	inner.AddItem(nameInput, 1, 0, false)
	inner.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	inner.AddItem(bioInput, 1, 0, false)
	inner.AddItem(errView, 1, 0, false)
	inner.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 0, 1, false)
	inner.AddItem(controls, 1, 0, true)

	inputs := []*tview.InputField{folderInput, avatarInput, nameInput, bioInput}
	focusIdx := 0

	page := centeredFlex(inner, 72, 22)
	page.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if event.Modifiers()&tcell.ModCtrl == 0 { // Real Tab, not Ctrl+I
				focusIdx = (focusIdx + 1) % len(inputs)
				sw.app.SetFocus(inputs[focusIdx])
				return nil
			}
		case tcell.KeyBacktab:
			focusIdx = (focusIdx - 1 + len(inputs)) % len(inputs)
			sw.app.SetFocus(inputs[focusIdx])
			return nil
		case tcell.KeyEnter:
			onNext()
			return nil
		case tcell.KeyEsc:
			sw.root.SwitchToPage("step2")
			sw.app.SetFocus(sw.root)
			return nil
		}
		return event
	})

	sw.root.AddPage("step3", page, true, true)
	sw.root.SwitchToPage("step3")
	sw.app.SetFocus(folderInput)
}

// ── Step 4 — Initial Playlist Setup ──────────────────────────────────────────

func (sw *SetupWizard) buildStep4() {
	langT := func(en, tr string) string { return state.T(sw.lang, en, tr) }

	folderInput := makeInput("  "+langT("Folder Name", "Klasör Adı")+":        ", "my-playlist")
	artInput := makeInput("  "+langT("Playlist Art path", "Playlist Görseli")+":  ", langT("optional", "isteğe bağlı"))
	nameInput := makeInput("  "+langT("Playlist Name", "Playlist Adı")+":     ", langT("My Playlist", "Listem"))
	descInput := makeInput("  "+langT("Description", "Açıklama")+":         ", langT("My favorite songs", "Favori şarkılarım"))
	errView := makeErrorView()
	hint := hintBar("  [Enter] " + langT("Finish", "Bitir") + "  |  [Esc] " + langT("Back", "Geri"))

	onNext := func() {
		folder := strings.TrimSpace(folderInput.GetText())
		if folder == "" {
			errView.SetText("[#FF4444]  ✗ " + langT("Folder name required", "Klasör adı gerekli") + "[-]")
			return
		}
		sw.playlist.folder = folder
		sw.playlist.art = strings.TrimSpace(artInput.GetText())
		sw.playlist.name = strings.TrimSpace(nameInput.GetText())
		if sw.playlist.name == "" {
			sw.playlist.name = folder
		}
		sw.playlist.bio = strings.TrimSpace(descInput.GetText())
		sw.finishSetup()
	}
	nextBtn := makeNextButton(onNext)

	controls := tview.NewFlex().SetDirection(tview.FlexColumn)
	controls.SetBackgroundColor(ui.ColorBackground)
	controls.AddItem(hint, 0, 1, false)
	controls.AddItem(nextBtn, 6, 0, true)

	inner := tview.NewFlex().SetDirection(tview.FlexRow)
	inner.SetBackgroundColor(ui.ColorBackground)
	inner.SetBorder(true)
	inner.SetBorderColor(ui.ColorAccent)
	inner.SetTitle(" " + langT("Playlist Setup", "İlk Playlist") + " ")
	inner.SetTitleColor(ui.ColorPrimary)

	titleBox := tview.NewTextView()
	titleBox.SetDynamicColors(true)
	titleBox.SetText("\n  " + logoText(sw.lang) + "  [#B3B3B3]— " + langT("Step 4: Create your first playlist", "Adım 4: İlk playlistinizi oluşturun") + "[-]")
	titleBox.SetBackgroundColor(ui.ColorBackground)

	inner.AddItem(titleBox, 3, 0, false)
	inner.AddItem(folderInput, 1, 0, true)
	inner.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	inner.AddItem(artInput, 1, 0, false)
	inner.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	inner.AddItem(nameInput, 1, 0, false)
	inner.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 1, 0, false)
	inner.AddItem(descInput, 1, 0, false)
	inner.AddItem(errView, 1, 0, false)
	inner.AddItem(tview.NewBox().SetBackgroundColor(ui.ColorBackground), 0, 1, false)
	inner.AddItem(controls, 1, 0, true)

	inputs := []*tview.InputField{folderInput, artInput, nameInput, descInput}
	focusIdx := 0

	page := centeredFlex(inner, 72, 22)
	page.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if event.Modifiers()&tcell.ModCtrl == 0 { // Real Tab, not Ctrl+I
				focusIdx = (focusIdx + 1) % len(inputs)
				sw.app.SetFocus(inputs[focusIdx])
				return nil
			}
		case tcell.KeyBacktab:
			focusIdx = (focusIdx - 1 + len(inputs)) % len(inputs)
			sw.app.SetFocus(inputs[focusIdx])
			return nil
		case tcell.KeyEnter:
			onNext()
			return nil
		case tcell.KeyEsc:
			sw.root.SwitchToPage("step3")
			sw.app.SetFocus(sw.root)
			return nil
		}
		return event
	})

	sw.root.AddPage("step4", page, true, true)
	sw.root.SwitchToPage("step4")
	sw.app.SetFocus(folderInput)
}

// ── Finish: persist everything and switch to home ─────────────────────────────

func (sw *SetupWizard) finishSetup() {
	langT := func(en, tr string) string { return state.T(sw.lang, en, tr) }
	state.Current.RootDir = sw.rootDir
	state.Current.Language = sw.lang

	// Create profile directory structure
	if err := state.Current.CreateProfileStructure(
		sw.profile.folder, sw.profile.name, sw.profile.bio, sw.profile.avatar, sw.lang,
	); err != nil {
		sw.showError(langT("Failed to create profile: ", "Profil oluşturulamadı: ") + err.Error())
		return
	}

	// Create playlist directory structure
	if err := state.Current.CreatePlaylistStructure(
		sw.profile.folder, sw.playlist.folder, sw.playlist.name, sw.playlist.bio, sw.playlist.art,
	); err != nil {
		sw.showError(langT("Failed to create playlist: ", "Playlist oluşturulamadı: ") + err.Error())
		return
	}

	// Persist config
	_ = state.Current.SaveConfig()

	// Reload profiles into state
	_ = state.Current.ScanProfiles()
	if len(state.Current.Profiles) > 0 {
		state.Current.CurrentProfile = &state.Current.Profiles[0]
		if len(state.Current.CurrentProfile.Playlists) > 0 {
			state.Current.CurrentPlaylist = &state.Current.CurrentProfile.Playlists[0]
		}
	}
	state.Current.IsFirstLaunch = false

	// Switch to home page — must be done in the UI goroutine
	sw.app.QueueUpdateDraw(func() {
		sw.pages.SwitchToPage("home")
	})
}

func (sw *SetupWizard) showError(msg string) {
	modal := tview.NewModal().
		SetText("⚠  " + msg).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			sw.pages.RemovePage("setupErr")
		})
	modal.SetBackgroundColor(tcell.NewRGBColor(30, 30, 30))
	modal.SetBorderColor(ui.ColorError)
	sw.pages.AddPage("setupErr", modal, false, true)
	sw.app.SetFocus(modal)
}
