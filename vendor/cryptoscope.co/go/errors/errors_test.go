/*
   This file is part of errors.

   errors is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   errors is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with errors.  If not, see <http://www.gnu.org/licenses/>.
*/

package errors

import (
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {
	inner := errors.New("inner")
	err1 := Wrap(inner, "something didn't work")

	if len(err1.stack) == 0 {
		t.Error("err1.stack is of length zero")
	}

	if len(err1.errs) == 0 {
		t.Error("err1.errs is of length zero")
	}

	if err1.errs[0] != inner {
		t.Error("err1.errs[0] is not inner")
	}

	err2 := Wrap(err1, "something didn't work")
	if err2 != err1 {
		t.Error("err1 != err2")
	}

	if len(err2.errs) != 2 {
		t.Error("err2.errs is not of length 2")
	}

	if err2.Error() != "inner"+"\nat:\n"+string(err2.Stack()) {
		t.Errorf("wrapped error doesn't return original cause on Error() but %s", err2.Error())
	}
}
