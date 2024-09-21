package main

import (
	"fmt"

	// FIXME: Remove after console portion of app is done and rewrite accordingly
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

// MINIMUM DOOR SIZE SINGLE
// 2400mm x 2400mm - 2400mm x 3600mm (could be wrong whiteboard says 3500mm)

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
	// Grab door size -> Mount Type -> Print static parts, then print dynamic parts.
	fmt.Println("GDAMA DOOR SIZE PART PICKER v0.1 (9/21/24) by wyattw")

	DoorSize()

	// FIXME: Remove after console portion of app is done
	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()
}

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
	fmt.Println("Select Door Mount Type: ")
	fmt.Println("1. Standard")
	fmt.Println("2. Front-Mount")
	fmt.Println("3. Rear-Mount Low Headroom")

	var choice int

	fmt.Scanf("%d", &choice)

	switch choice {
	case 1:
		fmt.Println("Mount Type Selected: Standard")
		mountType = 1
	case 2:
		fmt.Println("Mount Type Selected: Front-Mount")
		mountType = 2
	case 3:
		fmt.Println("Mount Type Selected: Rear-Mount Low Headroom")
		mountType = 3
	default:
		fmt.Println("Invalid selection!")
	}

	StaticParts()
}

// Static Parts will be listed first as parts that either
// Do certain part calculations depending on mount type (cables, etc)
// 1) Do not change size depending on the door size
// 2) Do not change amount depending on the door size
func StaticParts() {
	fmt.Println("Parts List:")
	fmt.Println("Track Brackets")
	fmt.Println("(L + R) Flag Brackets")
	fmt.Println("(L + R) Cable Drums")
	fmt.Println("Center Bearing Plate")
	fmt.Println("Torsion Pole")
	fmt.Println("2x Hinges")
	DynamicParts()
}

// This is where shit is gonna get messy
// Dynamic parts like Cable size will be dependant on the global door size entered
// Some of them are dependent on the mount type
// I will rip my hair out doing this probably.
func DynamicParts() {
	// Sort part list for standard + front mount first.
	if mountType == 1 || mountType == 2 {
		cableSize = doorHeight * 2
	} else if mountType == 3 {
		cableSize = doorHeight*2 + 500
	}

	// Update all dynamic parts to their minimum counter parts, and check doorsize and update
	midHingeCount = 9    // set MidHinge count to default for lowest possibility account
	wheelCount = 10      // Lowest possible amount of wheels is 10???
	wheelType = false    // Short wheels by default
	midHingeType = false // Single middle hinges by default

	var doorWidthMetre int = doorWidth / 1000

	for i := 1; i < doorWidthMetre; i++ {
		doorWidthMetre++
		break
	}

	// FIXME: potentionally incorrect value
	if doorWidth > 2500 {
		midHingeCount += doorWidthMetre
	}

	// If door width is higher than 4.5m then set double hinges and long wheels
	// Here check if door size is 4.5m+ and do 4+ hinges instead of two
	if doorWidth > 4500 {
		wheelType = true    // Long Wheels
		midHingeType = true // Double Hinge
	}

	// FIXME: Update count here.
	if wheelType {
		fmt.Println("10x Long Wheels")
	} else {
		fmt.Println("10x Short Wheels")
	}

	fmt.Println("2x Cables @ size:", cableSize)
	fmt.Println(midHingeCount, "Middle Hinges") // Need to be dynamic type, 2.5m x 3.5m = 9 by default? 4.5m+ door becomes double hinges
}
