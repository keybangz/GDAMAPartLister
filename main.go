package main

import (
	"fmt"
	"math"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	// FIXME: Remove after console portion of app is done and rewrite accordingly
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/signintech/gopdf"
)

// MINIMUM DOOR SIZE SINGLE
// 2400mm x 2400mm - 2400mm x 3600mm (could be wrong whiteboard says 3500mm)
// Panel height: 550 should be average, do we bother taking the remainder? probably not.

// These probably should be organized into some sort of door struct which the program can reference later
// Do we need to save some sort of door id for back reference in case they want a spreadsheet created?
// Am I overthinking it right now? probably.

// FIXME: 15/10/24 - SOFTWARE DUE FOR REFACTOR, ugly ass code.

// // For reference, every panel of door height should be 600mm high each, so 600 x 4 =
var goPDF *gopdf.GoPdf

type Door struct {
	mountType     int
	cableSize     float64
	wheelCount    int
	wheelType     bool
	midHingeType  bool
	midHingeCount int
}

func main() {
	// Handle PDF creator handles
	goPDF = &gopdf.GoPdf{}

	// Grab door size -> Mount Type -> Print static parts, then print dynamic parts.
	prog := app.New()
	mainWindow := prog.NewWindow("GDAMA Toolbox v0.1")

	eDoorHeight := widget.NewEntry()
	eDoorWidth := widget.NewEntry()
	ePanelHeight := widget.NewEntry()
	ePanelHeight.SetText("600") // Set default panel height to 550 for averaging.
	eOutput := widget.NewLabel("")
	ePrint := widget.NewButton("Print", nil)
	ePrint.Disable()

	var currentDoor Door

	eDoorType := widget.NewRadioGroup([]string{"Standard (STD)", "Front-mount", "Low Head-room (LHR) Rear Mount"}, func(value string) {
		if value == "Standard (STD)" {
			currentDoor.mountType = 1
		} else if value == "Front-mount" {
			currentDoor.mountType = 2
		} else if value == "Low Head-room (LHR) Rear Mount" {
			currentDoor.mountType = 3
		}
	})

	doorPartListerForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Door Height", Widget: eDoorHeight},
			{Text: "Door Width", Widget: eDoorWidth},
			{Text: "Panel Height", Widget: ePanelHeight},
			{Text: "Door Type", Widget: eDoorType},
		},
		OnSubmit: func() {
			// I imagine there are a few escape causes when tieing strings together like this
			// Not that it should it matter as the conversions are error checked
			panels := GetPanelCount(eDoorHeight.Text, ePanelHeight.Text)
			mountSpecs := fmt.Sprintf("Mount Specifications:\nDoor Height: %s Door Width: %s Panel Height: %s\nMount Type: %s\n%s\n", eDoorHeight.Text, eDoorWidth.Text, ePanelHeight.Text, eDoorType.Selected, panels)
			staticParts := StaticParts()
			dynamicParts := DynamicParts(currentDoor, eDoorWidth.Text, eDoorHeight.Text, ePanelHeight.Text)

			output := fmt.Sprintf("%s\nPart List: \n%s\n%s", mountSpecs, staticParts, dynamicParts)
			eOutput.SetText(output)
			ePrint.Enable()
			ePrint.OnTapped = func() {
				OnPressPrint(currentDoor, eDoorHeight.Text, eDoorWidth.Text, ePanelHeight.Text)
			}
		},
	}

	partListContent := container.NewVBox(doorPartListerForm)
	partListOutputContent := container.NewVBox(eOutput, ePrint)

	partLister := container.New(layout.NewGridLayout(1), partListContent, partListOutputContent)

	scrollPartLister := container.NewVScroll(partLister)

	tabs := container.NewAppTabs(
		container.NewTabItem("Part Lister", scrollPartLister),
		container.NewTabItem("Placeholder", widget.NewLabel("Placeholder")),
		container.NewTabItem("Placeholder", widget.NewLabel("Placeholder")),
	)

	tabs.SetTabLocation(container.TabLocationLeading)

	// text2.Move(fyne.NewPos(0, 20)) // This is to position on the layout, the app handles responsive change of the window size.

	// Set grid layout so Tabbed sections will be on the left
	// Middle content aka grid 1 will show tab contents
	// Right content aka grid 2 will show tab results (This case will be sectional door part list output.)
	appContent := tabs

	mainWindow.Resize(fyne.NewSize(800, 600)) // Resize window to a sane start size.
	mainWindow.SetContent(appContent)         // Set the content passed through
	mainWindow.ShowAndRun()                   // Profit
	// prog.Run()
	tidy()
}

