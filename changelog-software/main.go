package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
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

// Data structures - keeping it organized like my desk (mostly)
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

type ProjectsData struct {
	Projects   []Project `json:"projects"`
	Categories []string  `json:"categories"`
	Settings   struct {
		ShowFeaturedOnly bool `json:"showFeaturedOnly"`
		ItemsPerPage     int  `json:"itemsPerPage"`
	} `json:"settings"`
}

var (
	changelogFile = "changelog.json"
	projectsFile  = "projects.json"
	currentData   Changelog
	projectsData  ProjectsData
	mainWindow    fyne.Window
)

// Track selected project
var (
	selectedProjectIndex int = -1
	selectedProject      Project
)

// Modern color palette
var (
	darkBg        = color.NRGBA{R: 18, G: 18, B: 24, A: 255}      // Deeper dark
	darkSurface   = color.NRGBA{R: 30, G: 30, B: 38, A: 255}      // Slightly lighter
	darkCard      = color.NRGBA{R: 38, G: 38, B: 46, A: 255}      // Card background
	darkBorder    = color.NRGBA{R: 55, G: 55, B: 65, A: 255}      // Border color
	accentBlue    = color.NRGBA{R: 0, G: 122, B: 255, A: 255}     // Vibrant blue
	accentGreen   = color.NRGBA{R: 80, G: 200, B: 80, A: 255}     // Fresh green
	accentRed     = color.NRGBA{R: 255, G: 69, B: 58, A: 255}     // iOS-style red
	accentOrange  = color.NRGBA{R: 255, G: 159, B: 10, A: 255}    // Warm orange
	textPrimary   = color.NRGBA{R: 245, G: 245, B: 250, A: 255}   // Almost white
	textSecondary = color.NRGBA{R: 160, G: 160, B: 180, A: 255}   // Subtle gray
)

func main() {
	myApp := app.New()
	// maybe add later ¯\_(ツ)_/¯
	mainWindow = myApp.NewWindow("AndreiixeDev")
	mainWindow.Resize(fyne.NewSize(1400, 900)) // More breathing room
	
	// Load our data
	loadChangelogData()
	loadProjectsData()
	
	// Create the main layout
	background := canvas.NewRectangle(darkBg)
	
	header := createModernHeader()
	toolbar := createModernToolbar()
	
	// Create tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("📊 Dashboard", createStatsTab()),
		container.NewTabItem("📝 Changelog", createEntriesTab()),
		container.NewTabItem("🔮 Roadmap", createFutureTab()),
		container.NewTabItem("💻 Projects", createProjectsTab()),
		container.NewTabItem("⚙️ Settings", createSaveTab()),
	)
	
	tabs.SetTabLocation(container.TabLocationTop)
	
	// Stack everything together
	content := container.NewBorder(
		container.NewVBox(
			header,
			toolbar,
			createModernSeparator(),
		),
		nil,
		nil,
		nil,
		container.NewMax(
			background,
			container.NewPadded(tabs),
		),
	)
	
	mainWindow.SetContent(content)
	mainWindow.ShowAndRun()
}

// Modern header with gradient effect (kinda, rectangles are hard ok?)
func createModernHeader() fyne.CanvasObject {
	title := canvas.NewText("AndreiixeDev Studio", textPrimary)
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}
	
	subtitle := canvas.NewText("Content Management System", textSecondary)
	subtitle.TextSize = 14
	
	accentLine := canvas.NewRectangle(accentBlue)
	accentLine.SetMinSize(fyne.NewSize(1400, 3))
	
	return container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(subtitle),
		accentLine,
	)
}

func createModernToolbar() fyne.CanvasObject {
	saveBtn := widget.NewButtonWithIcon("Save Everything", theme.DocumentSaveIcon(), func() {
		saveAllData()
	})
	saveBtn.Importance = widget.HighImportance
	
	reloadBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		loadChangelogData()
		loadProjectsData()
		dialog.ShowInformation("✨ Refreshed", "All data has been reloaded from disk!", mainWindow)
	})
	
	aboutBtn := widget.NewButtonWithIcon("About", theme.InfoIcon(), func() {
		showAbout()
	})
	
	return container.NewHBox(
		layout.NewSpacer(),
		saveBtn,
		reloadBtn,
		aboutBtn,
		layout.NewSpacer(),
	)
}

func createModernSeparator() fyne.CanvasObject {
	sep := canvas.NewRectangle(darkBorder)
	sep.SetMinSize(fyne.NewSize(1400, 1))
	return sep
}

func createModernCard(title string, content fyne.CanvasObject) fyne.CanvasObject {
	card := widget.NewCard("", title, content)
	// Make the card look less boring
	return card
}

