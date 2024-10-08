package main

import (
	"fmt"
	"strconv"

	// FIXME: Remove after console portion of app is done and rewrite accordingly
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// MINIMUM DOOR SIZE SINGLE
// 2400mm x 2400mm - 2400mm x 3600mm (could be wrong whiteboard says 3500mm)
// Panel height: 550 should be average, do we bother taking the remainder? probably not.

// These probably should be organized into some sort of door struct which the program can reference later
// Do we need to save some sort of door id for back reference in case they want a spreadsheet created?
// Am I overthinking it right now? probably.

// Globals
var doorWidth int  // Door Width in mm
var doorHeight int // Door Height in mm
var mountType int  // Mount Type Choice in int

// For reference, every panel of door height should be 600mm high each, so 600 x 4 =
// Parts
var cableSize float64
var wheelCount int    // +1 Panel = 2+ Wheels
var wheelType bool    // If false, short wheels, if true, long wheels???? profit
var midHingeType bool // If false single hinge, if true double hinge
var midHingeCount int // Default to 9, divide door width by 1m and add extras when needed
var panelHeight int   // Panel size of door should be 600mm high each

func main() {
	// Grab door size -> Mount Type -> Print static parts, then print dynamic parts.
	// fmt.Println("GDAMA DOOR SIZE PART PICKER v0.1 (9/21/24) by wyattw")
	// DoorSize()

	prog := app.New()
	mainWindow := prog.NewWindow("GDAMA Toolbox v0.1")

	eDoorHeight := widget.NewEntry()
	eDoorWidth := widget.NewEntry()
	ePanelHeight := widget.NewEntry()
	ePanelHeight.SetText("600") // Set default panel height to 550 for averaging.
	eOutput := widget.NewLabel("")

	eDoorType := widget.NewRadioGroup([]string{"Standard", "Front-mount", "Low Head-room Rear Mount"}, func(value string) {
		if value == "Standard" {
			mountType = 1
		} else if value == "Front-mount" {
			mountType = 2
		} else if value == "Low Head-room Rear Mount" {
			mountType = 3
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
			mountSpecs := fmt.Sprintf("Mount Specifications:\nDoor Height: %s Door Width: %s Panel Height: %s Mount Type: %s\n", eDoorHeight.Text, eDoorWidth.Text, ePanelHeight.Text, eDoorType.Selected)
			staticParts := StaticParts()
			dynamicParts := DynamicParts(eDoorWidth.Text, eDoorHeight.Text, ePanelHeight.Text)

			output := fmt.Sprintf("%s\nPart List: \n%s\n%s", mountSpecs, staticParts, dynamicParts)
			eOutput.SetText(output)
		},
	}

	partListContent := container.NewVBox(doorPartListerForm)
	partListOutputContent := container.NewVBox(eOutput)

	partLister := container.New(layout.NewGridLayout(1), partListContent, partListOutputContent)

	tabs := container.NewAppTabs(
		container.NewTabItem("Part Lister", partLister),
		container.NewTabItem("Placeholder", widget.NewLabel("Placeholder")),
		container.NewTabItem("Placeholder", widget.NewLabel("Placeholder")),
	)

	tabs.SetTabLocation(container.TabLocationLeading)

	// text2.Move(fyne.NewPos(0, 20)) // This is to position on the layout, the app handles responsive change of the window size.

	// Set grid layout so Tabbed sections will be on the left
	// Middle content aka grid 1 will show tab contents
	// Right content aka grid 2 will show tab results (This case will be sectional door part list output.)
	appContent := tabs

	mainWindow.Resize(fyne.NewSize(800, 500)) // Resize window to a sane start size.
	mainWindow.SetContent(appContent)         // Set the content passed through
	mainWindow.ShowAndRun()                   // Profit
	// prog.Run()
	tidy()
}

// End GUI Loop
func tidy() {
	fmt.Println("Debug Exit")
}

// CONSOLE PROGRAM BELOW....
// Function for getting DoorSize
// Here we will call appriopriate part calculation sizes later.
// func DoorSize() {
// 	fmt.Println("Enter height(mm) of door: ")

// 	var height int
// 	var width int

// 	fmt.Scanf("%d", &height)

// 	fmt.Print("Height of door: ", height, "(mm)\n")
// 	fmt.Println("Enter width(mm) of door.")

// 	fmt.Scanf("%d", &width)
// 	fmt.Print("Width of door: ", width, "(mm)\n")

// 	// FIXME: silly way to do it?????
// 	doorWidth = width
// 	doorHeight = height

// 	fmt.Println("Door size: ", doorWidth, " x ", doorHeight)
// 	MountType()
// }

// func MountType() {
// 	var choice int
// 	badinput := false

// 	for menu := true; menu; {
// 		if !badinput {
// 			fmt.Println("Select Door Mount Type: ")
// 			fmt.Println("1. Standard")
// 			fmt.Println("2. Front-Mount")
// 			fmt.Println("3. Rear-Mount Low Headroom")
// 			fmt.Scanf("%d", &choice)
// 		} else {
// 			fmt.Println("Invalid selection!")
// 			fmt.Scanf("%d", &choice)
// 		}

// 		switch choice {
// 		case 1:
// 			fmt.Println("Mount Type Selected: Standard")
// 			mountType = 1
// 			menu = false
// 		case 2:
// 			fmt.Println("Mount Type Selected: Front-Mount")
// 			mountType = 2
// 			menu = false
// 		case 3:
// 			fmt.Println("Mount Type Selected: Rear-Mount Low Headroom")
// 			mountType = 3
// 			menu = false
// 		default:
// 			badinput = true
// 		}
// 	}

// 	StaticParts()
// }

// Static Parts will be listed first as parts that either
// Do certain part calculations depending on mount type (cables, etc)
// 1) Do not change size depending on the door size
// 2) Do not change amount depending on the door size
func StaticParts() (output string) {
	// fmt.Println("PARTS LIST:")

	output = fmt.Sprintf("Track Brackets\n(L + R) Flag Brackets\n(L + R) Cable Drums\nCenter Bearing Plate\nTorsion Pole")
	// fmt.Println("Track Brackets")
	// fmt.Println("(L + R) Flag Brackets")
	// fmt.Println("(L + R) Cable Drums")
	// fmt.Println("Center Bearing Plate")
	// fmt.Println("Torsion Pole")
	// fmt.Println("2x Hinges")

	return output
}

// This is where shit is gonna get messy
// Dynamic parts like Cable size will be dependant on the global door size entered
// Some of them are dependent on the mount type
func DynamicParts(width string, height string, panelHeight string) (output string) {
	// iHeight, err := strconv.Atoi(height)
	// iWidth, err := strconv.Atoi(width)
	// iPanelHeight, err := strconv.Atoi(panelHeight)

	// Convert to float for extra middle hinge count.
	fHeight, err := strconv.ParseFloat(height, 64)
	fWidth, err := strconv.ParseFloat(width, 64)
	fPanelHeight, err := strconv.ParseFloat(panelHeight, 64)

	// cheeky error handling by crashing the program lol
	if err != nil {
		panic(err)
	}

	// Sort part list for standard + front mount first.
	if mountType == 1 || mountType == 2 {
		cableSize = fHeight * 2.0
		output = fmt.Sprintf("(L + R) STD Bearing Plates\n2x STD Top Hinges\n(L + R) STD Bottom Hanger\n")
		// fmt.Println("(L, R) STD Bearing Plates")
		// fmt.Println("2x STD Top Hinges")
		// fmt.Println("(L, R) STD Bottom Hangers")
	} else if mountType == 3 {
		cableSize = fHeight*2 + 500
		output = fmt.Sprintf("(L + R) LHR Bearing Plates\n2x LHR Top Hinges\n(L + R) LHR Bottom Hangers\n")
		// fmt.Println("(L, R) LHR Bearing Plates")
		// fmt.Println("2x LHR Top Hinges")
		// fmt.Println("(L, R) LHR Bottom Hangers")
	}

	// use temporary output to store current dynamic list
	var tempOut = output

	// Check doorsize and update accordingly
	midHingeCount = 9    // set MidHinge count to default for lowest possibility account
	wheelCount = 10      // Lowest possible amount of wheels is 10???
	wheelType = false    // Short wheels by default
	midHingeType = false // Single middle hinges by default
	// iPanelHeight = 600   // Standard height for a panel

	// doorCheck := false

	var doorWidthMetre float64 = fWidth / 1000.0
	var doorPanelCount float64 = fHeight / fPanelHeight // Divide door height by 600
	fmt.Println("doorPanelCount: ", doorPanelCount)
	fmt.Println("doorWidthMetre: ", doorWidthMetre)

	// For every extra panel higher than 4 panels add 2 extra wheels
	for i := 0.0; i < doorPanelCount; i++ {
		if i >= 4.0 { // if door panel count is higher than 4
			wheelCount += 2
			midHingeCount += 2 // fix
			fmt.Println("wheelCount added 2 && midHingeCount added 2")
		}

		if fWidth >= 4500.0 && i >= 4.0 { // if door panel count is higher than 4 and also checks if door is over 4.5m wide
			midHingeCount += 4
			fmt.Println("midHingeCount added 4")
		}
	}

	// this brokie, for every meter of width above 3.6m we need to add an extra middle hinge
	// should probably use a float here, or across the board and convert from the users input
	// FIXME: Get width restraints at work on monday

	for i := 3.6; i <= doorWidthMetre; i++ {
		if i <= 4.49 {
			midHingeCount += 1
		} else {
			midHingeCount += 2
		}
	}

	// If door width is higher than 4.5m then set double hinges and long wheels
	// Here check if door size is 4.5m+ and do 4+ hinges instead of two
	if fWidth >= 4500.0 {
		wheelType = true    // Long Wheels
		midHingeType = true // Double Hinge
		// midHingeCount += 4
	} else if fWidth < 4500.0 {
		wheelType = false // we dont need to do this but fuck it do it anyways
		midHingeType = false
		// midHingeCount += 2
	}

	if wheelType {
		//fmt.Println(wheelCount, "Long Wheels")
		output = fmt.Sprintf("%s%dx Long Wheels\n", tempOut, wheelCount)
	} else {
		// fmt.Println(wheelCount, "Short Wheels")x``
		output = fmt.Sprintf("%s%dx Short Wheels\n", tempOut, wheelCount)
	}

	tempOut = output // Update temporary buffer

	if midHingeType {
		// fmt.Println(midHingeCount, "Double Middle Hinges")
		output = fmt.Sprintf("%s%dx Double Middle Hinges\n", tempOut, midHingeCount)
	} else {
		// fmt.Println(midHingeCount, "Single Middle Hinges")
		output = fmt.Sprintf("%s%dx Single Middle Hinges\n", tempOut, midHingeCount)
	}

	tempOut = output

	//fmt.Println("2x Cables @ size (mm):", cableSize)
	output = fmt.Sprintf("%s2x Cables @ size (mm): %f\n", tempOut, cableSize)

	fmt.Println("--- DOOR END ---")

	return output
}
