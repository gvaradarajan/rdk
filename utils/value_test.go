package utils

import (
	"testing"

	"github.com/pkg/errors"
	"go.viam.com/test"
)

func TestAssertType(t *testing.T) {
	one := 1
	_, err := AssertType[string](one)
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, err, test.ShouldBeError, NewUnexpectedTypeError[string](one))

	_, err = AssertType[myAssertIfc](one)
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, err, test.ShouldBeError, NewUnexpectedTypeError[myAssertIfc](one))

	asserted, err := AssertType[myAssertIfc](myAssertInt(one))
	test.That(t, err, test.ShouldBeNil)
	test.That(t, asserted.method1(), test.ShouldBeError, errors.New("cool 8)"))
}

type myAssertIfc interface {
	method1() error
}

type myAssertInt int

func (m myAssertInt) method1() error {
	return errors.New("cool 8)")
}