func showAbout() {
	aboutText := `> AndreiixeDev Manager

Version 2.0.0 • "Fixes and cleanup"

What's this thing do?
• Manage your changelog without touching JSON directly
• Keep your projects organized (unlike my desktop, LOL)

© 2026 AndreiixeDev`

	dialog.ShowInformation("> About", aboutText, mainWindow)
}

// Load data with better error handling (because files can be sneaky)
func loadChangelogData() {
	data, err := os.ReadFile(changelogFile)
	if err != nil {
		// No file? No problem :]]
		currentData = Changelog{
			Stats: Stats{
				TotalUpdates:  0,
				FeaturesAdded: 0,
				BugsFixed:     0,
				Improvements:  0,
			},
			Entries:        []Entry{},
			FutureFeatures: []FutureFeature{},
		}
		saveChangelogData()
		return
	}
	
	err = json.Unmarshal(data, &currentData)
	if err != nil {
		// JSON is broken? That's a you problem ;)
		dialog.ShowError(fmt.Errorf("Your changelog.json is corrupted."), mainWindow)
		currentData = Changelog{
			Stats:          Stats{},
			Entries:        []Entry{},
			FutureFeatures: []FutureFeature{},
		}
	}
}

func loadProjectsData() {
	data, err := os.ReadFile(projectsFile)
	if err != nil {
		projectsData = ProjectsData{
			Projects:   []Project{},
			Categories: []string{},
			Settings: struct {
				ShowFeaturedOnly bool `json:"showFeaturedOnly"`
				ItemsPerPage     int  `json:"itemsPerPage"`
			}{
				ShowFeaturedOnly: false,
				ItemsPerPage:     6,
			},
		}
		saveProjectsData()
		return
	}
	
	err = json.Unmarshal(data, &projectsData)
	if err != nil {
		dialog.ShowError(fmt.Errorf("projects.json is corrupted."), mainWindow)
		projectsData = ProjectsData{
			Projects:   []Project{},
			Categories: []string{},
			Settings: struct {
				ShowFeaturedOnly bool `json:"showFeaturedOnly"`
				ItemsPerPage     int  `json:"itemsPerPage"`
			}{
				ShowFeaturedOnly: false,
				ItemsPerPage:     6,
			},
		}
	}
}

func saveChangelogData() {
	recalculateStats()
	
	data, err := json.MarshalIndent(currentData, "", "  ")
	if err != nil {
		dialog.ShowError(err, mainWindow)
		return
	}
	
	err = os.WriteFile(changelogFile, data, 0644)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}
}

func saveProjectsData() {
	data, err := json.MarshalIndent(projectsData, "", "  ")
	if err != nil {
		dialog.ShowError(err, mainWindow)
		return
	}
	
	err = os.WriteFile(projectsFile, data, 0644)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}
}

func saveAllData() {
	saveChangelogData()
	saveProjectsData()
	dialog.ShowInformation("💾 Saved!", "All your data is safely stored.", mainWindow)
}

// Recalculate stats
func recalculateStats() {
	totalUpdates := len(currentData.Entries)
	featuresAdded := 0
	bugsFixed := 0
	improvements := 0
	
	for _, entry := range currentData.Entries {
		for _, change := range entry.Changes {
			switch change.Type {
			case "added":
				featuresAdded++
			case "fixed":
				bugsFixed++
			case "improved":
				improvements++
			}
		}
	}
	
	currentData.Stats = Stats{
		TotalUpdates:  totalUpdates,
		FeaturesAdded: featuresAdded,
		BugsFixed:     bugsFixed,
		Improvements:  improvements,
	}
}

// Dashboard tab
func createStatsTab() fyne.CanvasObject {
	// Create fancy shitt stat cards
	totalCard := widget.NewCard("📦 Total Updates", "",
		container.NewCenter(
			widget.NewLabelWithStyle(fmt.Sprintf("%d", currentData.Stats.TotalUpdates), 
				fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true}),
		),
	)
	
	featuresCard := widget.NewCard("✨ Features Added", "",
		container.NewCenter(
			widget.NewLabelWithStyle(fmt.Sprintf("%d", currentData.Stats.FeaturesAdded), 
				fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true}),
		),
	)
	
	bugsCard := widget.NewCard("🐛 Bugs Fixed", "",
		container.NewCenter(
			widget.NewLabelWithStyle(fmt.Sprintf("%d", currentData.Stats.BugsFixed), 
				fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true}),
		),
	)
	
	improvementsCard := widget.NewCard("📈 Improvements", "",
		container.NewCenter(
			widget.NewLabelWithStyle(fmt.Sprintf("%d", currentData.Stats.Improvements), 
				fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true}),
		),
	)
	
	// Grid layout for stats
	statsGrid := container.NewGridWithColumns(2, totalCard, featuresCard, bugsCard, improvementsCard)
	
	// Quick actions - for the lazy (like me)
	quickActions := createModernCard("⚡ Quick Actions",
		container.NewVBox(
			widget.NewButton("📝 Add New Version", func() {
				// Switch to changelog tab and scroll? Maybe future feature
				dialog.ShowInformation("Coming Soon", "This will auto-switch to Changelog tab in the next update.", mainWindow)
			}),
			widget.NewButton("💻 Add New Project", func() {
				dialog.ShowInformation("Coming Soon", "Auto-switch to Projects tab coming.", mainWindow)
			}),
			widget.NewButton("🔄 Refresh Stats", func() {
				recalculateStats()
				// Force refresh the tab by recreating?
				dialog.ShowInformation("Stats Updated", "Statistics recalculated! Click the tab again to see changes.", mainWindow)
				//Nah, just tell them to restart
			}),
		),
	)
	
	return container.NewVScroll(
		container.NewPadded(
			container.NewVBox(
				statsGrid,
				createModernSeparator(),
				quickActions,
			),
		),
	)
}

