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
	"fmt"
)

func New(s string) error {
	return errors.New(s)
}

func Errorf(f string, args ...interface{}) error {
	return fmt.Errorf(f, args...)
}

func Wrapf(err error, f string, args ...interface{}) error {
	return Wrap(err, fmt.Sprintf(f, args...))
}
