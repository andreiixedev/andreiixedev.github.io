package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Change struct {
	Type        string `json:"type"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

type Entry struct {
	Version string   `json:"version"`
	Date    string   `json:"date"`
	Title   string   `json:"title"`
	Icon    string   `json:"icon"`
	Changes []Change `json:"changes"`
}

type Stats struct {
	TotalUpdates  int `json:"totalUpdates"`
	FeaturesAdded int `json:"featuresAdded"`
	BugsFixed     int `json:"bugsFixed"`
	Improvements  int `json:"improvements"`
}

type FutureFeature struct {
	Icon string `json:"icon"`
	Name string `json:"name"`
}

type Changelog struct {
	Stats          Stats           `json:"stats"`
	Entries        []Entry         `json:"entries"`
	FutureFeatures []FutureFeature `json:"futureFeatures"`
}

type ProjectLinks struct {
	Github  string `json:"github,omitempty"`
	Live    string `json:"live,omitempty"`
	Release string `json:"release,omitempty"`
}

type Project struct {
	ID           int          `json:"id"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	Image        string       `json:"image"`
	Year         string       `json:"year"`
	Technologies []string     `json:"technologies"`
	Category     string       `json:"category"`
	Featured     bool         `json:"featured"`
	Links        ProjectLinks `json:"links"`
}

type ProjectsSettings struct {
	ShowFeaturedOnly bool `json:"showFeaturedOnly"`
	ItemsPerPage     int  `json:"itemsPerPage"`
}

type ProjectsData struct {
	Projects   []Project        `json:"projects"`
	Categories []string         `json:"categories"`
	Settings   ProjectsSettings `json:"settings"`
}

const (
	changelogFile = "changelog.json"
	projectsFile  = "projects.json"
)

var (
	currentData          Changelog
	projectsData         ProjectsData
	mainWindow           fyne.Window
	selectedProjectIndex = -1
	isDarkMode           = false

	// Neumorphic from index.html 
	bgColor      = color.NRGBA{R: 247, G: 242, B: 233, A: 255} // --bg-primary
	surfaceColor = color.NRGBA{R: 240, G: 232, B: 220, A: 255} // --bg-secondary
	cardColor    = color.NRGBA{R: 229, G: 217, B: 204, A: 255} // --bg-tertiary
	borderColor  = color.NRGBA{R: 214, G: 201, B: 184, A: 255}
	accentColor  = color.NRGBA{R: 184, G: 123, B: 93, A: 255}  // --accent
	mutedText    = color.NRGBA{R: 158, G: 125, B: 107, A: 255} // --text-muted
	primaryText  = color.NRGBA{R: 74, G: 74, B: 74, A: 255}    // --text-primary
	dangerColor  = color.NRGBA{R: 180, G: 78, B: 70, A: 255}
	successColor = color.NRGBA{R: 108, G: 145, B: 95, A: 255}
	warningColor = color.NRGBA{R: 194, G: 143, B: 79, A: 255}
)

type modernTheme struct{}

func (modernTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return bgColor
	case theme.ColorNameButton:
		return surfaceColor
	case theme.ColorNameDisabledButton:
		return cardColor
	case theme.ColorNameForeground:
		return primaryText
	case theme.ColorNameInputBackground:
		return surfaceColor
	case theme.ColorNamePlaceHolder:
		return mutedText
	case theme.ColorNamePrimary:
		return accentColor
	case theme.ColorNameHover:
		return color.NRGBA{R: 238, G: 225, B: 211, A: 255}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 218, G: 203, B: 188, A: 255}
	case theme.ColorNameSeparator:
		return borderColor
	case theme.ColorNameShadow:
		return color.NRGBA{R: 120, G: 96, B: 78, A: 90}
	case theme.ColorNameFocus:
		return accentColor
	case theme.ColorNameSelection:
		return color.NRGBA{R: 184, G: 123, B: 93, A: 80}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func applyLightPalette() {
	bgColor = color.NRGBA{R: 247, G: 242, B: 233, A: 255}
	surfaceColor = color.NRGBA{R: 240, G: 232, B: 220, A: 255}
	cardColor = color.NRGBA{R: 229, G: 217, B: 204, A: 255}
	borderColor = color.NRGBA{R: 214, G: 201, B: 184, A: 255}
	accentColor = color.NRGBA{R: 184, G: 123, B: 93, A: 255}
	mutedText = color.NRGBA{R: 158, G: 125, B: 107, A: 255}
	primaryText = color.NRGBA{R: 74, G: 74, B: 74, A: 255}
	dangerColor = color.NRGBA{R: 180, G: 78, B: 70, A: 255}
	successColor = color.NRGBA{R: 108, G: 145, B: 95, A: 255}
	warningColor = color.NRGBA{R: 194, G: 143, B: 79, A: 255}
}

func applyDarkPalette() {
	bgColor = color.NRGBA{R: 26, G: 26, B: 26, A: 255}
	surfaceColor = color.NRGBA{R: 45, G: 45, B: 45, A: 255}
	cardColor = color.NRGBA{R: 61, G: 61, B: 61, A: 255}
	borderColor = color.NRGBA{R: 82, G: 82, B: 82, A: 255}
	accentColor = color.NRGBA{R: 184, G: 123, B: 93, A: 255}
	mutedText = color.NRGBA{R: 138, G: 138, B: 138, A: 255}
	primaryText = color.NRGBA{R: 224, G: 224, B: 224, A: 255}
	dangerColor = color.NRGBA{R: 230, G: 100, B: 92, A: 255}
	successColor = color.NRGBA{R: 120, G: 190, B: 125, A: 255}
	warningColor = color.NRGBA{R: 220, G: 165, B: 95, A: 255}
}

