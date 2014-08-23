/*
	The core package contains all the types, 
	functions and interface's that all connectors, 
	leaders, and workers will contain inside an 
	emd distribution.
 */
package core

// This contains simply a name as a string to 
// identify this unique component from another 
// inside an emd distribution.
type Core struct {
	Name_ string
}
