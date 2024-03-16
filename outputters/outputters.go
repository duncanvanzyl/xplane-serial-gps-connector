package outputters

import (
	"github.com/duncanvanzyl/xplane-serial-gps-connector/nmea"
	"github.com/duncanvanzyl/xplane-serial-gps-connector/xplane"
)

// Outputter is the interface that takes a xplane Position and returns a string for use as output
// eg NMEA sentences
type Outputter interface {
	Output(xplane.Position) (string, error)
}

// GGA is an Outputter that returns a GPGGA NMEA sentence
type GGA struct{}

// Output returns a GPGGA NMEA sentence
func (g *GGA) Output(p xplane.Position) (string, error) {
	return nmea.ToGPGGA(p.Dat_lat, p.Dat_lon, p.Dat_ele), nil
}

// VTG is an Outputter that returns a GPVTG NMEA sentence
type VTG struct{}

// Output returns a GPVTG NMEA sentence
func (v *VTG) Output(p xplane.Position) (string, error) {
	return nmea.ToGPVTG(float64(p.Veh_psi_loc), p.SOG()), nil
}