func applySelectedPalette() {
	if isDarkMode {
		applyDarkPalette()
	} else {
		applyLightPalette()
	}
	fyne.CurrentApp().Settings().SetTheme(modernTheme{})
}

func toggleTheme() {
	isDarkMode = !isDarkMode
	applySelectedPalette()
	mainWindow.SetContent(createMainMenu())
}

func (modernTheme) Font(style fyne.TextStyle) fyne.Resource { return theme.DefaultTheme().Font(style) }
func (modernTheme) Icon(name fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(name) }
func (modernTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 10
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameText:
		return 15
	case theme.SizeNameHeadingText:
		return 24
	}
	return theme.DefaultTheme().Size(name)
}

func main() {
	myApp := app.New()
	applySelectedPalette()

	mainWindow = myApp.NewWindow("AndreiixeDev Studio")
	mainWindow.Resize(fyne.NewSize(760, 560))
	mainWindow.CenterOnScreen()

	loadChangelogData()
	loadProjectsData()

	mainWindow.SetContent(createMainMenu())
	mainWindow.ShowAndRun()
}

func createMainMenu() fyne.CanvasObject {
	intro := widget.NewLabelWithStyle(
		"Choose what you want to manage",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	menu := container.NewGridWithColumns(2,
		bigMenuButton("📊", "Dashboard", "Stats and quick overview", func() {
			openToolWindow("Dashboard", createStatsTab(), 980, 680)
		}),
		bigMenuButton("📝", "Changelog", "Add and manage releases", func() {
			openToolWindow("Changelog", createEntriesTab(), 1180, 760)
		}),
		bigMenuButton("🔮", "Roadmap", "Future features", func() {
			openToolWindow("Roadmap", createFutureTab(), 900, 620)
		}),
		bigMenuButton("💻", "Projects", "Portfolio projects", func() {
			openToolWindow("Projects", createProjectsTab(), 1220, 780)
		}),
		bigMenuButton("⚙️", "Settings", "Files and system actions", func() {
			openToolWindow("Settings", createSaveTab(), 860, 620)
		}),
		bigMenuButton("💾", "Save All", "Save changelog and projects", saveAllData),
	)

	return container.NewBorder(
		container.NewVBox(createHeader(), separator()),
		nil,
		nil,
		nil,
		container.NewMax(
			canvas.NewRectangle(bgColor),
			container.NewPadded(container.NewVBox(intro, menu)),
		),
	)
}

func bigMenuButton(icon, title, subtitle string, tapped func()) fyne.CanvasObject {
	iconText := canvas.NewText(icon, accentColor)
	iconText.TextSize = 38
	iconText.Alignment = fyne.TextAlignCenter

	titleText := canvas.NewText(title, primaryText)
	titleText.TextSize = 22
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.Alignment = fyne.TextAlignCenter

	subtitleText := canvas.NewText(subtitle, mutedText)
	subtitleText.TextSize = 14
	subtitleText.Alignment = fyne.TextAlignCenter

	btn := widget.NewButton("Open", tapped)
	btn.Importance = widget.HighImportance

	return panel("", container.NewCenter(container.NewVBox(iconText, titleText, subtitleText, softPill(btn))))
}

func openToolWindow(title string, content fyne.CanvasObject, width, height float32) {
	win := fyne.CurrentApp().NewWindow("AndreiixeDev Studio • " + title)
	win.Resize(fyne.NewSize(width, height))
	win.CenterOnScreen()

	win.SetContent(container.NewBorder(
		container.NewVBox(
			createWindowHeader(title),
			container.NewPadded(container.NewHBox(
				widget.NewButtonWithIcon("Save all", theme.DocumentSaveIcon(), saveAllData),
				widget.NewButtonWithIcon("Reload", theme.ViewRefreshIcon(), func() {
					loadChangelogData()
					loadProjectsData()
					win.SetContent(container.NewBorder(
						container.NewVBox(createWindowHeader(title), separator()),
						nil, nil, nil,
						container.NewMax(canvas.NewRectangle(bgColor), container.NewPadded(content)),
					))
				}),
				layout.NewSpacer(),
				widget.NewButtonWithIcon("Close", theme.CancelIcon(), func() { win.Close() }),
			)),
			separator(),
		),
		nil,
		nil,
		nil,
		container.NewMax(canvas.NewRectangle(bgColor), container.NewPadded(content)),
	))

	win.Show()
}

func createWindowHeader(title string) fyne.CanvasObject {
	logo := canvas.NewText("AndreiixeDev.", accentColor)
	logo.TextSize = 22
	logo.TextStyle = fyne.TextStyle{Bold: true}

	label := canvas.NewText(title, primaryText)
	label.TextSize = 26
	label.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Workspace window", mutedText)
	subtitle.TextSize = 13

	return container.NewPadded(container.NewVBox(logo, label, subtitle))
}

func openMessageWindow(title, message string) {
	win := fyne.CurrentApp().NewWindow(title)
	win.Resize(fyne.NewSize(520, 240))
	win.CenterOnScreen()

	titleText := canvas.NewText(title, primaryText)
	titleText.TextSize = 22
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	messageLabel := widget.NewLabel(message)
	messageLabel.Wrapping = fyne.TextWrapWord

	okBtn := widget.NewButtonWithIcon("OK", theme.ConfirmIcon(), func() {
		win.Close()
	})
	okBtn.Importance = widget.HighImportance

	win.SetContent(container.NewMax(
		canvas.NewRectangle(bgColor),
		container.NewPadded(container.NewVBox(
			titleText,
			separator(),
			messageLabel,
			layout.NewSpacer(),
			container.NewHBox(layout.NewSpacer(), okBtn),
		)),
	))
	win.Show()
}

func openErrorWindow(err error) {
	if err == nil {
		return
	}
	openMessageWindow("Error", err.Error())
}

func openConfirmWindow(title, message, dangerActionText string, onConfirm func()) {
	win := fyne.CurrentApp().NewWindow("Confirm • " + title)
	win.Resize(fyne.NewSize(520, 260))
	win.CenterOnScreen()

	titleText := canvas.NewText(title, primaryText)
	titleText.TextSize = 22
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	messageLabel := widget.NewLabel(message)
	messageLabel.Wrapping = fyne.TextWrapWord

	cancelBtn := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		win.Close()
	})

	deleteBtn := widget.NewButtonWithIcon(dangerActionText, theme.DeleteIcon(), func() {
		if onConfirm != nil {
			onConfirm()
		}
		win.Close()
	})
	deleteBtn.Importance = widget.DangerImportance

	win.SetContent(container.NewMax(
		canvas.NewRectangle(bgColor),
		container.NewPadded(container.NewVBox(
			titleText,
			separator(),
			messageLabel,
			layout.NewSpacer(),
			container.NewGridWithColumns(2, cancelBtn, deleteBtn),
		)),
	))
	win.Show()
}

