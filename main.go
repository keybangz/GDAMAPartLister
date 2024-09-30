package main

import (
	"fmt"
	"image/color"

	// FIXME: Remove after console portion of app is done and rewrite accordingly
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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
var cableSize int
var wheelCount int    // +1 Panel = 2+ Wheels
var wheelType bool    // If false, short wheels, if true, long wheels???? profit
var midHingeType bool // If false single hinge, if true double hinge
var midHingeCount int // Default to 9, divide door width by 1m and add extras when needed
var panelHeight int   // Panel size of door should be 600mm high each

func main() {
	// TODO: Take console application process and make useful for people who don't like black box consoles, aka everyone but me.
	// Grab door size -> Mount Type -> Print static parts, then print dynamic parts.
	// fmt.Println("GDAMA DOOR SIZE PART PICKER v0.1 (9/21/24) by wyattw")
	// DoorSize()

	prog := app.New()
	mainWindow := prog.NewWindow("GDAMA Part Picker v0.1")

	// Text needs a color input from a RGBA value passed through on creation... why?
	green := color.NRGBA{R: 0, G: 180, B: 0, A: 255}
	text1 := canvas.NewText("Hello", green)
	text2 := canvas.NewText("There", green)

	text2.Move(fyne.NewPos(0, 20)) // This is to position on the layout, the app handles responsive change of the window size.

	// Set grid layout so Tabbed sections will be on the left
	// Middle content aka grid 1 will show tab contents
	// Right content aka grid 2 will show tab results (This case will be sectional door part list output.)
	content := container.New(layout.NewGridLayout(2), text1, text2)

	mainWindow.Resize(fyne.NewSize(800, 500)) // Resize window to a sane start size.
	mainWindow.SetContent(content)            // Set the content passed through
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
func DoorSize() {
	fmt.Println("Enter height(mm) of door: ")

	var height int
	var width int

	fmt.Scanf("%d", &height)

	fmt.Print("Height of door: ", height, "(mm)\n")
	fmt.Println("Enter width(mm) of door.")

	fmt.Scanf("%d", &width)
	fmt.Print("Width of door: ", width, "(mm)\n")

	// FIXME: silly way to do it?????
	doorWidth = width
	doorHeight = height

	fmt.Println("Door size: ", doorWidth, " x ", doorHeight)
	MountType()
}

func MountType() {
	var choice int
	badinput := false

	for menu := true; menu; {
		if !badinput {
			fmt.Println("Select Door Mount Type: ")
			fmt.Println("1. Standard")
			fmt.Println("2. Front-Mount")
			fmt.Println("3. Rear-Mount Low Headroom")
			fmt.Scanf("%d", &choice)
		} else {
			fmt.Println("Invalid selection!")
			fmt.Scanf("%d", &choice)
		}

		switch choice {
		case 1:
			fmt.Println("Mount Type Selected: Standard")
			mountType = 1
			menu = false
		case 2:
			fmt.Println("Mount Type Selected: Front-Mount")
			mountType = 2
			menu = false
		case 3:
			fmt.Println("Mount Type Selected: Rear-Mount Low Headroom")
			mountType = 3
			menu = false
		default:
			badinput = true
		}
	}

	StaticParts()
}

// Static Parts will be listed first as parts that either
// Do certain part calculations depending on mount type (cables, etc)
// 1) Do not change size depending on the door size
// 2) Do not change amount depending on the door size
func StaticParts() {
	fmt.Println("PARTS LIST:")
	fmt.Println("Track Brackets")
	fmt.Println("(L + R) Flag Brackets")
	fmt.Println("(L + R) Cable Drums")
	fmt.Println("Center Bearing Plate")
	fmt.Println("Torsion Pole")
	// fmt.Println("2x Hinges")
	DynamicParts()
}

// This is where shit is gonna get messy
// Dynamic parts like Cable size will be dependant on the global door size entered
// Some of them are dependent on the mount type
func DynamicParts() {
	// Sort part list for standard + front mount first.
	if mountType == 1 || mountType == 2 {
		cableSize = doorHeight * 2
		fmt.Println("(L, R) STD Bearing Plates")
		fmt.Println("2x STD Top Hinges")
		fmt.Println("(L, R) STD Bottom Hangers")
	} else if mountType == 3 {
		cableSize = doorHeight*2 + 500
		fmt.Println("(L, R) LHR Bearing Plates")
		fmt.Println("2x LHR Top Hinges")
		fmt.Println("(L, R) LHR Bottom Hangers")
	}

	// Check doorsize and update accordingly
	midHingeCount = 9    // set MidHinge count to default for lowest possibility account
	wheelCount = 10      // Lowest possible amount of wheels is 10???
	wheelType = false    // Short wheels by default
	midHingeType = false // Single middle hinges by default
	panelHeight = 600    // Standard height for a panel

	var doorWidthMetre int = doorWidth / 1000
	var doorPanelCount int = doorHeight / panelHeight // Divide door height by 600
	fmt.Println("doorWidthMetre: ", doorWidthMetre)

	// For every extra panel higher than 4 panels add 2 extra wheels
	for i := 0; i < doorPanelCount; i++ {
		if i > 4 { // if door panel count is higher than 4
			wheelCount++
			wheelCount++
		}

		// FIXME: Go keeps giving me int unused when adding a value of 4, i'm probably doing it wrong but shitty syntax equals shitty problems
		if doorWidth >= 4500 && i > 4 { // if door panel count is higher than 4 and also checks if door is over 4.5m wide
			midHingeCount++
			midHingeCount++
			midHingeCount++
			midHingeCount++
		}
		continue
	}

	// this brokie, for every meter of width above 3.6m we need to add an extra middle hinge
	// should probably use a float here, or across the board and convert from the users input
	// FIXME: Get width restraints at work on monday
	for i := 4; i <= doorWidthMetre; i++ {
		fmt.Println("doorWidthMetre added: ", i)
		midHingeCount++
		continue
	}

	// FIXME: potentionally incorrect value
	// if doorWidth > 2500 {
	// 	midHingeCount += doorWidthMetre
	// }

	// If door width is higher than 4.5m then set double hinges and long wheels
	// Here check if door size is 4.5m+ and do 4+ hinges instead of two
	if doorWidth >= 4500 {
		wheelType = true    // Long Wheels
		midHingeType = true // Double Hinge
	}

	if wheelType {
		fmt.Println(wheelCount, "Long Wheels")
	} else {
		fmt.Println(wheelCount, "Short Wheels")
	}

	if midHingeType {
		fmt.Println(midHingeCount, "Double Middle Hinges")
	} else {
		fmt.Println(midHingeCount, "Single Middle Hinges")
	}

	fmt.Println("2x Cables @ size (mm):", cableSize)
}
