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

package errors // import "cryptoscope.co/go/errors"

import (
	"errors"
	"runtime/debug"
)

type Error struct {
	stack []byte
	errs  []error
}

func (e *Error) Error() string {
	return e.errs[0].Error() + "\nat:\n" + string(e.stack)
}

func (e *Error) Errors() []error {
	return e.errs
}

func (e *Error) Stack() []byte {
	return e.stack
}

func Wrap(err error, cause string) error {
	if er, ok := err.(*Error); ok {
		er.errs = append(er.errs, errors.New(cause))
		return er
	}

	return &Error{
		stack: debug.Stack(),
		errs:  []error{err},
	}
}
