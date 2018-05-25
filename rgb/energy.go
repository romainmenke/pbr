package rgb

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
)

// Energy stores RGB light energy as a 3D Vector.
// TODO: replace all Energy instances with Vector3 (simpler, less duplication)
type Energy geom.Vector3

// TODO: type Energy struct { geom.Vector3 }

var Full, White = Energy{1, 1, 1}, Energy{1, 1, 1}
var Empty, Black = Energy{0, 0, 0}, Energy{0, 0, 0}

// Merged merges energy b into energy a with a given signal strength.
func (a Energy) Merged(b Energy, signal Energy) Energy {
	return Energy{a.X + b.X*signal.X, a.Y + b.Y*signal.Y, a.Z + b.Z*signal.Z}
}

func (a Energy) Compressed(n float64) (b Energy, scale float64) {
	max := math.Max(a.X, math.Max(a.Y, a.Z))
	scale = max / n
	return a.Scaled(n / max), scale
}

// Scaled returns energy a scaled by n.
func (a Energy) Scaled(n float64) Energy {
	return Energy{a.X * n, a.Y * n, a.Z * n}
}

func (a Energy) Zero() bool {
	return a.X == 0 && a.Y == 0 && a.Z == 0
}

func (a Energy) Plus(b Energy) Energy {
	return Energy{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (a Energy) Minus(b Energy) Energy {
	return Energy{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func (a Energy) Size() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}

func (a Energy) Limit(n float64) Energy {
	return Energy{math.Min(a.X, n), math.Min(a.Y, n), math.Min(a.Z, n)}
}

// RandomGain randomly amplifies or destroys a signal.
// Strong signals are less likely to be destroyed and gain less amplification.
// Weak signals are more likely to be destroyed but gain more amplification.
// This creates greater overall system throughput (higher energy per signal, fewer signals).
func (a Energy) RandomGain(rnd *rand.Rand) Energy {
	greatest := geom.Vector3(a).Greatest()
	if rnd.Float64() > greatest {
		return Energy{}
	}
	return a.Scaled(1 / greatest)
}

// Times returns energy a multiplied by energy b.
func (a Energy) Times(b Energy) Energy {
	return Energy{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

// Diff returns the difference in two Energies
func (a Energy) Variance(b Energy) float64 {
	d := geom.Vector3(a).Minus(geom.Vector3(b))
	return d.X*d.X + d.Y*d.Y + d.Z*d.Z
}

func (a Energy) Mean() float64 {
	return (a.X + a.Y + a.Z) / 3
}

func (a Energy) Max() float64 {
	return math.Max(a.X, math.Max(a.Y, a.Z))
}

func (a Energy) Lerp(b Energy, n float64) Energy {
	return Energy(geom.Vector3(a).Lerp(geom.Vector3(b), n))
}

func (a *Energy) Set(b Energy) {
	a.X = b.X
	a.Y = b.Y
	a.Z = b.Z
}

func (a *Energy) UnmarshalText(b []byte) error {
	v, err := geom.ParseVector3(string(b))
	if err != nil {
		return err
	}
	a.Set(Energy(v))
	return nil
}

func ParseEnergy(s string) (e Energy, err error) {
	v, err := geom.ParseVector3(s)
	return Energy(v), err
}
