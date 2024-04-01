package classify

import (
	"fmt"
	"math"
)

type Complex struct {
	Re float64 // the real part
	Im float64 // the imaginary part
}

// create a new object with the given real and imaginary parts
func NewComplex(real, imag float64) Complex {
	return Complex{Re: real, Im: imag}
}

// return a string representation of the invoking Complex object
func (c Complex) String() string {
	if c.Im == 0 {
		return fmt.Sprintf("%g", c.Re)
	}
	if c.Re == 0 {
		return fmt.Sprintf("%gi", c.Im)
	}
	if c.Im < 0 {
		return fmt.Sprintf("%g - %gi", c.Re, -c.Im)
	}
	return fmt.Sprintf("%g + %gi", c.Re, c.Im)
}

// return abs/modulus/magnitude
func (c Complex) Abs() float64 {
	return math.Hypot(c.Re, c.Im)
}

// return a new Complex object whose value is (this + b)
func (c Complex) Plus(b Complex) Complex {
	return Complex{c.Re + b.Re, c.Im + b.Im}
}

// return a new Complex object whose value is (this - b)
func (c Complex) Minus(b Complex) Complex {
	return Complex{c.Re - b.Re, c.Im - b.Im}
}

// return a new Complex object whose value is (this * b)
func (c Complex) Times(b Complex) Complex {
	return Complex{c.Re*b.Re - c.Im*b.Im, c.Re*b.Im + c.Im*b.Re}
}

// return a new Complex object whose value is the reciprocal of this
func (c Complex) Reciprocal() Complex {
	scale := c.Re*c.Re + c.Im*c.Im
	return Complex{c.Re / scale, -c.Im / scale}
}

// return a / b
func (c Complex) Divides(b Complex) Complex {
	return c.Times(b.Reciprocal())
}

// return a new Complex object whose value is the complex sine of this
func (c Complex) Sin() Complex {
	return Complex{math.Sin(c.Re) * math.Cosh(c.Im), math.Cos(c.Re) * math.Sinh(c.Im)}
}

// return a new Complex object whose value is the complex cosine of this
func (c Complex) Cos() Complex {
	return Complex{math.Cos(c.Re) * math.Cosh(c.Im), -math.Sin(c.Re) * math.Sinh(c.Im)}
}

// equals returns true if the given Complex object is equal to the receiver
func (c Complex) Equals(x Complex) bool {
	return c.Re == x.Re && c.Im == x.Im
}

// FFT computes the FFT of a complex sequence x[] of length n.
func FFT(x []Complex) []Complex {
	n := len(x)
	if n == 1 {
		return []Complex{x[0]}
	}

	if n%2 != 0 {
		panic("n is not a power of 2")
	}

	// Compute FFT of even terms
	even := make([]Complex, n/2)
	for k := 0; k < n/2; k++ {
		even[k] = x[2*k]
	}
	q := FFT(even)

	// Compute FFT of odd terms
	odd := even // Reuse the array
	for k := 0; k < n/2; k++ {
		odd[k] = x[2*k+1]
	}
	r := FFT(odd)

	// Combine
	y := make([]Complex, n)
	for k := 0; k < n/2; k++ {
		kth := -2 * math.Pi * float64(k) / float64(n)
		wk := Complex{math.Cos(kth), math.Sin(kth)}
		y[k] = q[k].Plus(wk.Times(r[k]))
		y[k+n/2] = q[k].Minus(wk.Times(r[k]))
	}
	return y
}