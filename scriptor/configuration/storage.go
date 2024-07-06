package configuration

import "scriptor.test/scriptor/errs"

type ConstantsMapIntsStorage struct {
	list map[string]int
}

func NewConstantsMapIntsStorage(elems map[string]int) *ConstantsMapIntsStorage {
	return &ConstantsMapIntsStorage{
		list: elems,
	}
}

func (st *ConstantsMapIntsStorage) MapElement(elName string) (int, error) {
	if number, ok := st.list[elName]; ok {
		return number, nil
	}
	return 0, errs.ErrElementNotFound(elName)
}