// End GUI Loop
func tidy() {
	fmt.Println("DEBUG | Program Exited")
}

func GetPanelCount(doorHeight string, panelHeight string) (output string) {
	fHeight, err := strconv.ParseFloat(doorHeight, 64)
	fPanelHeight, err := strconv.ParseFloat(panelHeight, 64)

	fmt.Println(fPanelHeight, fHeight)

	// cheeky error handling by crashing the program lol
	if err != nil {
		panic(err)
	}

	var doorPanelCount float64 = fHeight / fPanelHeight // Divide door height by 600
	var realPanelCount int                              // Real Panel count showed in program and counting logically

	for i := 0.0; i < doorPanelCount; i++ {
		realPanelCount = int(i + 1)
	}

	output = fmt.Sprintf("Panels: %d", realPanelCount)

	return output
}

// Static Parts will be listed first as parts that either
// Do certain part calculations depending on mount type (cables, etc)
// 1) Do not change size depending on the door size
// 2) Do not change amount depending on the door size
func StaticParts() (output string) {
	output = fmt.Sprintf("2x Track Brackets\n(L + R) Flag Brackets\n(L + R) Cable Drums\n1x Center Bearing Plate")

	return output
}

func OnPressPrint(currentDoor Door, height string, width string, panelHeight string) {
	fmt.Println("Print Button Handler")

	goPDF.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	goPDF.AddPage()

	var fontPath string
	var runningOS string
	switch runtime.GOOS {
	case "windows":
		{
			fontPath = "C:\\Windows\\Fonts\\arial.ttf"
			runningOS = "windows"
		}
	case "darwin":
		{
			fontPath = "/Library/Fonts/Arial.ttf"
			runningOS = "mac"
		}
	default:
		{
			fontPath = "/usr/share/fonts/TTF/DejaVuSans.ttf"
			runningOS = "linux"
		}
	}

	fmt.Println(runningOS)

	var err error

	if runningOS == "windows" || runningOS == "mac" {
		err = goPDF.AddTTFFont("Arial", fontPath)

		if err != nil {
			fmt.Println("Error loading font Arial", err)
		}

		goPDF.SetFont("Arial", "", 18)
	} else if runningOS == "linux" {
		err = goPDF.AddTTFFont("DejaVu", fontPath)

		if err != nil {
			fmt.Println("Error loading font DejaVu Sans", err)
		}

		goPDF.SetFont("DejaVu", "", 18)
	}

	// 1240 x 1754 A4 in pixels (wrong apparently?)

	tempX := 4.0
	tempY := 5.0
	goPDF.SetX(tempX)
	goPDF.SetY(tempY)
	goPDF.Cell(nil, "Garage Doors & More | Sectional Door Partlist")

	goPDF.SetX(tempX)
	goPDF.SetY(tempY + 25.0)
	timeNow := time.Now()

	currentDate := timeNow.Format("2006-01-02")
	currentTime := timeNow.Format("15:04:05")

	generated := fmt.Sprintf("Generated on %s @ %s", currentDate, currentTime)
	goPDF.Cell(nil, generated)

	// Find RGB colors and hopefully no color means transparent / nothing, otherwise set to white.
	goPDF.SetStrokeColor(0, 0, 0)
	goPDF.SetLineWidth(2)
	goPDF.SetFillColor(0, 0, 0)

	fDoorHeight, err := strconv.ParseFloat(height, 64)
	if err != nil {
		fmt.Println(err.Error())
	}

	fDoorWidth, err := strconv.ParseFloat(width, 64)

	if err != nil {
		fmt.Println(err.Error())
	}

	fDoorHeight = fDoorHeight / 15.0
	fDoorWidth = fDoorWidth / 15.0

	// rec := goPDF.Rectangle(100, 100, fDoorWidth, fDoorHeight, "DF", 0, 0)
	// Add automatic padding for spacing of document
	// Padding might be 25 pixels? what are the constraints of an A4 document? idfk
	// Add width and height of door 25 pixels from desired side of door diagram
	// Part list can go underneath, How can I center the recentangle? Who fucking knows.

	opt := gopdf.CellOption{
		Align:  gopdf.Center | gopdf.Middle,
		Border: gopdf.AllBorders,
		Float:  gopdf.Center,
	}

	rec := &gopdf.Rect{
		W: fDoorWidth,
		H: fDoorHeight,
	}

	startX := 50.0
	startY := 60.0

	size := fmt.Sprintf("%s(h) x %s(w)", height, width)
	panels := GetPanelCount(height, panelHeight)

	goPDF.SetX(startX)
	goPDF.SetY(startY)
	goPDF.Cell(nil, panels)

	goPDF.SetX(startX)
	goPDF.SetY(startY + 25)
	goPDF.CellWithOption(rec, size, opt)

	// Part List, this is such a shitty way to do new lines its insane
	goPDF.SetX(startX)
	newY := rec.H + startY + 50
	oldY := newY
	goPDF.SetY(newY)

	staticParts := StaticParts()
	segments := strings.Split(staticParts, "\n")

	// FIXME: Clean up later
	for _, segment := range segments {
		goPDF.SetX(startX)
		goPDF.SetY(newY)
		goPDF.Cell(nil, segment)
		newY = oldY + 25.0
		oldY = newY
	}

	dynParts := DynamicParts(currentDoor, width, height, panelHeight)
	segments = strings.Split(dynParts, "\n")

	for _, segment := range segments {
		goPDF.SetX(startX)
		goPDF.SetY(newY)
		goPDF.Cell(nil, segment)
		newY = oldY + 25.0
		oldY = newY
	}

	goPDF.WritePdf("door.pdf")
	goPDF.Close()

	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", "door.pdf").Start()
	case "darwin":
		err = exec.Command("open", "door.pdf").Start()
	default:
		err = exec.Command("xdg-open", "door.pdf").Start()
	}
}