// Changelog tab
func createEntriesTab() fyne.CanvasObject {
	// Entry list with better display
	entryList := widget.NewList(
		func() int {
			return len(currentData.Entries)
		},
		func() fyne.CanvasObject {
			// Modern list item with icon and metadata
			icon := widget.NewIcon(theme.DocumentIcon())
			titleLabel := widget.NewLabel("")
			versionLabel := widget.NewLabel("")
			dateLabel := widget.NewLabel("")
			
			return container.NewBorder(
				nil, nil,
				icon,
				dateLabel,
				container.NewHBox(
					versionLabel,
					widget.NewLabel("•"),
					titleLabel,
				),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(currentData.Entries) {
				return
			}
			
			entry := currentData.Entries[i]
			border := o.(*fyne.Container)
			
			// The structure is: border (left icon, right date, center hbox)
			if len(border.Objects) >= 2 {
				if hbox, ok := border.Objects[0].(*fyne.Container); ok {
					if len(hbox.Objects) >= 3 {
						if versionLabel, ok := hbox.Objects[0].(*widget.Label); ok {
							versionLabel.SetText(entry.Version)
						}
						if titleLabel, ok := hbox.Objects[2].(*widget.Label); ok {
							titleLabel.SetText(entry.Title)
						}
					}
				}
				if dateLabel, ok := border.Objects[1].(*widget.Label); ok {
					dateLabel.SetText(entry.Date)
				}
			}
		},
	)
	
	// Form inputs
	versionEntry := widget.NewEntry()
	versionEntry.SetPlaceHolder("v2.0.0")
	
	dateEntry := widget.NewEntry()
	dateEntry.SetPlaceHolder("YYYY-MM-DD")
	dateEntry.SetText(time.Now().Format("2006-01-02"))
	
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("What's new?")
	
	iconSelect := widget.NewSelect([]string{
		"⭐ star", "🚀 rocket", "🎨 paint", "🔗 link",
		"👤 user", "📊 chart", "🎮 game", "📱 mobile",
	}, func(s string) {})
	
	// Changes section - where the details go
	changeType := widget.NewSelect([]string{"added", "fixed", "improved", "changed", "removed"}, nil)
	changeIcon := widget.NewEntry()
	changeIcon.SetPlaceHolder("Icon name (e.g., 'bug', 'rocket')")
	changeDesc := widget.NewEntry()
	changeDesc.SetPlaceHolder("Describe what changed (and why it matters)")
	
	changesList := []Change{}
	changesLabel := widget.NewLabel("📝 0 changes")
	
	addChangeBtn := widget.NewButton("➕ Add Change", func() {
		if changeDesc.Text == "" {
			dialog.ShowInformation("Missing Info", "You need to describe the change", mainWindow)
			return
		}
		
		newChange := Change{
			Type:        changeType.Selected,
			Icon:        changeIcon.Text,
			Description: changeDesc.Text,
		}
		changesList = append(changesList, newChange)
		changesLabel.SetText(fmt.Sprintf("✅ %d changes added", len(changesList)))
		
		// Clear inputs for next change
		changeDesc.SetText("")
		changeIcon.SetText("")
	})
	
	addEntryBtn := widget.NewButtonWithIcon("✨ Save Entry", theme.ConfirmIcon(), func() {
		if versionEntry.Text == "" || titleEntry.Text == "" {
			dialog.ShowInformation("Oops!", "Version and Title are required fields. Can't have a changelog entry without them!", mainWindow)
			return
		}
		
		newEntry := Entry{
			Version: versionEntry.Text,
			Date:    dateEntry.Text,
			Title:   titleEntry.Text,
			Icon:    strings.TrimPrefix(iconSelect.Selected, "⭐ "),
			Changes: changesList,
		}
		
		// Add to top of list (newest first)
		currentData.Entries = append([]Entry{newEntry}, currentData.Entries...)
		entryList.Refresh()
		saveChangelogData()
		
		// Reset form
		versionEntry.SetText("")
		titleEntry.SetText("")
		changesList = []Change{}
		changesLabel.SetText("📝 0 changes")
		
		dialog.ShowInformation("Success! 🎉", "Your changelog entry has been saved.", mainWindow)
	})
	
	deleteBtn := widget.NewButtonWithIcon("🗑️ Delete Last", theme.DeleteIcon(), func() {
		if len(currentData.Entries) == 0 {
			dialog.ShowInformation("Nothing to delete", "The changelog is empty.", mainWindow)
			return
		}
		
		dialog.ShowConfirm("Confirm Delete", "Are you sure you want to delete the most recent entry? This cannot be undone!", func(confirm bool) {
			if confirm {
				currentData.Entries = currentData.Entries[1:]
				entryList.Refresh()
				saveChangelogData()
				dialog.ShowInformation("Deleted", "Entry removed.", mainWindow)
			}
		}, mainWindow)
	})
	
	// Build the form
	form := container.NewVScroll(
		container.NewPadded(
			container.NewVBox(
				createModernCard("📝 New Release Entry",
					container.NewVBox(
						container.NewGridWithColumns(2,
							container.NewVBox(
								widget.NewLabel("Version:"),
								versionEntry,
							),
							container.NewVBox(
								widget.NewLabel("Release Date:"),
								dateEntry,
							),
						),
						widget.NewLabel("Title:"),
						titleEntry,
						widget.NewLabel("Icon:"),
						iconSelect,
					),
				),
				createModernCard("🔧 Changes & Improvements",
					container.NewVBox(
						container.NewGridWithColumns(3,
							container.NewVBox(
								widget.NewLabel("Type:"),
								changeType,
							),
							container.NewVBox(
								widget.NewLabel("Icon:"),
								changeIcon,
							),
							container.NewVBox(
								widget.NewLabel("Description:"),
								changeDesc,
							),
						),
						addChangeBtn,
						changesLabel,
					),
				),
				createModernCard("⚡ Actions",
					container.NewHBox(
						addEntryBtn,
						deleteBtn,
					),
				),
			),
		),
	)
	
	// Split view: list on left, form on right
	split := container.NewHSplit(
		container.NewBorder(
			createModernCard("📋 Version History", nil),
			nil, nil, nil,
			container.NewPadded(entryList),
		),
		form,
	)
	split.Offset = 0.4
	
	return split
}

// Future features tab
func createFutureTab() fyne.CanvasObject {
	featureList := widget.NewList(
		func() int {
			return len(currentData.FutureFeatures)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.ConfirmIcon()),
				widget.NewLabel(""),
				layout.NewSpacer(),
				widget.NewLabel(""),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(currentData.FutureFeatures) {
				return
			}
			
			feature := currentData.FutureFeatures[i]
			icon := feature.Icon
			if icon == "" {
				icon = "🔮"
			}
			
			box := o.(*fyne.Container)
			if len(box.Objects) > 1 {
				box.Objects[1].(*widget.Label).SetText(
					fmt.Sprintf("%s %s", icon, feature.Name),
				)
			}
		},
	)
	
	iconSelect := widget.NewSelect([]string{
		"💬 chat", "📈 analytics", "🔍 search", "🔄 sync",
		"❤️ love", "⭐ star", "🔔 notify", "✉️ email",
		"🎮 game", "📱 mobile", "💻 desktop", "🌐 web",
	}, nil)
	
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("What awesome feature is coming?")
	
	addBtn := widget.NewButtonWithIcon("➕ Add to Roadmap", theme.ContentAddIcon(), func() {
		if nameEntry.Text == "" {
			dialog.ShowInformation("Need a name", "Even future features need names.", mainWindow)
			return
		}
		
		newFeature := FutureFeature{
			Icon: strings.TrimPrefix(iconSelect.Selected, "💬 "),
			Name: nameEntry.Text,
		}
		
		currentData.FutureFeatures = append(currentData.FutureFeatures, newFeature)
		featureList.Refresh()
		saveChangelogData()
		
		nameEntry.SetText("")
		dialog.ShowInformation("Added to Roadmap! 🚀", "Future feature saved.", mainWindow)
	})
	
	deleteBtn := widget.NewButtonWithIcon("🗑️ Remove Last", theme.DeleteIcon(), func() {
		if len(currentData.FutureFeatures) == 0 {
			dialog.ShowInformation("Empty Roadmap", "No features to remove.", mainWindow)
			return
		}
		
		dialog.ShowConfirm("Remove Feature", "Are you sure? This feature might never get built now...", func(confirm bool) {
			if confirm {
				currentData.FutureFeatures = currentData.FutureFeatures[:len(currentData.FutureFeatures)-1]
				featureList.Refresh()
				saveChangelogData()
			}
		}, mainWindow)
	})
	
	form := container.NewPadded(
		createModernCard("🔮 Add Future Feature",
			container.NewVBox(
				iconSelect,
				nameEntry,
				container.NewHBox(addBtn, deleteBtn),
			),
		),
	)
	
	split := container.NewHSplit(
		container.NewBorder(
			createModernCard("🔮 Planned Features", nil),
			nil, nil, nil,
			container.NewPadded(featureList),
		),
		form,
	)
	split.Offset = 0.4
	
	return split
}