func createHeader() fyne.CanvasObject {
	logo := canvas.NewText("AndreiixeDev.", accentColor)
	logo.TextSize = 28
	logo.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Changelog • Projects • Links", mutedText)
	subtitle.TextSize = 14

	clockBadge := widget.NewLabelWithStyle(time.Now().Format("15:04:05"), fyne.TextAlignCenter, fyne.TextStyle{Monospace: true, Bold: true})
	clockCard := softPill(container.NewHBox(widget.NewIcon(theme.HistoryIcon()), clockBadge))

	themeSymbol := "☾"
	if isDarkMode {
		themeSymbol = "☀"
	}

	themeText := canvas.NewText(themeSymbol, accentColor)
	themeText.TextSize = 24
	themeText.Alignment = fyne.TextAlignCenter

	themeCircle := canvas.NewCircle(cardColor)
	themeCircle.StrokeColor = borderColor
	themeCircle.StrokeWidth = 1
	themeVisual := container.NewGridWrap(
		fyne.NewSize(54, 54),
		container.NewStack(
			themeCircle,
			container.NewCenter(themeText),
		),
	)

	themeTap := widget.NewButton("", toggleTheme)
	themeTap.Importance = widget.LowImportance

	themeButton := container.NewStack(themeVisual, themeTap)

	right := container.NewHBox(themeButton, clockCard)

	return container.NewPadded(container.NewBorder(nil, nil, container.NewVBox(logo, subtitle), right))
}

func createToolbar() fyne.CanvasObject {
	saveBtn := widget.NewButtonWithIcon("Save all", theme.DocumentSaveIcon(), saveAllData)
	saveBtn.Importance = widget.HighImportance

	reloadBtn := widget.NewButtonWithIcon("Reload", theme.ViewRefreshIcon(), func() {
		loadChangelogData()
		loadProjectsData()
		mainWindow.SetContent(createMainMenu())
		openMessageWindow("Message", "Action completed or needs attention.")
	})

	aboutBtn := widget.NewButtonWithIcon("About", theme.InfoIcon(), showAbout)

	return container.NewPadded(container.NewHBox(layout.NewSpacer(), saveBtn, reloadBtn, aboutBtn))
}

func separator() fyne.CanvasObject {
	line := canvas.NewRectangle(borderColor)
	line.SetMinSize(fyne.NewSize(1, 1))
	return line
}

func panel(title string, content fyne.CanvasObject) fyne.CanvasObject {
	if content == nil {
		content = widget.NewLabel("")
	}

	outer := canvas.NewRectangle(borderColor)
	outer.CornerRadius = 30

	inner := canvas.NewRectangle(surfaceColor)
	inner.CornerRadius = 28

	var body fyne.CanvasObject
	if strings.TrimSpace(title) == "" {
		body = content
	} else {
		header := canvas.NewText(title, accentColor)
		header.TextSize = 22
		header.TextStyle = fyne.TextStyle{Bold: true}
		body = container.NewVBox(header, dashedSeparator(), content)
	}

	return container.NewPadded(container.NewMax(
		outer,
		container.NewPadded(container.NewMax(inner, container.NewPadded(body))),
	))
}

func softPill(content fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(cardColor)
	bg.CornerRadius = 40
	return container.NewPadded(container.NewMax(bg, container.NewPadded(content)))
}

func dashedSeparator() fyne.CanvasObject {
	line := canvas.NewRectangle(mutedText)
	line.SetMinSize(fyne.NewSize(1, 2))
	return line
}