// This is where shit is gonna get messy
// Dynamic parts like Cable size will be dependant on the global door size entered
// Some of them are dependent on the mount type
func DynamicParts(currentDoor Door, width string, height string, panelHeight string) (output string) {
	// Convert to float for extra middle hinge count.
	fHeight, err := strconv.ParseFloat(height, 64)
	fWidth, err := strconv.ParseFloat(width, 64)
	fPanelHeight, err := strconv.ParseFloat(panelHeight, 64)

	// cheeky error handling by crashing the program lol
	if err != nil {
		panic(err)
	}

	// Sort part list for standard + front mount first.
	if currentDoor.mountType == 1 || currentDoor.mountType == 2 {
		currentDoor.cableSize = fHeight + 600.0
		output = fmt.Sprintf("(L + R) STD Bearing Plates\n2x STD Top Hinges\n(L + R) STD Bottom Hanger\n")
	} else if currentDoor.mountType == 3 {
		currentDoor.cableSize = fHeight*2 + 500
		output = fmt.Sprintf("(L + R) LHR Bearing Plates\n2x LHR Top Hinges\n(L + R) LHR Bottom Hangers\n")
	}

	// use temporary output to store current dynamic list
	var tempOut = output

	torsionWidth := fWidth + 500.0
	output = fmt.Sprintf("%s1x Torsion Pole - %f\n", tempOut, torsionWidth)

	tempOut = output // Update temporary buffer

	// Check doorsize and update accordingly
	currentDoor.midHingeCount = 9 // set MidHinge count to default for lowest possibility account
	midHingeAdd := 0
	currentDoor.wheelCount = 10      // Lowest possible amount of wheels is 10???
	currentDoor.wheelType = false    // Short wheels by default
	currentDoor.midHingeType = false // Single middle hinges by default
	loopCounter := 0.0

	var doorWidthMetre float64 = fWidth / 1000.0
	var doorPanelCount float64 = fHeight / fPanelHeight // Divide door height by 600
	fmt.Println("doorPanelCount: ", doorPanelCount)
	fmt.Println("doorWidthMetre: ", doorWidthMetre)

	// For every extra panel higher than 4 panels add 2 extra wheels
	// Four panels is anything over 3.5

	// doorHeight divided by 600 default panel height is going to always return a floating point number
	// Algorithmic will assume that anything over .5 of a divided doorHeight via 600 will mean it is the next rounded number up.

	tempPanelCount := 0
	tempPanelCheck := 0
	for i := 0.0; i < doorPanelCount; i++ {
		tempPanelCount = int(i + 1)

		if tempPanelCount > 4 {
			currentDoor.wheelCount += 2
		}

		// Four panel door will add three middle hinges per extra metre from 2000
		// Five panel = 4+ per metre
		// Six panel = 5+ per metre

		if tempPanelCount == 4 {
			midHingeAdd = 3
		} else if tempPanelCount == 5 {
			midHingeAdd = 4
		} else if tempPanelCount == 6 {
			midHingeAdd = 5
		} else if tempPanelCount == 7 {
			midHingeAdd = 6
		}
	}

	for i := 0.0; i < fWidth; i++ {
		loopCounter++

		if loopCounter < 3000.0 || tempPanelCount == 4 && loopCounter == 3000.0 {
			continue
		}

		if math.Mod(loopCounter, 1000.0) == 0 {
			fmt.Println("Thousandth found:", loopCounter)

			if tempPanelCheck == 0 {
				// THIS COULD BE WRONG??? BUT IT SPITS OUT DIVIDED NUMBERS THAT WORK
				if tempPanelCount == 5 {
					currentDoor.midHingeCount = 8 // add 4 to 8 to get right number for 5 panel doors
				} else if tempPanelCount == 6 {
					currentDoor.midHingeCount = 10
				} else if tempPanelCount == 7 {
					currentDoor.midHingeCount = 12
				} else if tempPanelCount == 8 {
					currentDoor.midHingeCount = 14
				}
				tempPanelCheck++
			}

			fmt.Println("Middle Hinges:", currentDoor.midHingeCount)
			currentDoor.midHingeCount += midHingeAdd
		}
	}

	fmt.Println("Middle Hinge Add:", midHingeAdd)
	fmt.Println("Panel Count:", tempPanelCount)

	// If door width is higher than 4.5m then set double hinges and long wheels
	// Here check if door size is 4.5m+ and do 4+ hinges instead of two
	if fWidth >= 4500.0 {
		currentDoor.wheelType = true    // Long Wheels
		currentDoor.midHingeType = true // Double Hinge
	} else if fWidth < 4500.0 {
		currentDoor.wheelType = false // we dont need to do this but fuck it do it anyways
		currentDoor.midHingeType = false
	}

	if currentDoor.wheelType {
		output = fmt.Sprintf("%s%dx Long Wheels\n", tempOut, currentDoor.wheelCount)
	} else {
		output = fmt.Sprintf("%s%dx Short Wheels\n", tempOut, currentDoor.wheelCount)
	}

	tempOut = output // Update temporary buffer

	if currentDoor.midHingeType {
		output = fmt.Sprintf("%s%dx Double Middle Hinges\n", tempOut, currentDoor.midHingeCount)
	} else {
		output = fmt.Sprintf("%s%dx Single Middle Hinges\n", tempOut, currentDoor.midHingeCount)
	}

	tempOut = output // Update temporary buffer

	output = fmt.Sprintf("%s2x Cables @ size (mm): %f\n", tempOut, currentDoor.cableSize)

	fmt.Println("--- DOOR END ---")

	return output
}
