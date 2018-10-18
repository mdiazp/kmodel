package kmodel

func fComa(x ...string) string {
	s := ""
	ln := len(x)
	for i := 0; i < ln; i++ {
		s += x[i]
		if i+1 < ln {
			s += ", "
		}
	}
	return s
}

func columnNames(o ObjectModel, pk ...bool) []string {
	names := o.ColumnNames()
	if len(pk) > 0 && pk[0] {
		names = append(names, o.PkeyName())
	}
	return names
}

func columnValues(o ObjectModel, pk ...bool) []interface{} {
	values := o.ColumnValues()
	if len(pk) > 0 && pk[0] {
		values = append(values, o.PkeyValue())
	}
	return values
}

func columnPointers(o ObjectModel, pk ...bool) []interface{} {
	pointers := o.ColumnPointers()
	if len(pk) > 0 && pk[0] {
		pointers = append(pointers, o.PkeyPointer())
	}
	return pointers
}