func showAbout() {
	dialog.ShowInformation("About", "AndreiixeDev Studio\n\n© 2026 AndreiixeDev", mainWindow)
}

func defaultChangelog() Changelog {
	return Changelog{Stats: Stats{}, Entries: []Entry{}, FutureFeatures: []FutureFeature{}}
}

func defaultProjectsData() ProjectsData {
	return ProjectsData{
		Projects:   []Project{},
		Categories: []string{"web", "mobile", "desktop", "game", "tool"},
		Settings: ProjectsSettings{
			ShowFeaturedOnly: false,
			ItemsPerPage:     6,
		},
	}
}

func loadChangelogData() {
	data, err := os.ReadFile(changelogFile)
	if err != nil {
		currentData = defaultChangelog()
		saveChangelogData()
		return
	}
	if err := json.Unmarshal(data, &currentData); err != nil {
		currentData = defaultChangelog()
		dialog.ShowError(fmt.Errorf("changelog.json is corrupted: %w", err), mainWindow)
		return
	}
	if currentData.Entries == nil {
		currentData.Entries = []Entry{}
	}
	if currentData.FutureFeatures == nil {
		currentData.FutureFeatures = []FutureFeature{}
	}
	recalculateStats()
}

func loadProjectsData() {
	data, err := os.ReadFile(projectsFile)
	if err != nil {
		projectsData = defaultProjectsData()
		saveProjectsData()
		return
	}
	if err := json.Unmarshal(data, &projectsData); err != nil {
		projectsData = defaultProjectsData()
		dialog.ShowError(fmt.Errorf("projects.json is corrupted: %w", err), mainWindow)
		return
	}
	if projectsData.Projects == nil {
		projectsData.Projects = []Project{}
	}
	if projectsData.Categories == nil || len(projectsData.Categories) == 0 {
		projectsData.Categories = defaultProjectsData().Categories
	}
	if projectsData.Settings.ItemsPerPage <= 0 {
		projectsData.Settings.ItemsPerPage = 6
	}
}

func saveChangelogData() {
	recalculateStats()
	data, err := json.MarshalIndent(currentData, "", "  ")
	if err != nil {
		openMessageWindow("Error", "An error occurred.")
		return
	}
	if err := os.WriteFile(changelogFile, data, 0644); err != nil {
		openMessageWindow("Error", "An error occurred.")
	}
}

func saveProjectsData() {
	data, err := json.MarshalIndent(projectsData, "", "  ")
	if err != nil {
		openMessageWindow("Error", "An error occurred.")
		return
	}
	if err := os.WriteFile(projectsFile, data, 0644); err != nil {
		openMessageWindow("Error", "An error occurred.")
	}
}

func saveAllData() {
	saveChangelogData()
	saveProjectsData()
	openMessageWindow("Message", "Action completed or needs attention.")
}

func recalculateStats() {
	stats := Stats{TotalUpdates: len(currentData.Entries)}
	for _, entry := range currentData.Entries {
		for _, change := range entry.Changes {
			switch strings.ToLower(strings.TrimSpace(change.Type)) {
			case "added":
				stats.FeaturesAdded++
			case "fixed":
				stats.BugsFixed++
			case "improved":
				stats.Improvements++
			}
		}
	}
	currentData.Stats = stats
}

