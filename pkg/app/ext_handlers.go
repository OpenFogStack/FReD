package app

// ExtHandleCreateKeygroup creates a new keygroup with the specified name in Storage.
func (a *App) ExtHandleCreateKeygroup(kgname string) error {
	if err := checkParameters(kgname); err != nil {
		return err
	}

	return a.createKeygroup(kgname)
}

// ExtHandleDeleteKeygroup removes the keygroup with the specified name from Storage.
func (a *App) ExtHandleDeleteKeygroup(kgname string) error {
	if err := checkParameters(kgname); err != nil {
		return err
	}

	return a.deleteKeygroup(kgname)

	return nil
}

// ExtHandleRead returns an item with the specified id from the specified keygroup.
func (a *App) ExtHandleRead(kgname string, id string) (string, error) {
	if err := checkParameters(kgname, id); err != nil {
		return "", err
	}

	return a.read(kgname, id)
}

// ExtHandleUpdate updates the item with the specified id in the specified keygroup.
func (a *App) ExtHandleUpdate(kgname string, id string, data string) error {
	if err := checkParameters(kgname, id, data); err != nil {
		return err
	}

	return a.update(kgname, id, data)
}

// ExtHandleDelete deletes the item with the specified id from the specified keygroup.
func (a *App) ExtHandleDelete(kgname string, id string) error {
	if err := checkParameters(kgname, id); err != nil {
		return err
	}

	return a.delete(kgname, id)
}

