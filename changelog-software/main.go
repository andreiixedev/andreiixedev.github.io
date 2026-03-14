package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
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

var (
	selectedProjectIndex int = -1
	selectedProject      Project
)

var (
	darkBg        = color.NRGBA{R: 30, G: 30, B: 30, A: 255}
	darkSurface   = color.NRGBA{R: 45, G: 45, B: 45, A: 255}
	darkCard      = color.NRGBA{R: 55, G: 55, B: 55, A: 255}
	darkBorder    = color.NRGBA{R: 80, G: 80, B: 80, A: 255}
	accentBlue    = color.NRGBA{R: 0, G: 120, B: 212, A: 255}
	accentGreen   = color.NRGBA{R: 30, G: 150, B: 30, A: 255}
	accentRed     = color.NRGBA{R: 200, G: 50, B: 50, A: 255}
	accentOrange  = color.NRGBA{R: 240, G: 140, B: 0, A: 255}
	textPrimary   = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	textSecondary = color.NRGBA{R: 180, G: 180, B: 180, A: 255}
)

func main() {
	myApp := app.New()
	mainWindow = myApp.NewWindow("AndreiixeDev - Website JSON Editor")
	mainWindow.Resize(fyne.NewSize(1300, 800))

	loadChangelogData()
	loadProjectsData()

	background := canvas.NewRectangle(darkBg)
	background.SetMinSize(fyne.NewSize(1300, 800))

	header := createDarkHeader()

	toolbar := createDarkToolbar()

	tabs := container.NewAppTabs(
		container.NewTabItem("📊 DASHBOARD", createStatsTab()),
		container.NewTabItem("📝 CHANGELOG", createEntriesTab()),
		container.NewTabItem("🔮 FUTURE", createFutureTab()),
		container.NewTabItem("💻 PROJECTS", createProjectsTab()),
		container.NewTabItem("⚙️ SETTINGS", createSaveTab()),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	content := container.NewBorder(
		container.NewVBox(
			header,
			toolbar,
			createDarkSeparator(),
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

func createDarkHeader() fyne.CanvasObject {
	title := canvas.NewText("AndreiixeDev - Website JSON Manager", textPrimary)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Manage your changelog and projects with ease", textSecondary)
	subtitle.TextSize = 14

	line := canvas.NewRectangle(accentBlue)
	line.SetMinSize(fyne.NewSize(1300, 3))

	return container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(subtitle),
		line,
	)
}

func createDarkToolbar() fyne.CanvasObject {
	saveBtn := widget.NewButtonWithIcon("Save All", theme.DocumentSaveIcon(), func() {
		saveAllData()
	})

	reloadBtn := widget.NewButtonWithIcon("Reload", theme.ViewRefreshIcon(), func() {
		loadChangelogData()
		loadProjectsData()
		dialog.ShowInformation("Success", "All data reloaded!", mainWindow)
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

func createDarkSeparator() fyne.CanvasObject {
	sep := canvas.NewRectangle(darkBorder)
	sep.SetMinSize(fyne.NewSize(1300, 1))
	return sep
}

func createDarkCard(title string, content fyne.CanvasObject) fyne.CanvasObject {
	if content == nil {
		return widget.NewCard("", title, nil)
	}
	return widget.NewCard("", title, content)
}

func showAbout() {
	aboutText := `AndreiixeDev Website JSON Editor

Version 1.0.10

Purpose:
This tool helps you manage the JSON data files for your portfolio website:
• changelog.json - Track version history and updates
• projects.json - Manage your projects and categories

Created by AndreiixeDev
© 2026 All rights reserved`

	dialog.ShowInformation("ℹ️ About AndreiixeDev JSON Editor", aboutText, mainWindow)
}

// Load functions
func loadChangelogData() {
	file, err := ioutil.ReadFile(changelogFile)
	if err != nil {
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

	err = json.Unmarshal(file, &currentData)
	if err != nil {
		log.Fatal("Error parsing changelog JSON:", err)
	}
}

func loadProjectsData() {
	file, err := ioutil.ReadFile(projectsFile)
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

	err = json.Unmarshal(file, &projectsData)
	if err != nil {
		log.Fatal("Error parsing projects JSON:", err)
	}
}

func saveChangelogData() {
	recalculateStats()

	data, err := json.MarshalIndent(currentData, "", "  ")
	if err != nil {
		dialog.ShowError(err, mainWindow)
		return
	}

	err = ioutil.WriteFile(changelogFile, data, 0644)
	if err != nil {
		dialog.ShowError(err, mainWindow)
		return
	}
}

func saveProjectsData() {
	data, err := json.MarshalIndent(projectsData, "", "  ")
	if err != nil {
		dialog.ShowError(err, mainWindow)
		return
	}

	err = ioutil.WriteFile(projectsFile, data, 0644)
	if err != nil {
		dialog.ShowError(err, mainWindow)
		return
	}
}

func saveAllData() {
	saveChangelogData()
	saveProjectsData()
	dialog.ShowInformation("Success", "✅ All data saved successfully!", mainWindow)
}

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

func createStatsTab() fyne.CanvasObject {
	bg := canvas.NewRectangle(darkBg)

	totalCard := widget.NewCard("📦 TOTAL UPDATES", "",
		container.NewCenter(
			widget.NewLabelWithStyle(fmt.Sprintf("%d", currentData.Stats.TotalUpdates), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		),
	)

	featuresCard := widget.NewCard("✨ FEATURES ADDED", "",
		container.NewCenter(
			widget.NewLabelWithStyle(fmt.Sprintf("%d", currentData.Stats.FeaturesAdded), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		),
	)

	bugsCard := widget.NewCard("🐛 BUGS FIXED", "",
		container.NewCenter(
			widget.NewLabelWithStyle(fmt.Sprintf("%d", currentData.Stats.BugsFixed), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		),
	)

	improvementsCard := widget.NewCard("📈 IMPROVEMENTS", "",
		container.NewCenter(
			widget.NewLabelWithStyle(fmt.Sprintf("%d", currentData.Stats.Improvements), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		),
	)

	statsGrid := container.NewGridWithColumns(2,
		totalCard,
		featuresCard,
		bugsCard,
		improvementsCard,
	)

	quickActions := createDarkCard("⚡ QUICK ACTIONS",
		container.NewVBox(
			widget.NewButton("📝 Add New Changelog Entry", func() {
			}),
			widget.NewButton("💻 Add New Project", func() {
			}),
			widget.NewButton("🔄 Refresh Statistics", func() {
				recalculateStats()
				dialog.ShowInformation("Success", "Statistics updated!\nRestart tab to see changes", mainWindow)
			}),
		),
	)

	return container.NewMax(
		bg,
		container.NewVScroll(
			container.NewPadded(
				container.NewVBox(
					statsGrid,
					createDarkSeparator(),
					quickActions,
				),
			),
		),
	)
}

func createEntriesTab() fyne.CanvasObject {
	bg := canvas.NewRectangle(darkBg)

	entryList := widget.NewList(
		func() int {
			return len(currentData.Entries)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel(""),
				layout.NewSpacer(),
				widget.NewLabel(""),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			entry := currentData.Entries[i]
			box := o.(*fyne.Container)
			if len(box.Objects) > 1 {
				box.Objects[1].(*widget.Label).SetText(
					fmt.Sprintf("%s - %s", entry.Version, entry.Title),
				)
			}
			if len(box.Objects) > 3 {
				box.Objects[3].(*widget.Label).SetText(entry.Date)
			}
		},
	)

	versionEntry := widget.NewEntry()
	versionEntry.SetPlaceHolder("e.g., v2.0.0")

	dateEntry := widget.NewEntry()
	dateEntry.SetPlaceHolder("2026-03-14")
	dateEntry.SetText(time.Now().Format("2006-01-02"))

	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("e.g., 3D Model & Animations")

	iconEntry := widget.NewSelect([]string{
		"⭐ star", "🚀 rocket", "🎨 paint-brush", "🔗 link",
		"👤 user", "📊 diagram", "🖌️ roller", "🌿 code-branch",
	}, func(s string) {})

	// Changes section
	changeType := widget.NewSelect([]string{"added", "changed", "fixed", "improved", "removed"}, nil)
	changeIcon := widget.NewEntry()
	changeIcon.SetPlaceHolder("e.g., cube, bug, leaf")
	changeDesc := widget.NewEntry()
	changeDesc.SetPlaceHolder("Describe the change")

	changesList := []Change{}
	changesLabel := widget.NewLabel("📋 0 changes added")

	addChangeBtn := widget.NewButton("➕ Add Change", func() {
		if changeDesc.Text == "" {
			dialog.ShowInformation("Error", "Please enter a description", mainWindow)
			return
		}

		newChange := Change{
			Type:        changeType.Selected,
			Icon:        changeIcon.Text,
			Description: changeDesc.Text,
		}
		changesList = append(changesList, newChange)
		changesLabel.SetText(fmt.Sprintf("📋 %d changes added", len(changesList)))
		changeDesc.SetText("")
		changeIcon.SetText("")
	})

	addEntryBtn := widget.NewButtonWithIcon("💾 Save Entry", theme.ConfirmIcon(), func() {
		if versionEntry.Text == "" || titleEntry.Text == "" {
			dialog.ShowInformation("Error", "Version and Title are required", mainWindow)
			return
		}

		newEntry := Entry{
			Version: versionEntry.Text,
			Date:    dateEntry.Text,
			Title:   titleEntry.Text,
			Icon:    strings.TrimPrefix(iconEntry.Selected, "⭐ "),
			Changes: changesList,
		}

		currentData.Entries = append([]Entry{newEntry}, currentData.Entries...)
		entryList.Refresh()
		saveChangelogData()

		versionEntry.SetText("")
		titleEntry.SetText("")
		changesList = []Change{}
		changesLabel.SetText("📋 0 changes added")

		dialog.ShowInformation("Success", "✅ Entry added successfully!", mainWindow)
	})

	deleteBtn := widget.NewButtonWithIcon("🗑️ Delete Last", theme.DeleteIcon(), func() {
		if len(currentData.Entries) == 0 {
			return
		}
		dialog.ShowConfirm("Confirm Delete", "Are you sure you want to delete the last entry?", func(confirm bool) {
			if confirm {
				currentData.Entries = currentData.Entries[1:]
				entryList.Refresh()
				saveChangelogData()
			}
		}, mainWindow)
	})

	form := container.NewMax(bg,
		container.NewVScroll(
			container.NewPadded(
				container.NewVBox(
					createDarkCard("📝 ADD NEW ENTRY",
						container.NewVBox(
							container.NewGridWithColumns(2,
								container.NewVBox(
									widget.NewLabel("Version:"),
									versionEntry,
								),
								container.NewVBox(
									widget.NewLabel("Date:"),
									dateEntry,
								),
							),
							widget.NewLabel("Title:"),
							titleEntry,
							widget.NewLabel("Icon:"),
							iconEntry,
						),
					),
					createDarkCard("🔧 CHANGES",
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
					createDarkCard("⚡ ACTIONS",
						container.NewHBox(
							addEntryBtn,
							deleteBtn,
						),
					),
				),
			),
		),
	)

	split := container.NewHSplit(
		container.NewBorder(
			createDarkCard("📋 ENTRIES LIST", nil),
			nil, nil, nil,
			container.NewMax(bg, container.NewPadded(entryList)),
		),
		form,
	)
	split.Offset = 0.35

	return split
}

func createFutureTab() fyne.CanvasObject {
	bg := canvas.NewRectangle(darkBg)

	featureList := widget.NewList(
		func() int {
			return len(currentData.FutureFeatures)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.ConfirmIcon()),
				widget.NewLabel(""),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
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

	iconEntry := widget.NewSelect([]string{
		"💬 comment", "📈 chart-line", "🔍 search", "🔄 share",
		"❤️ heart", "⭐ star", "🔔 bell", "✉️ envelope",
	}, nil)

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("e.g., Dark mode")

	addBtn := widget.NewButtonWithIcon("➕ Add", theme.ContentAddIcon(), func() {
		if nameEntry.Text == "" {
			dialog.ShowInformation("Error", "Please enter a feature name", mainWindow)
			return
		}

		newFeature := FutureFeature{
			Icon: strings.TrimPrefix(iconEntry.Selected, "💬 "),
			Name: nameEntry.Text,
		}

		currentData.FutureFeatures = append(currentData.FutureFeatures, newFeature)
		featureList.Refresh()
		saveChangelogData()

		nameEntry.SetText("")
		dialog.ShowInformation("Success", "✅ Feature added!", mainWindow)
	})

	deleteBtn := widget.NewButtonWithIcon("🗑️ Delete Last", theme.DeleteIcon(), func() {
		if len(currentData.FutureFeatures) == 0 {
			return
		}
		dialog.ShowConfirm("Confirm Delete", "Are you sure?", func(confirm bool) {
			if confirm {
				currentData.FutureFeatures = currentData.FutureFeatures[:len(currentData.FutureFeatures)-1]
				featureList.Refresh()
				saveChangelogData()
			}
		}, mainWindow)
	})

	form := container.NewMax(bg,
		container.NewPadded(
			createDarkCard("🔮 ADD FUTURE FEATURE",
				container.NewVBox(
					iconEntry,
					nameEntry,
					container.NewHBox(addBtn, deleteBtn),
				),
			),
		),
	)

	split := container.NewHSplit(
		container.NewBorder(
			createDarkCard("🔮 PLANNED FEATURES", nil),
			nil, nil, nil,
			container.NewMax(bg, container.NewPadded(featureList)),
		),
		form,
	)
	split.Offset = 0.4

	return split
}

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

func createProjectsTab() fyne.CanvasObject {
	bg := canvas.NewRectangle(darkBg)

	projectList := widget.NewList(
		func() int {
			return len(projectsData.Projects)
		},
		func() fyne.CanvasObject {
			itemBg := canvas.NewRectangle(darkSurface)
			itemBg.CornerRadius = 4
			itemBg.SetMinSize(fyne.NewSize(200, 50))

			iconWidget := widget.NewIcon(theme.ComputerIcon())
			titleLabel := widget.NewLabel("Loading...")
			categoryLabel := widget.NewLabel("")

			contentBox := container.NewHBox(
				iconWidget,
				titleLabel,
				layout.NewSpacer(),
				categoryLabel,
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
								if label, ok := contentBox.Objects[1].(*widget.Label); ok {
									label.SetText(fmt.Sprintf("%d. %s%s", i+1, project.Title, featuredStar))
								}
							}
							if len(contentBox.Objects) > 3 {
								if label, ok := contentBox.Objects[3].(*widget.Label); ok {
									label.SetText(project.Category)
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

	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("e.g., Portfolio Website")

	descEntry := widget.NewEntry()
	descEntry.SetPlaceHolder("Project description...")
	descEntry.MultiLine = true

	imageEntry := widget.NewEntry()
	imageEntry.SetPlaceHolder("images/project.jpg")

	yearEntry := widget.NewEntry()
	yearEntry.SetPlaceHolder("2024")

	techEntry := widget.NewEntry()
	techEntry.SetPlaceHolder("html, css, js, go (comma separated)")

	categorySelect := widget.NewSelect(projectsData.Categories, nil)
	if len(projectsData.Categories) > 0 {
		categorySelect.SetSelected(projectsData.Categories[0])
	}

	featuredCheck := widget.NewCheck("Featured Project (⭐)", nil)

	githubEntry := widget.NewEntry()
	githubEntry.SetPlaceHolder("https://github.com/...")

	liveEntry := widget.NewEntry()
	liveEntry.SetPlaceHolder("https://demo.com/...")

	releaseEntry := widget.NewEntry()
	releaseEntry.SetPlaceHolder("https://release.com/...")

	moveUpBtn := widget.NewButtonWithIcon("↑ Move Up", theme.MoveUpIcon(), func() {
		if selectedProjectIndex <= 0 {
			dialog.ShowInformation("Info", "Cannot move first item up", mainWindow)
			return
		}

		projectsData.Projects[selectedProjectIndex], projectsData.Projects[selectedProjectIndex-1] =
			projectsData.Projects[selectedProjectIndex-1], projectsData.Projects[selectedProjectIndex]

		selectedProjectIndex--
		projectList.Refresh()
		saveProjectsData()
	})

	moveDownBtn := widget.NewButtonWithIcon("↓ Move Down", theme.MoveDownIcon(), func() {
		if selectedProjectIndex < 0 || selectedProjectIndex >= len(projectsData.Projects)-1 {
			dialog.ShowInformation("Info", "Cannot move last item down", mainWindow)
			return
		}

		projectsData.Projects[selectedProjectIndex], projectsData.Projects[selectedProjectIndex+1] =
			projectsData.Projects[selectedProjectIndex+1], projectsData.Projects[selectedProjectIndex]

		selectedProjectIndex++
		projectList.Refresh()
		saveProjectsData()
	})

	editProjectBtn := widget.NewButtonWithIcon("✏️ Edit Selected", theme.DocumentCreateIcon(), func() {
		if selectedProjectIndex < 0 {
			dialog.ShowInformation("Info", "Please select a project to edit", mainWindow)
			return
		}

		populateProjectForm(
			titleEntry, descEntry, imageEntry, yearEntry, techEntry,
			githubEntry, liveEntry, releaseEntry,
			categorySelect, featuredCheck,
			selectedProject,
		)

		dialog.ShowInformation("Info", "You can now edit the project details", mainWindow)
	})

	updateProjectBtn := widget.NewButtonWithIcon("🔄 Update Project", theme.ConfirmIcon(), func() {
		if selectedProjectIndex < 0 {
			dialog.ShowInformation("Info", "Please select a project to update", mainWindow)
			return
		}

		if titleEntry.Text == "" || descEntry.Text == "" {
			dialog.ShowInformation("Error", "Title and Description are required", mainWindow)
			return
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

		titleEntry.SetText("")
		descEntry.SetText("")
		imageEntry.SetText("")
		yearEntry.SetText("")
		techEntry.SetText("")
		githubEntry.SetText("")
		liveEntry.SetText("")
		releaseEntry.SetText("")
		featuredCheck.SetChecked(false)

		dialog.ShowInformation("Success", "✅ Project updated successfully!", mainWindow)
	})

	addProjectBtn := widget.NewButtonWithIcon("➕ Add New Project", theme.ContentAddIcon(), func() {
		if titleEntry.Text == "" || descEntry.Text == "" {
			dialog.ShowInformation("Error", "Title and Description are required", mainWindow)
			return
		}

		if len(projectsData.Categories) == 0 {
			dialog.ShowInformation("Error", "Please add at least one category first", mainWindow)
			return
		}

		if categorySelect.Selected == "" {
			dialog.ShowInformation("Error", "Please select a category", mainWindow)
			return
		}

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

		titleEntry.SetText("")
		descEntry.SetText("")
		imageEntry.SetText("")
		yearEntry.SetText("")
		techEntry.SetText("")
		githubEntry.SetText("")
		liveEntry.SetText("")
		releaseEntry.SetText("")
		featuredCheck.SetChecked(false)

		dialog.ShowInformation("Success", "✅ Project added!", mainWindow)
	})

	deleteProjectBtn := widget.NewButtonWithIcon("🗑️ Delete Selected", theme.DeleteIcon(), func() {
		if selectedProjectIndex < 0 {
			dialog.ShowInformation("Info", "Please select a project to delete", mainWindow)
			return
		}

		dialog.ShowConfirm("Confirm Delete", fmt.Sprintf("Delete project '%s'?", selectedProject.Title), func(confirm bool) {
			if confirm {
				projectsData.Projects = append(
					projectsData.Projects[:selectedProjectIndex],
					projectsData.Projects[selectedProjectIndex+1:]...,
				)
				selectedProjectIndex = -1
				projectList.Refresh()
				saveProjectsData()
			}
		}, mainWindow)
	})

	categoryEntry := widget.NewEntry()
	categoryEntry.SetPlaceHolder("e.g., web, mobile, desktop, game")

	categoriesCount := widget.NewLabel(fmt.Sprintf("Total categories: %d", len(projectsData.Categories)))

	categoriesList := container.NewVBox()
	refreshCategoriesList := func() {
		categoriesList.Objects = nil
		for i, cat := range projectsData.Categories {
			catLabel := widget.NewLabel(fmt.Sprintf("%d. 🏷️ %s", i+1, cat))
			categoriesList.AddObject(catLabel)
		}
		categoriesCount.SetText(fmt.Sprintf("Total categories: %d", len(projectsData.Categories)))
		categoriesList.Refresh()
	}

	refreshCategoriesList()

	addCategoryBtn := widget.NewButton("➕ Add Category", func() {
		if categoryEntry.Text == "" {
			dialog.ShowInformation("Error", "Enter category name", mainWindow)
			return
		}

		for _, cat := range projectsData.Categories {
			if strings.EqualFold(cat, categoryEntry.Text) {
				dialog.ShowInformation("Info", "Category already exists", mainWindow)
				return
			}
		}

		projectsData.Categories = append(projectsData.Categories, categoryEntry.Text)

		categorySelect.Options = projectsData.Categories
		if len(projectsData.Categories) > 0 {
			categorySelect.SetSelected(projectsData.Categories[0])
		}
		categorySelect.Refresh()

		refreshCategoriesList()

		categoryEntry.SetText("")
		saveProjectsData()
		dialog.ShowInformation("Success", "✅ Category added!", mainWindow)
	})

	deleteCategoryBtn := widget.NewButtonWithIcon("🗑️ Delete Last Category", theme.DeleteIcon(), func() {
		if len(projectsData.Categories) == 0 {
			return
		}
		dialog.ShowConfirm("Confirm Delete", "Delete last category?\nProjects using this category will need reassignment!", func(confirm bool) {
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
			}
		}, mainWindow)
	})

	categoriesScroll := container.NewVScroll(categoriesList)
	categoriesScroll.SetMinSize(fyne.NewSize(300, 150))

	moveToolbar := container.NewHBox(
		moveUpBtn,
		moveDownBtn,
		editProjectBtn,
		updateProjectBtn,
	)

	form := container.NewMax(bg,
		container.NewVScroll(
			container.NewPadded(
				container.NewVBox(
					createDarkCard("🎯 PROJECT ACTIONS",
						container.NewVBox(
							widget.NewLabel("Select a project from the list, then:"),
							moveToolbar,
						),
					),

					createDarkCard("🏷️ MANAGE CATEGORIES",
						container.NewVBox(
							widget.NewLabelWithStyle("Add your own categories:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
							container.NewHBox(
								categoryEntry,
								addCategoryBtn,
							),
							widget.NewLabelWithStyle("Existing categories:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
							categoriesCount,
							categoriesScroll,
							deleteCategoryBtn,
						),
					),

					createDarkSeparator(),

					createDarkCard("💻 ADD / EDIT PROJECT",
						container.NewVBox(
							titleEntry,
							descEntry,
							container.NewGridWithColumns(2,
								imageEntry,
								yearEntry,
							),
							techEntry,
							container.NewVBox(
								widget.NewLabelWithStyle("Select Category:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
								categorySelect,
								featuredCheck,
							),
						),
					),

					createDarkCard("🔗 LINKS (optional)",
						container.NewVBox(
							githubEntry,
							liveEntry,
							releaseEntry,
						),
					),

					createDarkCard("⚡ ACTIONS",
						container.NewHBox(
							addProjectBtn,
							deleteProjectBtn,
						),
					),
				),
			),
		),
	)

	split := container.NewHSplit(
		container.NewBorder(
			createDarkCard("📋 PROJECTS LIST (click to select)", nil),
			nil, nil, nil,
			container.NewMax(bg, container.NewPadded(projectList)),
		),
		form,
	)
	split.Offset = 0.4

	return split
}

func createSaveTab() fyne.CanvasObject {
	bg := canvas.NewRectangle(darkBg)

	dir, _ := os.Getwd()

	stats := createDarkCard("📊 FILE STATISTICS",
		container.NewGridWithColumns(2,
			container.NewVBox(
				widget.NewLabelWithStyle("📝 Changelog:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel(fmt.Sprintf("%d entries", len(currentData.Entries))),
				widget.NewLabel(fmt.Sprintf("%d future features", len(currentData.FutureFeatures))),
			),
			container.NewVBox(
				widget.NewLabelWithStyle("💻 Projects:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel(fmt.Sprintf("%d projects", len(projectsData.Projects))),
				widget.NewLabel(fmt.Sprintf("%d categories", len(projectsData.Categories))),
			),
		),
	)

	files := createDarkCard("📁 FILE LOCATIONS",
		container.NewVBox(
			widget.NewLabel(fmt.Sprintf("📂 Directory: %s", dir)),
			widget.NewLabel(fmt.Sprintf("📄 %s", changelogFile)),
			widget.NewLabel(fmt.Sprintf("📄 %s", projectsFile)),
		),
	)

	actions := createDarkCard("⚙️ ACTIONS",
		container.NewVBox(
			widget.NewButtonWithIcon("💾 Save All Files", theme.DocumentSaveIcon(), func() {
				saveAllData()
			}),
			widget.NewButtonWithIcon("🔄 Reload Data", theme.ViewRefreshIcon(), func() {
				loadChangelogData()
				loadProjectsData()
				dialog.ShowInformation("Success", "Data reloaded!", mainWindow)
			}),
		),
	)

	return container.NewMax(
		bg,
		container.NewVScroll(
			container.NewPadded(
				container.NewVBox(
					stats,
					files,
					actions,
				),
			),
		),
	)
}
