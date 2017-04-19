package ui

import "github.com/cheggaaa/pb"

// StartBar creates a progress bar object and starts it.
func StartBar(total int64, prefix string, units pb.Units, preconfigure bool, start bool) *pb.ProgressBar {
	// Create the progress bar and defer its start.
	bar := pb.New64(total).Prefix(prefix)
	if start {
		defer bar.Start()
	}

	// Set the units.
	bar.SetUnits(units)

	// Configure it if the caller wishes.
	if preconfigure {
		bar.ShowSpeed = true
	}

	// Return the progress bar object.
	return bar
}