// Helper to populate project form (because DRY is nice)
func populateProjectForm(
	titleEntry, descEntry, imageEntry, yearEntry, techEntry, githubEntry, liveEntry, releaseEntry *widget.Entry,
	categorySelect *widget.Select,
	featuredCheck *widget.Check,
	project Project,
) {
	titleEntry.SetText(project.Title)
	descEntry.SetText(project.Description)
	imageEntry.SetText(project.Image)
	yearEntry.SetText(project.Year)
	techEntry.SetText(strings.Join(project.Technologies, ", "))
	githubEntry.SetText(project.Links.Github)
	liveEntry.SetText(project.Links.Live)
	releaseEntry.SetText(project.Links.Release)
	categorySelect.SetSelected(project.Category)
	featuredCheck.SetChecked(project.Featured)
}

// Projects tab
func createProjectsTab() fyne.CanvasObject {
	// Project list with better visual feedback
	projectList := widget.NewList(
		func() int {
			return len(projectsData.Projects)
		},
		func() fyne.CanvasObject {
			itemBg := canvas.NewRectangle(darkSurface)
			itemBg.CornerRadius = 8 // Rounded corners because we're fancy
			
			iconWidget := widget.NewIcon(theme.ComputerIcon())
			titleLabel := widget.NewLabel("")
			categoryLabel := widget.NewLabel("")
			
			// Modern layout
			contentBox := container.NewHBox(
				iconWidget,
				container.NewVBox(
					titleLabel,
					categoryLabel,
				),
				layout.NewSpacer(),
			)
			
			return container.NewMax(
				itemBg,
				container.NewPadded(contentBox),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(projectsData.Projects) {
				return
			}
			
			project := projectsData.Projects[i]
			featuredStar := ""
			if project.Featured {
				featuredStar = " ⭐"
			}
			
			mainContainer := o.(*fyne.Container)
			if len(mainContainer.Objects) > 0 {
				if bgRect, ok := mainContainer.Objects[0].(*canvas.Rectangle); ok {
					if i == selectedProjectIndex {
						bgRect.FillColor = accentBlue
					} else {
						bgRect.FillColor = darkSurface
					}
					bgRect.Refresh()
				}
			}
			
			if len(mainContainer.Objects) > 1 {
				if paddedContainer, ok := mainContainer.Objects[1].(*fyne.Container); ok {
					if len(paddedContainer.Objects) > 0 {
						if contentBox, ok := paddedContainer.Objects[0].(*fyne.Container); ok {
							if len(contentBox.Objects) > 1 {
								if vbox, ok := contentBox.Objects[1].(*fyne.Container); ok {
									if len(vbox.Objects) > 0 {
										if titleLabel, ok := vbox.Objects[0].(*widget.Label); ok {
											titleLabel.SetText(fmt.Sprintf("%s%s", project.Title, featuredStar))
										}
										if len(vbox.Objects) > 1 {
											if categoryLabel, ok := vbox.Objects[1].(*widget.Label); ok {
												categoryLabel.SetText(project.Category)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		},
	)
	
	projectList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(projectsData.Projects) {
			selectedProjectIndex = id
			selectedProject = projectsData.Projects[id]
			projectList.Refresh()
		}
	}
	
	// Form inputs
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Project name")
	
	descEntry := widget.NewEntry()
	descEntry.SetPlaceHolder("What does this project do?")
	descEntry.MultiLine = true
	
	imageEntry := widget.NewEntry()
	imageEntry.SetPlaceHolder("images/project.jpg")
	
	yearEntry := widget.NewEntry()
	yearEntry.SetPlaceHolder("2024 • When did you build this?")
	
	techEntry := widget.NewEntry()
	techEntry.SetPlaceHolder("Go, React, Python, etc. • comma separated")
	
	// Category selector with live updates
	categorySelect := widget.NewSelect(projectsData.Categories, nil)
	if len(projectsData.Categories) > 0 {
		categorySelect.SetSelected(projectsData.Categories[0])
	}
	
	featuredCheck := widget.NewCheck("⭐ Featured Project", nil)
	
	// Links section
	githubEntry := widget.NewEntry()
	githubEntry.SetPlaceHolder("https://github.com/username/project")
	
	liveEntry := widget.NewEntry()
	liveEntry.SetPlaceHolder("https://demo.project.com")
	
	releaseEntry := widget.NewEntry()
	releaseEntry.SetPlaceHolder("https://releases.project.com")
	
	// Buttons with better UX
	moveUpBtn := widget.NewButtonWithIcon("↑ Move Up", theme.MoveUpIcon(), func() {
		if selectedProjectIndex <= 0 {
			dialog.ShowInformation("Can't Move Up", "This project is already at the top of the list!", mainWindow)
			return
		}
		
		projectsData.Projects[selectedProjectIndex], projectsData.Projects[selectedProjectIndex-1] =
			projectsData.Projects[selectedProjectIndex-1], projectsData.Projects[selectedProjectIndex]
		
		selectedProjectIndex--
		projectList.Refresh()
		saveProjectsData()
		dialog.ShowInformation("Reordered", "Project moved up.", mainWindow)
	})
	
	moveDownBtn := widget.NewButtonWithIcon("↓ Move Down", theme.MoveDownIcon(), func() {
		if selectedProjectIndex < 0 || selectedProjectIndex >= len(projectsData.Projects)-1 {
			dialog.ShowInformation("Can't Move Down", "This project is already at the bottom!", mainWindow)
			return
		}
		
		projectsData.Projects[selectedProjectIndex], projectsData.Projects[selectedProjectIndex+1] =
			projectsData.Projects[selectedProjectIndex+1], projectsData.Projects[selectedProjectIndex]
		
		selectedProjectIndex++
		projectList.Refresh()
		saveProjectsData()
	})
	
	editProjectBtn := widget.NewButtonWithIcon("✏️ Load to Edit", theme.DocumentCreateIcon(), func() {
		if selectedProjectIndex < 0 {
			dialog.ShowInformation("Select First", "Click on a project in the list to edit it!", mainWindow)
			return
		}
		
		populateProjectForm(
			titleEntry, descEntry, imageEntry, yearEntry, techEntry,
			githubEntry, liveEntry, releaseEntry,
			categorySelect, featuredCheck,
			selectedProject,
		)
		
		dialog.ShowInformation("Ready to Edit", "Project loaded! Make your changes and click Update.", mainWindow)
	})
	
	updateProjectBtn := widget.NewButtonWithIcon("🔄 Update Project", theme.ConfirmIcon(), func() {
		if selectedProjectIndex < 0 {
			dialog.ShowInformation("No Selection", "Select a project first (click on it in the list)", mainWindow)
			return
		}
		
		if titleEntry.Text == "" || descEntry.Text == "" {
			dialog.ShowInformation("Missing Info", "Title and Description are required.", mainWindow)
			return
		}
		
		// Parse technologies
		technologies := []string{}
		if techEntry.Text != "" {
			techs := strings.Split(techEntry.Text, ",")
			for _, tech := range techs {
				if trimmed := strings.TrimSpace(tech); trimmed != "" {
					technologies = append(technologies, trimmed)
				}
			}
		}
		
		updatedProject := Project{
			ID:           selectedProject.ID,
			Title:        titleEntry.Text,
			Description:  descEntry.Text,
			Image:        imageEntry.Text,
			Year:         yearEntry.Text,
			Technologies: technologies,
			Category:     categorySelect.Selected,
			Featured:     featuredCheck.Checked,
			Links: ProjectLinks{
				Github:  githubEntry.Text,
				Live:    liveEntry.Text,
				Release: releaseEntry.Text,
			},
		}
		
		projectsData.Projects[selectedProjectIndex] = updatedProject
		selectedProject = updatedProject
		
		projectList.Refresh()
		saveProjectsData()
		
		// Clear form for next project
		titleEntry.SetText("")
		descEntry.SetText("")
		imageEntry.SetText("")
		yearEntry.SetText("")
		techEntry.SetText("")
		githubEntry.SetText("")
		liveEntry.SetText("")
		releaseEntry.SetText("")
		featuredCheck.SetChecked(false)
		
		dialog.ShowInformation("Updated! 🎉", "Project changes saved.", mainWindow)
	})
	
	addProjectBtn := widget.NewButtonWithIcon("➕ Add New Project", theme.ContentAddIcon(), func() {
		if titleEntry.Text == "" || descEntry.Text == "" {
			dialog.ShowInformation("Missing Info", "Title and Description are required for new projects!", mainWindow)
			return
		}
		
		if len(projectsData.Categories) == 0 {
			dialog.ShowInformation("Need Categories", "Create at least one category before adding projects!", mainWindow)
			return
		}
		
		if categorySelect.Selected == "" {
			dialog.ShowInformation("Select Category", "Pick a category for this project!", mainWindow)
			return
		}
		
		// Generate new ID
		newID := 1
		if len(projectsData.Projects) > 0 {
			newID = projectsData.Projects[len(projectsData.Projects)-1].ID + 1
		}
		
		technologies := []string{}
		if techEntry.Text != "" {
			techs := strings.Split(techEntry.Text, ",")
			for _, tech := range techs {
				if trimmed := strings.TrimSpace(tech); trimmed != "" {
					technologies = append(technologies, trimmed)
				}
			}
		}
		
		newProject := Project{
			ID:           newID,
			Title:        titleEntry.Text,
			Description:  descEntry.Text,
			Image:        imageEntry.Text,
			Year:         yearEntry.Text,
			Technologies: technologies,
			Category:     categorySelect.Selected,
			Featured:     featuredCheck.Checked,
			Links: ProjectLinks{
				Github:  githubEntry.Text,
				Live:    liveEntry.Text,
				Release: releaseEntry.Text,
			},
		}
		
		projectsData.Projects = append(projectsData.Projects, newProject)
		projectList.Refresh()
		saveProjectsData()
		
		// Clear form
		titleEntry.SetText("")
		descEntry.SetText("")
		imageEntry.SetText("")
		yearEntry.SetText("")
		techEntry.SetText("")
		githubEntry.SetText("")
		liveEntry.SetText("")
		releaseEntry.SetText("")
		featuredCheck.SetChecked(false)
		
		dialog.ShowInformation("Project Added! 🚀", fmt.Sprintf("%s has been added.", newProject.Title), mainWindow)
	})
	
	deleteProjectBtn := widget.NewButtonWithIcon("🗑️ Delete Selected", theme.DeleteIcon(), func() {
		if selectedProjectIndex < 0 {
			dialog.ShowInformation("Nothing Selected", "Click on a project to select it first!", mainWindow)
			return
		}
		
		dialog.ShowConfirm("Delete Project", fmt.Sprintf("Are you sure you want to delete '%s'? This can't be undone!", selectedProject.Title), func(confirm bool) {
			if confirm {
				projectsData.Projects = append(
					projectsData.Projects[:selectedProjectIndex],
					projectsData.Projects[selectedProjectIndex+1:]...,
				)
				selectedProjectIndex = -1
				projectList.Refresh()
				saveProjectsData()
				dialog.ShowInformation("Deleted", "Project removed.", mainWindow)
			}
		}, mainWindow)
	})
	
	// Category management section
	categoryEntry := widget.NewEntry()
	categoryEntry.SetPlaceHolder("e.g., web, mobile, ai, game")
	
	categoriesCount := widget.NewLabel(fmt.Sprintf("📁 %d categories", len(projectsData.Categories)))
	
	categoriesList := container.NewVBox()
	refreshCategoriesList := func() {
		categoriesList.Objects = nil
		for i, cat := range projectsData.Categories {
			catLabel := widget.NewLabel(fmt.Sprintf("%d. 🏷️ %s", i+1, cat))
			categoriesList.AddObject(catLabel)
		}
		categoriesCount.SetText(fmt.Sprintf("📁 %d categories", len(projectsData.Categories)))
		categoriesList.Refresh()
	}
	
	refreshCategoriesList()
	
	addCategoryBtn := widget.NewButton("➕ Add Category", func() {
		if categoryEntry.Text == "" {
			dialog.ShowInformation("Need a Name", "What should we call this category?", mainWindow)
			return
		}
		
		// Check for duplicates
		for _, cat := range projectsData.Categories {
			if strings.EqualFold(cat, categoryEntry.Text) {
				dialog.ShowInformation("Already Exists", fmt.Sprintf("'%s' category already exists!", cat), mainWindow)
				return
			}
		}
		
		projectsData.Categories = append(projectsData.Categories, categoryEntry.Text)
		
		// Update selector
		categorySelect.Options = projectsData.Categories
		if len(projectsData.Categories) > 0 {
			categorySelect.SetSelected(projectsData.Categories[0])
		}
		categorySelect.Refresh()
		
		refreshCategoriesList()
		
		categoryEntry.SetText("")
		saveProjectsData()
		dialog.ShowInformation("Category Added! 🎯", fmt.Sprintf("'%s' is now available for projects.", categoryEntry.Text), mainWindow)
	})
	
	deleteCategoryBtn := widget.NewButtonWithIcon("🗑️ Remove Last", theme.DeleteIcon(), func() {
		if len(projectsData.Categories) == 0 {
			return
		}
		
		lastCategory := projectsData.Categories[len(projectsData.Categories)-1]
		dialog.ShowConfirm("Delete Category", fmt.Sprintf("Delete '%s' category?\nProjects using this category will lose their category!", lastCategory), func(confirm bool) {
			if confirm {
				projectsData.Categories = projectsData.Categories[:len(projectsData.Categories)-1]
				
				categorySelect.Options = projectsData.Categories
				if len(projectsData.Categories) > 0 {
					categorySelect.SetSelected(projectsData.Categories[0])
				} else {
					categorySelect.SetSelected("")
				}
				categorySelect.Refresh()
				
				refreshCategoriesList()
				saveProjectsData()
				dialog.ShowInformation("Category Removed", "Make sure to reassign any affected projects!", mainWindow)
			}
		}, mainWindow)
	})
	
	categoriesScroll := container.NewVScroll(categoriesList)
	categoriesScroll.SetMinSize(fyne.NewSize(300, 150))
	
	// Assemble the form
	moveToolbar := container.NewHBox(
		moveUpBtn,
		moveDownBtn,
		editProjectBtn,
		updateProjectBtn,
	)
	
	form := container.NewVScroll(
		container.NewPadded(
			container.NewVBox(
				createModernCard("🎯 Project Actions",
					container.NewVBox(
						widget.NewLabel("Select a project from the list, then:"),
						moveToolbar,
					),
				),
				
				createModernCard("🏷️ Manage Categories",
					container.NewVBox(
						widget.NewLabel("Add categories to organize your work:"),
						container.NewHBox(
							categoryEntry,
							addCategoryBtn,
						),
						categoriesCount,
						categoriesScroll,
						deleteCategoryBtn,
					),
				),
				
				createModernSeparator(),
				
				createModernCard("💻 Project Details",
					container.NewVBox(
						titleEntry,
						descEntry,
						container.NewGridWithColumns(2,
							imageEntry,
							yearEntry,
						),
						techEntry,
						container.NewVBox(
							widget.NewLabel("Category:"),
							categorySelect,
							featuredCheck,
						),
					),
				),
				
				createModernCard("🔗 Links (optional)",
					container.NewVBox(
						githubEntry,
						liveEntry,
						releaseEntry,
					),
				),
				
				createModernCard("⚡ Actions",
					container.NewHBox(
						addProjectBtn,
						deleteProjectBtn,
					),
				),
			),
		),
	)
	
	split := container.NewHSplit(
		container.NewBorder(
			createModernCard("📋 Your Projects (click to select)", nil),
			nil, nil, nil,
			container.NewPadded(projectList),
		),
		form,
	)
	split.Offset = 0.4
	
	return split
}

// Settings/Save tab
func createSaveTab() fyne.CanvasObject {
	dir, _ := os.Getwd()
	
	stats := createModernCard("📊 Portfolio Statistics",
		container.NewGridWithColumns(2,
			container.NewVBox(
				widget.NewLabelWithStyle("Changelog:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel(fmt.Sprintf("📝 %d version entries", len(currentData.Entries))),
				widget.NewLabel(fmt.Sprintf("🔮 %d planned features", len(currentData.FutureFeatures))),
			),
			container.NewVBox(
				widget.NewLabelWithStyle("Projects:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel(fmt.Sprintf("💻 %d projects", len(projectsData.Projects))),
				widget.NewLabel(fmt.Sprintf("🏷️ %d categories", len(projectsData.Categories))),
			),
		),
	)
	
	files := createModernCard("📁 File Locations",
		container.NewVBox(
			widget.NewLabel(fmt.Sprintf("📂 Working Directory: %s", dir)),
			widget.NewLabel(fmt.Sprintf("📄 %s • Changelog data", changelogFile)),
			widget.NewLabel(fmt.Sprintf("📄 %s • Projects data", projectsFile)),
		),
	)
	
	actions := createModernCard("⚙️ System Actions",
		container.NewVBox(
			widget.NewButtonWithIcon("💾 Save Everything", theme.DocumentSaveIcon(), func() {
				saveAllData()
			}),
			widget.NewButtonWithIcon("🔄 Reload from Disk", theme.ViewRefreshIcon(), func() {
				loadChangelogData()
				loadProjectsData()
				dialog.ShowInformation("Reloaded! 🔄", "All data has been refreshed from disk. Changes made elsewhere are now visible.", mainWindow)
			}),
		),
	)
	
	funFact := widget.NewLabelWithStyle("💡 Did you know? Nothing..", 
		fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	
	return container.NewVScroll(
		container.NewPadded(
			container.NewVBox(
				stats,
				files,
				actions,
				createModernSeparator(),
				funFact,
			),
		),
	)
}