func splitOption(value string) string {
	parts := strings.Fields(value)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func parseTechnologies(input string) []string {
	items := strings.Split(input, ",")
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

func nextProjectID() int {
	maxID := 0
	for _, p := range projectsData.Projects {
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	return maxID + 1
}

func clearProjectForm(fields ...*widget.Entry) {
	for _, field := range fields {
		field.SetText("")
	}
}

func createStatsTab() fyne.CanvasObject {
	cards := container.NewGridWithColumns(4,
		statCard("Total updates", currentData.Stats.TotalUpdates, accentColor),
		statCard("Features", currentData.Stats.FeaturesAdded, successColor),
		statCard("Bugs fixed", currentData.Stats.BugsFixed, dangerColor),
		statCard("Improvements", currentData.Stats.Improvements, warningColor),
	)

	recent := container.NewVBox()
	limit := 5
	if len(currentData.Entries) < limit {
		limit = len(currentData.Entries)
	}
	for i := 0; i < limit; i++ {
		e := currentData.Entries[i]
		recent.Add(widget.NewLabel(fmt.Sprintf("%s  %s — %s", e.Icon, e.Version, e.Title)))
	}
	if limit == 0 {
		recent.Add(widget.NewLabel("No releases yet."))
	}

	quickActions := container.NewGridWithColumns(3,
		widget.NewButtonWithIcon("Add changelog", theme.ContentAddIcon(), func() {
			openMessageWindow("Message", "Action completed or needs attention.")
		}),
		widget.NewButtonWithIcon("Add project", theme.ContentAddIcon(), func() {
			openMessageWindow("Message", "Action completed or needs attention.")
		}),
		widget.NewButtonWithIcon("Refresh stats", theme.ViewRefreshIcon(), func() {
			recalculateStats()
			mainWindow.SetContent(createMainMenu())
		}),
	)

	return container.NewVScroll(container.NewPadded(container.NewVBox(
		cards,
		panel("Recent releases", recent),
		panel("Quick actions", quickActions),
	)))
}

func statCard(title string, value int, c color.Color) fyne.CanvasObject {
	valueText := canvas.NewText(fmt.Sprintf("%d", value), c)
	valueText.TextSize = 34
	valueText.TextStyle = fyne.TextStyle{Bold: true}
	label := canvas.NewText(title, mutedText)
	label.TextSize = 14
	return panel("", container.NewCenter(container.NewVBox(valueText, label)))
}

func createEntriesTab() fyne.CanvasObject {
	entryList := widget.NewList(
		func() int { return len(currentData.Entries) },
		func() fyne.CanvasObject {
			icon := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
			version := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			title := widget.NewLabel("")
			date := widget.NewLabel("")
			return container.NewBorder(nil, nil, icon, date, container.NewVBox(version, title))
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			if i < 0 || i >= len(currentData.Entries) {
				return
			}
			e := currentData.Entries[i]
			box := obj.(*fyne.Container)
			icon := box.Objects[1].(*widget.Label)
			date := box.Objects[2].(*widget.Label)
			center := box.Objects[0].(*fyne.Container)
			version := center.Objects[0].(*widget.Label)
			title := center.Objects[1].(*widget.Label)

			if e.Icon == "" {
				e.Icon = "📝"
			}
			icon.SetText(e.Icon)
			version.SetText(e.Version)
			title.SetText(e.Title)
			date.SetText(e.Date)
		},
	)
	entryList.SetItemHeight(0, 70)

	versionEntry := widget.NewEntry()
	versionEntry.SetPlaceHolder("v2.0.0")
	dateEntry := widget.NewEntry()
	dateEntry.SetText(time.Now().Format("2006-01-02"))
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Release title")
	iconSelect := widget.NewSelect([]string{"⭐ star", "🚀 rocket", "🎨 paint", "🔗 link", "📊 chart", "🎮 game", "📱 mobile", "🛠️ tool"}, nil)
	iconSelect.SetSelected("🚀 rocket")

	changeType := widget.NewSelect([]string{"added", "fixed", "improved", "changed", "removed"}, nil)
	changeType.SetSelected("added")
	changeIcon := widget.NewEntry()
	changeIcon.SetPlaceHolder("🐛 / 🚀 / 🛠️")
	changeDesc := widget.NewEntry()
	changeDesc.SetPlaceHolder("Describe the change")

	changes := []Change{}
	changesBox := container.NewVBox(widget.NewLabel("No changes added yet."))
	changesCount := widget.NewLabel("0 changes")
	refreshChanges := func() {
		changesBox.Objects = nil
		if len(changes) == 0 {
			changesBox.Add(widget.NewLabel("No changes added yet."))
		} else {
			for i, change := range changes {
				changesBox.Add(widget.NewLabel(fmt.Sprintf("%d. %s [%s] %s", i+1, change.Icon, change.Type, change.Description)))
			}
		}
		changesCount.SetText(fmt.Sprintf("%d changes", len(changes)))
		changesBox.Refresh()
	}

	addChangeBtn := widget.NewButtonWithIcon("Add change", theme.ContentAddIcon(), func() {
		if strings.TrimSpace(changeDesc.Text) == "" {
			openMessageWindow("Message", "Action completed or needs attention.")
			return
		}
		icon := strings.TrimSpace(changeIcon.Text)
		if icon == "" {
			icon = "•"
		}
		changes = append(changes, Change{Type: changeType.Selected, Icon: icon, Description: strings.TrimSpace(changeDesc.Text)})
		changeDesc.SetText("")
		changeIcon.SetText("")
		refreshChanges()
	})

	addEntryBtn := widget.NewButtonWithIcon("Save release", theme.ConfirmIcon(), func() {
		if strings.TrimSpace(versionEntry.Text) == "" || strings.TrimSpace(titleEntry.Text) == "" {
			openMessageWindow("Message", "Action completed or needs attention.")
			return
		}
		entry := Entry{
			Version: strings.TrimSpace(versionEntry.Text),
			Date:    strings.TrimSpace(dateEntry.Text),
			Title:   strings.TrimSpace(titleEntry.Text),
			Icon:    splitOption(iconSelect.Selected),
			Changes: changes,
		}
		currentData.Entries = append([]Entry{entry}, currentData.Entries...)
		saveChangelogData()
		entryList.Refresh()
		versionEntry.SetText("")
		titleEntry.SetText("")
		changes = []Change{}
		refreshChanges()
		openMessageWindow("Message", "Action completed or needs attention.")
	})
	addEntryBtn.Importance = widget.HighImportance

	deleteBtn := widget.NewButtonWithIcon("Delete newest", theme.DeleteIcon(), func() {
		if len(currentData.Entries) == 0 {
			return
		}
		dialog.ShowConfirm("Delete newest release", "This cannot be undone.", func(ok bool) {
			if !ok {
				return
			}
			currentData.Entries = currentData.Entries[1:]
			saveChangelogData()
			entryList.Refresh()
		}, mainWindow)
	})

	form := container.NewVScroll(container.NewPadded(container.NewVBox(
		panel("New release", container.NewVBox(
			container.NewGridWithColumns(2, versionEntry, dateEntry),
			titleEntry,
			iconSelect,
		)),
		panel("Changes", container.NewVBox(
			container.NewGridWithColumns(3, changeType, changeIcon, changeDesc),
			container.NewHBox(addChangeBtn, changesCount),
			changesBox,
		)),
		panel("Actions", container.NewHBox(addEntryBtn, deleteBtn)),
	)))

	split := container.NewHSplit(panel("Version history", entryList), form)
	split.Offset = 0.38
	return split
}

func createFutureTab() fyne.CanvasObject {
	featureList := widget.NewList(
		func() int { return len(currentData.FutureFeatures) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			if i < 0 || i >= len(currentData.FutureFeatures) {
				return
			}
			f := currentData.FutureFeatures[i]
			if f.Icon == "" {
				f.Icon = "🔮"
			}
			obj.(*widget.Label).SetText(fmt.Sprintf("%s  %s", f.Icon, f.Name))
		},
	)

	iconSelect := widget.NewSelect([]string{"💬 chat", "📈 analytics", "🔍 search", "🔄 sync", "⭐ star", "🔔 notify", "✉️ email", "🎮 game", "📱 mobile", "🌐 web"}, nil)
	iconSelect.SetSelected("🔮 roadmap")
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Future feature name")

	addBtn := widget.NewButtonWithIcon("Add feature", theme.ContentAddIcon(), func() {
		name := strings.TrimSpace(nameEntry.Text)
		if name == "" {
			openMessageWindow("Message", "Action completed or needs attention.")
			return
		}
		currentData.FutureFeatures = append(currentData.FutureFeatures, FutureFeature{Icon: splitOption(iconSelect.Selected), Name: name})
		nameEntry.SetText("")
		saveChangelogData()
		featureList.Refresh()
	})
	addBtn.Importance = widget.HighImportance

	deleteBtn := widget.NewButtonWithIcon("Remove last", theme.DeleteIcon(), func() {
		if len(currentData.FutureFeatures) == 0 {
			return
		}
		currentData.FutureFeatures = currentData.FutureFeatures[:len(currentData.FutureFeatures)-1]
		saveChangelogData()
		featureList.Refresh()
	})

	form := panel("Add roadmap item", container.NewVBox(iconSelect, nameEntry, container.NewHBox(addBtn, deleteBtn)))
	split := container.NewHSplit(panel("Planned features", featureList), container.NewPadded(form))
	split.Offset = 0.45
	return split
}

func createProjectsTab() fyne.CanvasObject {
	selectedProjectIndex = -1

	var refreshProjectList func()

	projectList := widget.NewList(
		func() int { return len(projectsData.Projects) },
		func() fyne.CanvasObject {
			icon := widget.NewIcon(theme.ComputerIcon())

			title := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			title.Wrapping = fyne.TextWrapWord

			category := widget.NewLabel("")
			category.Wrapping = fyne.TextWrapWord

			meta := widget.NewLabel("")
			meta.Wrapping = fyne.TextWrapWord

			content := container.NewVBox(title, category, meta)
			return container.NewPadded(container.NewBorder(nil, nil, icon, nil, content))
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			if i < 0 || i >= len(projectsData.Projects) {
				return
			}

			p := projectsData.Projects[i]
			padded := obj.(*fyne.Container)
			box := padded.Objects[0].(*fyne.Container)
			content := box.Objects[0].(*fyne.Container)

			title := content.Objects[0].(*widget.Label)
			category := content.Objects[1].(*widget.Label)
			meta := content.Objects[2].(*widget.Label)

			star := ""
			if p.Featured {
				star = " ⭐"
			}

			title.SetText(p.Title + star)
			category.SetText("Category: " + p.Category)
			meta.SetText(fmt.Sprintf("Year: %s  •  ID: #%d", p.Year, p.ID))
		},
	)

	refreshProjectList = func() {
		for i := range projectsData.Projects {
			projectList.SetItemHeight(i, 104)
		}
		projectList.Refresh()
	}

	projectList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(projectsData.Projects) {
			selectedProjectIndex = id
		}
	}
	refreshProjectList()

	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Project title")

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Project description")
	descEntry.SetMinRowsVisible(4)

	imageEntry := widget.NewEntry()
	imageEntry.SetPlaceHolder("images/project.jpg")

	yearEntry := widget.NewEntry()
	yearEntry.SetPlaceHolder("2026")

	techEntry := widget.NewEntry()
	techEntry.SetPlaceHolder("Go, Fyne, React")

	categorySelect := widget.NewSelect(projectsData.Categories, nil)
	if len(projectsData.Categories) > 0 {
		categorySelect.SetSelected(projectsData.Categories[0])
	}

	featuredCheck := widget.NewCheck("Featured project", nil)

	githubEntry := widget.NewEntry()
	githubEntry.SetPlaceHolder("GitHub URL")

	liveEntry := widget.NewEntry()
	liveEntry.SetPlaceHolder("Live URL")

	releaseEntry := widget.NewEntry()
	releaseEntry.SetPlaceHolder("Release URL")

	loadForm := func(p Project) {
		titleEntry.SetText(p.Title)
		descEntry.SetText(p.Description)
		imageEntry.SetText(p.Image)
		yearEntry.SetText(p.Year)
		techEntry.SetText(strings.Join(p.Technologies, ", "))
		categorySelect.SetSelected(p.Category)
		featuredCheck.SetChecked(p.Featured)
		githubEntry.SetText(p.Links.Github)
		liveEntry.SetText(p.Links.Live)
		releaseEntry.SetText(p.Links.Release)
	}

	buildProject := func(id int) Project {
		return Project{
			ID:           id,
			Title:        strings.TrimSpace(titleEntry.Text),
			Description:  strings.TrimSpace(descEntry.Text),
			Image:        strings.TrimSpace(imageEntry.Text),
			Year:         strings.TrimSpace(yearEntry.Text),
			Technologies: parseTechnologies(techEntry.Text),
			Category:     categorySelect.Selected,
			Featured:     featuredCheck.Checked,
			Links: ProjectLinks{
				Github:  strings.TrimSpace(githubEntry.Text),
				Live:    strings.TrimSpace(liveEntry.Text),
				Release: strings.TrimSpace(releaseEntry.Text),
			},
		}
	}

	validateProject := func() bool {
		if strings.TrimSpace(titleEntry.Text) == "" || strings.TrimSpace(descEntry.Text) == "" {
			openMessageWindow("Message", "Action completed or needs attention.")
			return false
		}
		if categorySelect.Selected == "" {
			openMessageWindow("Message", "Action completed or needs attention.")
			return false
		}
		return true
	}

	clearForm := func() {
		clearProjectForm(titleEntry, descEntry, imageEntry, yearEntry, techEntry, githubEntry, liveEntry, releaseEntry)
		featuredCheck.SetChecked(false)
		if len(projectsData.Categories) > 0 {
			categorySelect.SetSelected(projectsData.Categories[0])
		}
	}

	loadBtn := widget.NewButtonWithIcon("Load selected", theme.DocumentCreateIcon(), func() {
		if selectedProjectIndex < 0 || selectedProjectIndex >= len(projectsData.Projects) {
			openMessageWindow("Message", "Action completed or needs attention.")
			return
		}
		loadForm(projectsData.Projects[selectedProjectIndex])
	})

	addBtn := widget.NewButtonWithIcon("Add project", theme.ContentAddIcon(), func() {
		if !validateProject() {
			return
		}
		projectsData.Projects = append(projectsData.Projects, buildProject(nextProjectID()))
		saveProjectsData()
		refreshProjectList()
		clearForm()
	})
	addBtn.Importance = widget.HighImportance

	updateBtn := widget.NewButtonWithIcon("Update selected", theme.ConfirmIcon(), func() {
		if selectedProjectIndex < 0 || selectedProjectIndex >= len(projectsData.Projects) {
			openMessageWindow("Message", "Action completed or needs attention.")
			return
		}
		if !validateProject() {
			return
		}
		id := projectsData.Projects[selectedProjectIndex].ID
		projectsData.Projects[selectedProjectIndex] = buildProject(id)
		saveProjectsData()
		refreshProjectList()
		clearForm()
	})

	moveUpBtn := widget.NewButtonWithIcon("Move up", theme.MoveUpIcon(), func() {
		if selectedProjectIndex <= 0 || selectedProjectIndex >= len(projectsData.Projects) {
			return
		}
		projectsData.Projects[selectedProjectIndex], projectsData.Projects[selectedProjectIndex-1] = projectsData.Projects[selectedProjectIndex-1], projectsData.Projects[selectedProjectIndex]
		selectedProjectIndex--
		saveProjectsData()
		refreshProjectList()
		projectList.Select(selectedProjectIndex)
	})

	moveDownBtn := widget.NewButtonWithIcon("Move down", theme.MoveDownIcon(), func() {
		if selectedProjectIndex < 0 || selectedProjectIndex >= len(projectsData.Projects)-1 {
			return
		}
		projectsData.Projects[selectedProjectIndex], projectsData.Projects[selectedProjectIndex+1] = projectsData.Projects[selectedProjectIndex+1], projectsData.Projects[selectedProjectIndex]
		selectedProjectIndex++
		saveProjectsData()
		refreshProjectList()
		projectList.Select(selectedProjectIndex)
	})

	categoryEntry := widget.NewEntry()
	categoryEntry.SetPlaceHolder("new category")

	categoriesBox := container.NewVBox()
	categoriesCount := widget.NewLabel("")

	refreshCategories := func() {
		sort.Strings(projectsData.Categories)
		categoriesBox.Objects = nil
		for _, category := range projectsData.Categories {
			categoriesBox.Add(widget.NewLabel("🏷️ " + category))
		}
		categoriesCount.SetText(fmt.Sprintf("%d categories", len(projectsData.Categories)))
		categorySelect.Options = projectsData.Categories
		categorySelect.Refresh()
		categoriesBox.Refresh()
	}
	refreshCategories()

	addCategoryBtn := widget.NewButtonWithIcon("Add category", theme.ContentAddIcon(), func() {
		newCategory := strings.TrimSpace(categoryEntry.Text)
		if newCategory == "" {
			openMessageWindow("Message", "Action completed or needs attention.")
			return
		}
		for _, category := range projectsData.Categories {
			if strings.EqualFold(category, newCategory) {
				dialog.ShowInformation("Already exists", fmt.Sprintf("%q already exists.", category), mainWindow)
				return
			}
		}
		projectsData.Categories = append(projectsData.Categories, newCategory)
		categoryEntry.SetText("")
		refreshCategories()
		categorySelect.SetSelected(newCategory)
		saveProjectsData()
	})

	categoryScroll := container.NewVScroll(categoriesBox)
	categoryScroll.SetMinSize(fyne.NewSize(320, 170))

	deleteProjectBtn := widget.NewButtonWithIcon("Delete selected project", theme.DeleteIcon(), func() {
		if selectedProjectIndex < 0 || selectedProjectIndex >= len(projectsData.Projects) {
			openMessageWindow("Message", "Action completed or needs attention.")
			return
		}

		name := projectsData.Projects[selectedProjectIndex].Title
		openConfirmWindow("Delete project", fmt.Sprintf("Delete project %q? This cannot be undone.", name), "Delete project", func() {
			if selectedProjectIndex < 0 || selectedProjectIndex >= len(projectsData.Projects) {
				return
			}

			projectsData.Projects = append(projectsData.Projects[:selectedProjectIndex], projectsData.Projects[selectedProjectIndex+1:]...)
			selectedProjectIndex = -1
			saveProjectsData()
			refreshProjectList()
			clearForm()
		})
	})
	deleteProjectBtn.Importance = widget.DangerImportance

	removeCategoryBtn := widget.NewButtonWithIcon("Remove last category", theme.DeleteIcon(), func() {
		if len(projectsData.Categories) == 0 {
			openMessageWindow("Message", "Action completed or needs attention.")
			return
		}

		last := projectsData.Categories[len(projectsData.Categories)-1]
		openConfirmWindow("Remove category", fmt.Sprintf("Remove category %q? Existing projects keep the text, but it will disappear from the selector.", last), "Remove category", func() {
			if len(projectsData.Categories) == 0 {
				return
			}

			projectsData.Categories = projectsData.Categories[:len(projectsData.Categories)-1]
			refreshCategories()
			if len(projectsData.Categories) > 0 {
				categorySelect.SetSelected(projectsData.Categories[0])
			} else {
				categorySelect.SetSelected("")
			}
			saveProjectsData()
		})
	})
	removeCategoryBtn.Importance = widget.DangerImportance

	leftSide := panel("Projects", container.NewPadded(projectList))

	manageTab := container.NewVScroll(container.NewPadded(container.NewVBox(
		panel("Project actions", container.NewVBox(
			container.NewGridWithColumns(2, loadBtn, updateBtn),
			container.NewGridWithColumns(2, moveUpBtn, moveDownBtn),
		)),
		panel("Categories", container.NewVBox(
			widget.NewLabel("Add a new category:"),
			container.NewBorder(nil, nil, nil, addCategoryBtn, container.NewGridWrap(fyne.NewSize(280, 42), categoryEntry)),
			categoriesCount,
			categoryScroll,
		)),
		panel("Project details", container.NewVBox(
			titleEntry,
			descEntry,
			container.NewGridWithColumns(2, imageEntry, yearEntry),
			techEntry,
			categorySelect,
			featuredCheck,
		)),
		panel("Links", container.NewVBox(githubEntry, liveEntry, releaseEntry)),
		panel("Save", container.NewGridWithColumns(2, addBtn, widget.NewButtonWithIcon("Clear form", theme.ContentClearIcon(), clearForm))),
	)))

	deleteTab := container.NewVScroll(container.NewPadded(container.NewVBox(
		panel("Delete zone", container.NewVBox(
			widget.NewLabel("Select a project from the list on the left, then delete it here."),
			deleteProjectBtn,
		)),
		panel("Category delete zone", container.NewVBox(
			widget.NewLabel("This removes only the last category from the category list."),
			widget.NewLabel("Projects that already use that category will not be deleted."),
			removeCategoryBtn,
		)),
	)))

	rightTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Manage", theme.SettingsIcon(), manageTab),
		container.NewTabItemWithIcon("Delete Zone", theme.DeleteIcon(), deleteTab),
	)
	rightTabs.SetTabLocation(container.TabLocationTop)

	split := container.NewHSplit(leftSide, rightTabs)
	split.Offset = 0.38
	return split
}

func createSaveTab() fyne.CanvasObject {
	dir, _ := os.Getwd()

	stats := container.NewGridWithColumns(2,
		container.NewVBox(
			widget.NewLabelWithStyle("Changelog", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(fmt.Sprintf("%d releases", len(currentData.Entries))),
			widget.NewLabel(fmt.Sprintf("%d roadmap items", len(currentData.FutureFeatures))),
		),
		container.NewVBox(
			widget.NewLabelWithStyle("Projects", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(fmt.Sprintf("%d projects", len(projectsData.Projects))),
			widget.NewLabel(fmt.Sprintf("%d categories", len(projectsData.Categories))),
		),
	)

	files := container.NewVBox(
		widget.NewLabel("Working directory:"),
		widget.NewLabel(dir),
		widget.NewLabel("Data files:"),
		widget.NewLabel("• "+changelogFile),
		widget.NewLabel("• "+projectsFile),
	)

	actions := container.NewGridWithColumns(2,
		widget.NewButtonWithIcon("Save all", theme.DocumentSaveIcon(), saveAllData),
		widget.NewButtonWithIcon("Reload", theme.ViewRefreshIcon(), func() {
			loadChangelogData()
			loadProjectsData()
			mainWindow.SetContent(createMainMenu())
		}),
	)

	return container.NewVScroll(container.NewPadded(container.NewVBox(
		panel("Portfolio statistics", stats),
		panel("File locations", files),
		panel("System actions", actions),
		widget.NewLabelWithStyle("Tip: keep changelog.json and projects.json near the executable.", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
	)))
}
